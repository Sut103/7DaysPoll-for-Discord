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
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bwmarrin/discordgo"
)

type Lambda struct {
	publicKey ed25519.PublicKey
	endpoint  string
}

func NewLambda(key string, endpoint string) (*Lambda, error) {
	publicKey, err := hex.DecodeString(key)
	if err != nil {
		return nil, err
	}

	return &Lambda{
		publicKey,
		endpoint,
	}, nil
}

type lambdaSession struct {
	l *Lambda
}

func (s *lambdaSession) InteractionRespond(interaction *discordgo.Interaction, resp *discordgo.InteractionResponse, options ...discordgo.RequestOption) error {
	res, err := toAPIGatewayProxyResponse(resp, http.StatusOK)
	if err != nil {
		return err
	}
	err = s.l.request(res)
	return err
}

func (s *lambdaSession) InteractionResponse(interaction *discordgo.Interaction, options ...discordgo.RequestOption) (*discordgo.Message, error) {
	return &discordgo.Message{}, nil
}

func (s *lambdaSession) MessageReactionAdd(channelID, messageID, emojiID string, options ...discordgo.RequestOption) error {
	return nil
}

func toAPIGatewayProxyResponse(body *discordgo.InteractionResponse, statusCode int) (events.APIGatewayProxyResponse, error) {
	json, err := discordgo.Marshal(&body)
	if err != nil {
		log.Println(http.StatusInternalServerError, "json marshal error:", err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
	}
	return events.APIGatewayProxyResponse{Body: string(json), StatusCode: http.StatusOK}, nil
}

func (l *Lambda) request(body events.APIGatewayProxyResponse) error {
	r := strings.NewReader(body.Body)

	req, err := http.NewRequest("POST", l.endpoint, r)
	if err != nil {
		return err
	}
	req.Header.Set("Lambda-Runtime-Function-Response-Mode", "streaming")
	req.Header.Set("Transfer-Encoding", "chunked")

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		return err
	}
	return nil
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
		s := lambdaSession{
			l,
		}
		err := command.Poll(&s, &interaction)
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
		}
		return events.APIGatewayProxyResponse{StatusCode: http.StatusOK}, nil

	default:
		log.Printf("%+v", interaction)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest}, nil
	}
}

func (l *Lambda) Run() error {
	lambda.Start(l.handler)
	return nil
}
