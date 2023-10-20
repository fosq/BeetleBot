## Getting Started

1. Get your Bot token from [https://discord.com/build/app-developers](https://discord.com/build/app-developers)
2. Copy your server channel's ID where the formatted patch notes updates will be sent to
3. Enter your channel ID in bot/bot.go to
```
flag.StringVar(&ChannelId, "channelid", "CHANNEL ID NUMBER HERE", "Channel id where message is sent to")
```
4. Enter your Bot token in bot/bot.go to
```
flag.StringVar(&Token, "token", "TOKEN HERE", "Bot Token")
```

### Installation

To build BeetleBot as a Windows 64-bit .exe file, use the command:
```
go build GOOS=windows GOARCH=amd64
```
For Linux:
```
go build GOOS=linux GOARCH=amd64
```

Run the executable file