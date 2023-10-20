package bot

import (
	"discordbot/tft"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

var Token string
var Prefix string
var ChannelId string

func init() {
	flag.StringVar(&Token, "token", "TOKEN HERE", "Bot Token")
	flag.StringVar(&Prefix, "prefix", "!", "Chat prefix, e.g. '!'' for '!print'")
	flag.StringVar(&ChannelId, "channelid", "CHANNEL ID HERE", "Channel id where message is sent to")
	flag.Parse()
}

func StartBot() {
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Handler for !print
	dg.AddHandler(printNotes)

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	// Initial check for updates
	if tft.CheckForUpdate() {
		if found := foundPreviousMessage(dg); !found {
			sendUpdate(dg)
		}
	}

	ticker := time.NewTicker(10 * time.Second)
	go func() {
		for range ticker.C {
			found := foundPreviousMessage(dg)
			if tft.CheckForUpdate() || !found {
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
	if m.ChannelID != ChannelId {
		return
	}

	if m.Content == (Prefix + "print") {
		sendUpdate(dg)
	}
}

func foundPreviousMessage(dg *discordgo.Session) bool {
	messages, err := dg.ChannelMessages(ChannelId, 5, "", "", "")
	if err != nil {
		fmt.Println("error fetching channel messages,", err)
		return true
	}

	for _, message := range messages {
		if message.Content == tft.GetPatchNotes() {
			return true
		}
	}

	return false
}

func sendUpdate(dg *discordgo.Session) {
	dg.ChannelMessageSend(ChannelId, tft.GetPatchNotes())
}
