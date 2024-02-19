package main

import (
	"7DaysPoll-interactions/runtime"
	"log"
	"os"
)

func main() {
	token := os.Getenv("DISCORD_BOT_TOKEN")
	err := runtime.RunBot(token)
	if err != nil {
		log.Fatalln(err)
	}
}
