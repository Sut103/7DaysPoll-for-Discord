package main

import (
	"7DaysPoll/runtime"
	"log"
	"os"
)

func main() {
	token := os.Getenv("DISCORD_BOT_TOKEN")
	bot := runtime.NewBot(token)

	err := bot.Run()
	if err != nil {
		log.Fatalln(err)
	}
}
