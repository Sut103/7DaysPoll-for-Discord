package main

import (
	"7DaysPoll-interactions/runtime"
	"log"
	"os"
)

func main() {
	var r runtime.Runtime

	token := os.Getenv("DISCORD_BOT_TOKEN")
	r = runtime.NewBot(token)

	err := r.Run()
	if err != nil {
		log.Fatalln(err)
	}
}
