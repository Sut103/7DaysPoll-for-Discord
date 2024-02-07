package main

import (
	"7DaysPoll-interactions/runtime"
	"log"
	"os"
)

func main() {
	if os.Getenv("IS_LOCAL") == "true" {
		token := os.Getenv("DISCORD_BOT_TOKEN")
		err := runtime.RunBot(token)
		log.Fatalln(err)
		return
	}
	err := runtime.RunLambda()
	log.Fatalln(err)
}
