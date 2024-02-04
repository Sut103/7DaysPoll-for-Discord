package main

import (
	"bytes"
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

func verify(r *events.APIGatewayProxyRequest, key ed25519.PublicKey) bool {
	var msg bytes.Buffer

	xse, exists := r.Headers["x-signature-ed25519"]
	if !exists {
		return false
	}

	xst, exists := r.Headers["x-signature-timestamp"]
	if !exists {
		return false
	}

	sig, err := hex.DecodeString(xse)
	if err != nil {
		return false
	}

	if len(sig) != ed25519.SignatureSize {
		return false
	}

	msg.WriteString(xst)
	msg.WriteString(r.Body)

	return ed25519.Verify(key, msg.Bytes(), sig)
}

func pong() (events.APIGatewayProxyResponse, error) {
	body := discordgo.InteractionResponse{
		Type: discordgo.InteractionResponsePong,
	}
	json, err := discordgo.Marshal(body)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	log.Println(http.StatusOK, "PONG!")
	return events.APIGatewayProxyResponse{Body: string(json), StatusCode: http.StatusOK}, nil
}

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	publicKey := os.Getenv("DISCORD_PUBLIC_KEY")
	keyString, err := hex.DecodeString(publicKey)
	if err != nil {
		log.Println(http.StatusInternalServerError, "pub_key decode error", err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
	}

	if !verify(&event, keyString) {
		log.Println(http.StatusUnauthorized, "invalid request signature")
		return events.APIGatewayProxyResponse{StatusCode: http.StatusUnauthorized}, nil
	}

	interaction := discordgo.Interaction{}
	err = discordgo.Unmarshal([]byte(event.Body), &interaction)
	if err != nil {
		log.Println(http.StatusInternalServerError, "json.unmarshal error", err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
	}

	if err != nil {
		log.Println(http.StatusInternalServerError, "json.marshal error", err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
	}

	switch interaction.Type {
	case discordgo.InteractionApplicationCommand:
		log.Println(http.StatusOK, "COMMAND!")

	case discordgo.InteractionPing:
		return pong()
	}

	return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, nil
}

func main() {
	lambda.Start(handler)
}
