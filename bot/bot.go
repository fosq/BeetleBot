package bot

import (
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"syscall"
	"time"

	"discordbot/logs"
	"discordbot/tft"

	"github.com/bwmarrin/discordgo"
)

func StartBot() {
	dg, err := discordgo.New("Bot " + globalConfig.Token)
	if err != nil {
		logs.WriteLogFile("error creating Discord session,", err)
		return
	}

	// Handler for !print, !update
	dg.AddHandler(printNotes)
	dg.AddHandler(updateNotes)
	dg.AddHandler(purge)

	err = dg.Open()
	if err != nil {
		logs.WriteLogFile("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	logs.WriteLogFile("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	// Initial check for updates
	if tft.UpdatePatches() {
		if found := foundPreviousMessage(dg); !found {
			sendUpdate(dg)
		}
	}

	ticker := time.NewTicker(30 * time.Minute)
	go func() {
		for range ticker.C {
			if tft.UpdatePatches() {
				sendUpdate(dg)
			}
		}
	}()
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

func printNotes(dg *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == dg.State.User.ID {
		return
	}

	// Ignore all messages not in the given channel
	if m.ChannelID != globalConfig.ChannelId {
		return
	}

	if m.Content == (globalConfig.Prefix + "print") {
		logs.WriteLogFile("Received 'PRINT' command")
		sendUpdate(dg)
	}
}

func updateNotes(dg *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == dg.State.User.ID {
		return
	}

	// Ignore all messages not in the given channel
	if m.ChannelID != globalConfig.ChannelId {
		return
	}

	if m.Content == (globalConfig.Prefix + "update") {
		logs.WriteLogFile("Received 'UPDATE' command")
		if tft.UpdatePatches() {
			sendUpdate(dg)
			return
		}
		dg.ChannelMessageSend(globalConfig.ChannelId, "No new patch notes were found.")
	}
}

func foundPreviousMessage(dg *discordgo.Session) bool {
	notes := tft.GetPatchNotes()
	notesChunks := splitMessage(notes)

	messages, err := dg.ChannelMessages(globalConfig.ChannelId, 15, "", "", "")
	if err != nil {
		logs.WriteLogFile("error fetching channel messages,", err)
		return true
	}

	for i, message := range messages {
		if message.Content == notesChunks[len(notesChunks)-1] && i <= len(notesChunks) {
			logs.WriteLogFile("Found equal patch notes in previously sent message")
			return true
		}
	}
	return false
}

func sendUpdate(dg *discordgo.Session) {
	notes := tft.GetPatchNotes()

	if len(notes) >= 2000 {
		messageArr := splitMessage(notes)
		for i := range messageArr {
			dg.ChannelMessageSend(globalConfig.ChannelId, messageArr[i])
		}
	} else {
		dg.ChannelMessageSend(globalConfig.ChannelId, notes)
	}
	dg.ChannelMessageSend(globalConfig.ChannelId, "======================================================")
	dg.ChannelMessageSend(globalConfig.ChannelId, "======================================================")
}

// Split messages that reach Discord's message character limit of 2000
func splitMessage(input string) []string {
	// Regular expression for finding headers "__**EXAMPLE_HEADING**__"
	re := regexp.MustCompile(`__\*\*[^*]+\*\*__`)
	var chunks []string

	start := 0
	for start < len(input) {
		// Determine the end index for the current chunk
		end := start + 2000
		if end > len(input) {
			end = len(input)
		}

		// Find the last header before the end index
		lastHeaderIndex := start
		for _, indexes := range re.FindAllStringIndex(input[start:end], -1) {
			lastHeaderIndex = start + indexes[0]
		}

		// If a header is found and it's not at the very start of the chunk,
		// use it as the split point. Otherwise, split at the character limit.
		if lastHeaderIndex > start {
			chunks = append(chunks, "_ _\n"+input[start:lastHeaderIndex]) // "_ _" for a blank line
			start = lastHeaderIndex
		} else {
			chunks = append(chunks, "_ _\n"+input[start:end])
			start = end
		}
	}

	return chunks
}

func purge(dg *discordgo.Session, m *discordgo.MessageCreate) {
	var bulkDelete bool

	if m.Author.ID == dg.State.User.ID {
		return
	}

	if m.ChannelID != globalConfig.ChannelId {
		return
	}

	baseExpr := `^(%vpurge)(?: (\d+))?$`
	expr := fmt.Sprintf(baseExpr, globalConfig.Prefix)

	r, _ := regexp.Compile(expr)

	if r.MatchString(m.Content) {
		matches := r.FindStringSubmatch(m.Content)
		numToDelete, err := strconv.Atoi(matches[2])
		if err != nil {
			numToDelete = 1
		}

		if numToDelete > 20 {
			bulkDelete = true
		}

		if numToDelete > 100 {
			dg.ChannelMessageSend(m.ChannelID, "Too many messages to delete, select a number less than 100.")
			return
		}
		logs.WriteLogFile("Received 'PURGE' command")

		messagesToDelete, err := dg.ChannelMessages(m.ChannelID, numToDelete+1, "", "", "")
		if !logs.Check(err) {
			return
		}

		var messageIds []string
		for _, message := range messagesToDelete {
			if !bulkDelete {
				err := dg.ChannelMessageDelete(m.ChannelID, message.ID)
				if !logs.Check(err) {
					return
				}
			} else {
				messageIds = append(messageIds, message.ID)
			}
		}

		if bulkDelete {
			err = dg.ChannelMessagesBulkDelete(m.ChannelID, messageIds)
			if !logs.Check(err) {
				return
			}
		}
	}
}
