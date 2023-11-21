package main

import (
	"discordbot/bot"
	"discordbot/tft"
)

func main() {
	tft.UpdatePatches()
	bot.StartBot()
}
