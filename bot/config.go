package bot

import (
	"discordbot/logs"
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

var (
	globalConfig   Config
	logFileName    = "logs.txt"
	errorFileName  = "errors.txt"
	configFileName = "config.json"
)

type Config struct {
	Token     string `json:"discord_bot_token"`
	ChannelId string `json:"discord_channel_id"`
	Prefix    string `json:"prefix"`
}

func SetConfig() {
	//// logs.Check if any of the given files exist, creates them if not found
	// Logging file
	_, err := os.Stat(logFileName)
	if errors.Is(err, os.ErrNotExist) {
		createFile(logFileName)
	} else {
		logs.CheckDataRetention(30, logFileName)
	}

	// Error file
	_, err = os.Stat(errorFileName)
	if errors.Is(err, os.ErrNotExist) {
		createFile(errorFileName)
	} else {
		logs.CheckDataRetention(0, errorFileName)
	}

	// Configuration file creation if not found
	_, err = os.Stat(configFileName)
	if errors.Is(err, os.ErrNotExist) {
		logs.WriteLogFile(fmt.Sprintf("Configuration file '%v' not found. Creating a new one...", configFileName))
		createFile(configFileName)
		PromptAndSetConfig()
		return
	}

	file, err := os.ReadFile(configFileName)
	if !logs.Check(err) {
		os.Exit(1)
	}

	// Configuration file creation if file empty
	if len(file) == 0 {
		logs.WriteLogFile(fmt.Sprintf("Configuration file '%v' is empty. Creating a new one.\n", configFileName))
		PromptAndSetConfig()
		return
	}

	var config Config
	err = json.Unmarshal(file, &config)
	if !logs.Check(err) {
		fmt.Printf("Configuration file '%v' is corrupted. Please correct the file or re-enter the config prompts:\n", configFileName)
		PromptAndSetConfig()
		return
	}

	// Set config to global config variable
	globalConfig = config
}

func WriteConfig(config Config) {
	file, err := os.OpenFile(configFileName, os.O_WRONLY, 0644)
	if !logs.Check(err) {
		os.Exit(1)
	}
	defer file.Close()

	configJSON, err := json.MarshalIndent(config, "", "    ")
	if !logs.Check(err) {
		os.Exit(1)
	}

	file.Write(configJSON)
}

func askInputs() Config {
	var config Config

	fmt.Print("Enter your bot's Discord API token:\n")
	_, err := fmt.Scan(&config.Token)
	if !logs.Check(err) {
		os.Exit(1)
	}

	fmt.Print("\nEnter the channel, where the bot will be sending updates:\n")
	_, err = fmt.Scan(&config.ChannelId)
	if !logs.Check(err) {
		os.Exit(1)
	}

	fmt.Print("\nEnter the preferred prefix for commands (e.g. '!' for '!purge 3'):\n")
	_, err = fmt.Scan(&config.Prefix)
	if !logs.Check(err) {
		os.Exit(1)
	}
	fmt.Println()

	return config
}

func PromptAndSetConfig() {
	config := askInputs()
	WriteConfig(config)
	globalConfig = config
}

func createFile(name string) {
	file, err := os.Create(name)
	logs.Check(err)
	defer file.Close()
}
