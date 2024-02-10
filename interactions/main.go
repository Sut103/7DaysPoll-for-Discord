package main

import (
	"7DaysPoll-interactions/runtime"
	"log"
	"os"
)

func main() {
	if _, ok := os.LookupEnv("AWS_LAMBDA_RUNTIME_API"); ok {
		key := os.Getenv("DISCORD_PUBLIC_KEY")
		lambda, err := runtime.NewLambda(key)
		if err != nil {
			log.Fatalln(err)
			return
		}
		err = lambda.Run()
		if err != nil {
			log.Fatalln(err)
			return
		}
		return
	}

	token := os.Getenv("DISCORD_BOT_TOKEN")
	bot := runtime.NewBot(token)
	err := bot.Run()
	if err != nil {
		log.Fatalln(err)
	}
}
