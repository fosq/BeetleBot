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
		dg.ChannelMessageSend(ChannelId, tft.GetPatchNotes())
	}

	ticker := time.NewTicker(1 * time.Hour)
	go func() {
		for range ticker.C {
			if tft.CheckForUpdate() {
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

func sendUpdate(dg *discordgo.Session) {
	dg.ChannelMessageSend(ChannelId, tft.GetPatchNotes())
}
