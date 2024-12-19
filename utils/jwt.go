package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
)

var secretKey = "mysecretkey" // Ganti dengan kunci rahasia yang lebih aman

func ValidateJWT(token string) (map[string]interface{}, error) {
	// Pisahkan token menjadi 3 bagian: header, payload, dan signature
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid token format")
	}

	headerEncoded, payloadEncoded, signatureEncoded := parts[0], parts[1], parts[2]

	// Validasi signature
	tokenWithoutSignature := headerEncoded + "." + payloadEncoded
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(tokenWithoutSignature))
	expectedSignature := base64.RawURLEncoding.EncodeToString(h.Sum(nil))

	if signatureEncoded != expectedSignature {
		return nil, errors.New("invalid token signature")
	}

	// Decode payload
	payloadJSON, err := base64.RawURLEncoding.DecodeString(payloadEncoded)
	if err != nil {
		return nil, errors.New("invalid payload encoding")
	}

	// Parse payload menjadi map
	var payload map[string]interface{}
	if err := json.Unmarshal(payloadJSON, &payload); err != nil {
		return nil, errors.New("invalid payload JSON")
	}

	return payload, nil
}
