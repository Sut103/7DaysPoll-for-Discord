package main

import (
	"7DaysPoll-interactions/command"
	"7DaysPoll-interactions/ping"
	"7DaysPoll-interactions/util"
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bwmarrin/discordgo"
)

var publicKey ed25519.PublicKey

func init() {
	key, err := hex.DecodeString(os.Getenv("DISCORD_PUBLIC_KEY"))
	if err != nil {
		log.Fatalln(err)
	}
	publicKey = key
}

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if !util.Verify(&event, publicKey) {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusUnauthorized}, nil
	}

	interaction := discordgo.Interaction{}
	discordgo.Unmarshal([]byte(event.Body), &interaction)

	switch interaction.Type {
	case discordgo.InteractionPing:
		return ping.Pong()

	case discordgo.InteractionApplicationCommand:
		return command.Poll(&interaction)

	default:
		log.Printf("%+v", interaction)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest}, nil
	}
}

func main() {
	lambda.Start(handler)
}
