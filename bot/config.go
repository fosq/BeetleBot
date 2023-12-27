package bot

import (
	"discordbot/helpers"
	"discordbot/logs"
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

var (
	globalConfig Config
)

type Config struct {
	Token     string `json:"discord_bot_token"`
	ChannelId string `json:"discord_channel_id"`
	Prefix    string `json:"prefix"`
}

var (
	LogFileName    = "logs.txt"
	ErrorFileName  = "errors.txt"
	ConfigFileName = "config.json"
)

func SetConfig() {
	//// Check if any of the given files exist, creates them if not found
	// Logging file
	_, err := os.Stat(LogFileName)
	if errors.Is(err, os.ErrNotExist) {
		helpers.CreateFile(LogFileName)
		logs.WriteLogFile(fmt.Sprintf("Created a new logging file %v.\n",
			LogFileName))
	} else {
		logs.CheckDataRetention(30, LogFileName)
	}

	// Error file
	_, err = os.Stat(ErrorFileName)
	if errors.Is(err, os.ErrNotExist) {
		helpers.CreateFile(ErrorFileName)
		logs.WriteLogFile(fmt.Sprintf("Created a new error logging file '%v'.\n",
			ErrorFileName))
	} else {
		logs.CheckDataRetention(0, ErrorFileName)
	}

	// Configuration file
	_, err = os.Stat(ConfigFileName)
	if errors.Is(err, os.ErrNotExist) {
		logs.WriteLogFile(fmt.Sprintf("Configuration file '%v' not found. Creating a new one...\n",
			ConfigFileName))
		helpers.CreateFile(ConfigFileName)

		logs.WriteLogFile(fmt.Sprintf("Created a new configuration file '%v'.\n",
			ConfigFileName))

		PromptAndSetConfig()
		logs.WriteLogFile("Configuration file filled.\n")
		return
	}

	file, err := os.ReadFile(ConfigFileName)
	if !logs.Check(err) {
		os.Exit(1)
	}

	// Configuration file creation if file empty
	if len(file) == 0 {
		logs.WriteLogFile(fmt.Sprintf("Configuration file '%v' is empty. Creating a new one.\n",
			ConfigFileName))
		PromptAndSetConfig()
		return
	}

	var config Config
	err = json.Unmarshal(file, &config)
	if !logs.Check(err) {
		logs.WriteLogFile(fmt.Sprintf("Configuration file '%v' is corrupted. Please correct the file or re-enter the config prompts:\n",
			ConfigFileName))
		PromptAndSetConfig()
		logs.WriteLogFile("Configuration file filled.\n")
		return
	}

	// Set config to global config variable
	globalConfig = config
}

func WriteConfig(config Config) {
	file, err := os.OpenFile(ConfigFileName, os.O_WRONLY, 0644)
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

func askInputs() {
	fmt.Print("Enter your bot's Discord API token:\n")
	_, err := fmt.Scan(&globalConfig.Token)
	if !logs.Check(err) {
		os.Exit(1)
	}

	fmt.Print("\nEnter the channel, where the bot will be sending updates:\n")
	_, err = fmt.Scan(&globalConfig.ChannelId)
	if !logs.Check(err) {
		os.Exit(1)
	}

	fmt.Print("\nEnter the preferred prefix for commands (e.g. '!' for '!purge 3'):\n")
	_, err = fmt.Scan(&globalConfig.Prefix)
	if !logs.Check(err) {
		os.Exit(1)
	}
	fmt.Println()
}

func PromptAndSetConfig() {
	askInputs()
	WriteConfig(globalConfig)
}
