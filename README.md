## About

BeetleBot is a simple lolchess.gg scraper with Discord API integration which fetches TFT patch notes every half an hour and sends the latest patch notes to the designated Discord channel. 

## Getting Started

1. Clone the repo
```
git clone https://github.com/fosq/BeetleBot.git
```
2. Get your Bot token from [https://discord.com/build/app-developers](https://discord.com/build/app-developers)
3. Copy your server channel's ID where the formatted patch notes updates will be sent to
4. Enter your bot token and channel id when the terminal prompts you to upon running the bot

### Installation

To build BeetleBot as a Windows 64-bit .exe file, use the command:
```
GOOS=windows GOARCH=amd64 go build
```
For Linux:
```
GOOS=linux GOARCH=amd64 go build
```
OR
Download from Releases

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

To purge previous messages (up to 100), type:
```
!purge <number> (e.g. !purge 30)
```