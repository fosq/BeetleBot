## About

BeetleBot is a simple lolchess.gg scraper with Discord API integration which fetches TFT patch notes every hour and sends the latest patch notes (if !print was invoked or patch was newer than previous saved patch) to the designated Discord channel. 

## Getting Started

1. Clone the repo
```
git clone https://github.com/fosq/BeetleBot.git
```
2. Get your Bot token from [https://discord.com/build/app-developers](https://discord.com/build/app-developers)
3. Copy your server channel's ID where the formatted patch notes updates will be sent to
4. Enter your channel ID in bot/bot.go to
```
flag.StringVar(&ChannelId, "channelid", "CHANNEL ID NUMBER HERE", "Channel id where message is sent to")
```
5. Enter your Bot token in bot/bot.go to
```
flag.StringVar(&Token, "token", "TOKEN HERE", "Bot Token")
```

### Installation

To build BeetleBot as a Windows 64-bit .exe file, use the command:
```
GOOS=windows GOARCH=amd64 go build
```
For Linux:
```
GOOS=linux GOARCH=amd64 go build
```

Run the executable file

### Usage

BeetleBot automatically scrapes latest patches every hour and sends the newest patch notes if a difference to the previous saved patch was found or
```
!print
```
is invoked. (doesn't check for latest patch, just sends the currently saved patch notes)