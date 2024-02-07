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
		if err != nil {
			log.Fatalln(err)
			return
		}
		return
	}
	err := runtime.RunLambda()
	log.Fatalln(err)
}
