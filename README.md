## About

BeetleBot is a simple lolchess.gg scraper with Discord API integration which fetches TFT patch notes every half an hour and sends the latest patch notes to the designated Discord channel. 

## Getting Started
1. Download the executable file for your operating system in [Releases](https://github.com/fosq/BeetleBot/releases) tab.
2. Get your bot token from [https://discord.com/build/app-developers](https://discord.com/build/app-developers).
3. Get the `channel id` where the patch notes will be sent to.
4. Enter your `bot token` and `channel id` when the terminal prompts you to do so upon running the executable.
<br>

## Compiling your own build
1. Clone the repository
```
git clone https://github.com/fosq/BeetleBot.git
```

2. Find your platform for your operating system and architecture from `go tool`:
```
go tool dist list
```

3. Compile the build to your platform with the command:
```
GOOS=youroperatingsystem GOARCH=yourarchitecture go build
```


---
- e.g. For Windows 64-bit:
```
GOOS=windows GOARCH=amd64 go build
```
- For Linux:
```
GOOS=linux GOARCH=amd64 go build
```

4. Run the executable file with the name `discordbot`

## Usage

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