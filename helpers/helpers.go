package helpers

import (
	"discordbot/logs"
	"os"
)

func CreateFile(name string) {
	file, err := os.Create(name)
	logs.Check(err)
	defer file.Close()
}
