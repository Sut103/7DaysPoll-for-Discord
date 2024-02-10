package main

import (
	"7DaysPoll-interactions/runtime"
	"log"
	"os"
)

func main() {
	if _, ok := os.LookupEnv("AWS_LAMBDA_RUNTIME_API"); ok {
		key := os.Getenv("DISCORD_PUBLIC_KEY")
		err := runtime.RunLambda(key)
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
