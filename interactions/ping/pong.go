package ping

import (
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/bwmarrin/discordgo"
)

func Pong() (events.APIGatewayProxyResponse, error) {
	body := discordgo.InteractionResponse{
		Type: discordgo.InteractionResponsePong,
	}

	json, err := discordgo.Marshal(body)
	if err != nil {
		log.Println(http.StatusInternalServerError, "json marshal error:", err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
	}

	return events.APIGatewayProxyResponse{Body: string(json), StatusCode: http.StatusOK}, nil
}
