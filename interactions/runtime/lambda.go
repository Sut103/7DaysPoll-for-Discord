package runtime

import (
	"7DaysPoll-interactions/command"
	"7DaysPoll-interactions/ping"
	"7DaysPoll-interactions/util"
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bwmarrin/discordgo"
)

type Lambda struct {
	publicKey ed25519.PublicKey
}

func NewLambda(key string) (*Lambda, error) {
	publicKey, err := hex.DecodeString(key)
	if err != nil {
		return nil, err
	}

	return &Lambda{
		publicKey,
	}, nil
}

func toAPIGatewayProxyResponse(body *discordgo.InteractionResponse, statusCode int) (events.APIGatewayProxyResponse, error) {
	json, err := discordgo.Marshal(&body)
	if err != nil {
		log.Println(http.StatusInternalServerError, "json marshal error:", err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
	}
	return events.APIGatewayProxyResponse{Body: string(json), StatusCode: http.StatusOK}, nil
}

func (l *Lambda) handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if !util.Verify(&event, l.publicKey) {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusUnauthorized}, nil
	}

	interaction := discordgo.Interaction{}
	discordgo.Unmarshal([]byte(event.Body), &interaction)

	switch interaction.Type {
	case discordgo.InteractionPing:
		return ping.Pong()

	case discordgo.InteractionApplicationCommand:
		body, err := command.Poll(&interaction)
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
		}
		return toAPIGatewayProxyResponse(body, http.StatusOK)

	default:
		log.Printf("%+v", interaction)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest}, nil
	}
}

func (l *Lambda) Run() error {
	lambda.Start(l.handler)
	return nil
}
