package main

import (
	"7DaysPoll-interactions/runtime"
	"log"
	"os"
)

func main() {
	if _, ok := os.LookupEnv("AWS_LAMBDA_RUNTIME_API"); ok {
		err := runtime.RunLambda()
		if err != nil {
			log.Fatalln(err)
			return
		}
		return
	}

	token := os.Getenv("DISCORD_BOT_TOKEN")
	err := runtime.RunBot(token)
	if err != nil {
		log.Fatalln(err)
	}
}
