## About

BeetleBot is a simple lolchess.gg scraper with Discord API integration which fetches TFT patch notes every half an hour and sends the latest patch notes to the designated Discord channel. 

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

BeetleBot automatically scrapes latest patches every half an hour and sends the newest patch notes if a difference to the previously saved patch was found OR when
```
!print
```
is invoked. (doesn't check for latest patch, just sends the currently saved patch notes)

To check for updates immediately, type:
```
!update
```