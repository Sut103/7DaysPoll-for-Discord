package util

import (
	"bytes"
	"crypto/ed25519"
	"encoding/hex"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
)

func Verify(r *events.APIGatewayProxyRequest, pk ed25519.PublicKey) bool {
	if os.Getenv("ENV") == "dev" {
		return true
	}

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

	var msg bytes.Buffer
	msg.WriteString(xst)
	msg.WriteString(r.Body)

	return ed25519.Verify(pk, msg.Bytes(), sig)
}

func GetTimeZone(lang string) (*time.Location, error) {
	timezone := map[string]string{
		"Japanese": "Asia/Tokyo",
	}

	tz, ok := timezone[lang]
	if !ok {
		return time.Local, nil
	}

	return time.LoadLocation(tz)
}
