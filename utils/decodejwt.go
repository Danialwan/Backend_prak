package utils

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// DecodeJWT untuk memproses payload JWT
func DecodeJWT(token string) (map[string]interface{}, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid JWT format")
	}

	// Decode payload
	payloadDecoded, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, errors.New("failed to decode JWT payload")
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(payloadDecoded, &payload); err != nil {
		return nil, errors.New("failed to parse JWT payload")
	}

	// Log isi payload untuk debugging
	fmt.Printf("Decoded JWT payload: %+v\n", payload)

	return payload, nil
}
