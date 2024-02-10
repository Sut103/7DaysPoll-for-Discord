package main

import (
	"7DaysPoll-interactions/runtime"
	"log"
	"os"
)

func main() {
	var r runtime.Runtime

	if _, ok := os.LookupEnv("AWS_LAMBDA_RUNTIME_API"); ok {
		key := os.Getenv("DISCORD_PUBLIC_KEY")
		lambda, err := runtime.NewLambda(key)
		if err != nil {
			log.Fatalln(err)
			return
		}
		r = lambda
	} else {
		token := os.Getenv("DISCORD_BOT_TOKEN")
		r = runtime.NewBot(token)
	}

	err := r.Run()
	if err != nil {
		log.Fatalln(err)
	}
}
