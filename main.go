package main

import (
	tft "discordbot/bot"
)

func main() {
	tft.CreatePatchNotes()
	if tft.ComparePatchNotes() {
		// send bot data
	}
}
