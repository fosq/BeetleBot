package logs

import (
	"discordbot/bot"
	"fmt"
	"log"
	"os"
	"time"
)

func WriteLogFile(a ...any) {
	file, err := os.OpenFile(bot.LogFileName, os.O_WRONLY|os.O_APPEND, 0644)
	if !Check(err) {
		os.Exit(1)
	}
	defer file.Close()
	log.SetOutput(file)

	for i := range a {
		log.Print(a[i])
	}
}

func WriteErrorLogFile(mainErr error) {
	file, err := os.OpenFile(bot.ErrorFileName, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	log.SetOutput(file)
	log.Println(mainErr)
}

func CheckDataRetention(daysToKeep int, fileName string) {
	info, err := os.Stat(fileName)
	if !Check(err) {
		os.Exit(1)
	}

	modTime := info.ModTime()
	deletionTime := time.Now().AddDate(0, 0, -daysToKeep)

	if modTime.Before(deletionTime) || daysToKeep == 0 {
		file, err := os.Create(fileName)
		Check(err)
		defer file.Close()
		WriteLogFile(fmt.Sprintf("The file %v has been cleared because the data stored is older than %v days.", fileName, daysToKeep))
	}
}

func Check(err error) bool {
	if err != nil {
		WriteErrorLogFile(err)
		WriteLogFile(err)
		fmt.Println(err)
		return false
	}
	return true
}
