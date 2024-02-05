package main

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bwmarrin/discordgo"
)

func Verify(r *events.APIGatewayProxyRequest, key ed25519.PublicKey) bool {
	if os.Getenv("ENV") == "dev" {
		return true
	}

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

		log.Println(http.StatusInternalServerError, "json marshal error:", err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
	}

	return events.APIGatewayProxyResponse{Body: string(json), StatusCode: http.StatusOK}, nil
}

func get7Days(day time.Time) [7]time.Time {
	days := [7]time.Time{}

	for i := range days {
		days[i] = day.AddDate(0, 0, i)
	}

	return days
}

func getTZ(lang string) (*time.Location, error) {
	timezone := map[string]string{
		"Japanese": "Asia/Tokyo",
	}

	tz, ok := timezone[lang]
	if !ok {
		return time.Local, nil
	}

	return time.LoadLocation(tz)
}

func poll(Locale string) (events.APIGatewayProxyResponse, error) {
	content := ""
	timezone, err := getTZ(Locale)
	if err != nil {
		log.Println(http.StatusInternalServerError, "timezone error", err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
	}
	days := get7Days(time.Now().Local().In(timezone))

	emoji7 := []string{
		"one",
		"two",
		"three",
		"four",
		"five",
		"six",
		"seven",
	}

	for i, day := range days {
		emoji := emoji7[i]
		content += fmt.Sprintf(":%s:: %d/%d/%d\n", string(emoji), day.Year(), day.Month(), day.Day())
	}

	body := discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	}

	json, err := discordgo.Marshal(body)
	if err != nil {
		log.Println(http.StatusInternalServerError, "json marshal error:", err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
	}

	return events.APIGatewayProxyResponse{Body: string(json), StatusCode: http.StatusOK}, nil
}

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	publicKey := os.Getenv("DISCORD_PUBLIC_KEY")
	keyString, err := hex.DecodeString(publicKey)
	if err != nil {
		log.Println(http.StatusInternalServerError, "publickey decode error:", err)
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
	}

	if !Verify(&event, keyString) {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusUnauthorized}, nil
	}

	interaction := discordgo.Interaction{}
	discordgo.Unmarshal([]byte(event.Body), &interaction)
	// if err != nil {
	// 	log.Println(http.StatusInternalServerError, "json unmarshal error:", err)
	// 	return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
	// }

	switch interaction.Type {
	case discordgo.InteractionApplicationCommand:
		return poll(discordgo.Locales[interaction.Locale])

	case discordgo.InteractionPing:
		return pong()

	default:
		log.Println(http.StatusBadRequest, "bad request")
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest}, nil
	}
}

func main() {
	lambda.Start(handler)
}
