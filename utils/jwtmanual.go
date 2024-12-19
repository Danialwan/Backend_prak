// package utils

// import (
// 	"crypto/hmac"
// 	"crypto/sha256"
// 	"encoding/base64"
// 	"encoding/json"
// 	"time"
// )

// // Fungsi untuk membuat JWT manual
// func GenerateJWT(username string) (string, error) {
// 	// Buat header
// 	header := map[string]string{
// 		"alg": "HS256",
// 		"typ": "JWT",
// 	}
// 	headerJSON, _ := json.Marshal(header)
// 	headerEncoded := base64.RawURLEncoding.EncodeToString(headerJSON)

// 	// Buat payload (misal dengan username dan exp)
// 	payload := map[string]interface{}{
// 		"username": username,
// 		"exp":      time.Now().Add(time.Hour * 1).Unix(), // token akan expire dalam 1 jam
// 	}
// 	payloadJSON, _ := json.Marshal(payload)
// 	payloadEncoded := base64.RawURLEncoding.EncodeToString(payloadJSON)

// 	// Gabungkan header dan payload
// 	token := headerEncoded + "." + payloadEncoded

// 	// Buat signature dengan HMAC-SHA256
// 	secretKey := "mysecretkey" // ganti dengan kunci rahasia yang lebih kuat
// 	h := hmac.New(sha256.New, []byte(secretKey))
// 	h.Write([]byte(token))
// 	signature := base64.RawURLEncoding.EncodeToString(h.Sum(nil))

// 	// Gabungkan semuanya
// 	jwtToken := token + "." + signature

// 	return jwtToken, nil
// }

package utils

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Fungsi untuk membuat JWT secara manual

func GenerateJWT(username string, role int, id primitive.ObjectID, idJenisUser int) (string, error) {
	// Header
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))

	// Konversi ObjectID menjadi string
	idStr := id.Hex()

	// Payload
	payload := map[string]interface{}{
		"username":      username,
		"role":          role,
		"id":            idStr, // Simpan ObjectID sebagai string
		"id_jenis_user": idJenisUser,
		"exp":           time.Now().Add(time.Hour * 1).Unix(),
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	payloadBase64 := base64.RawURLEncoding.EncodeToString(payloadBytes)

	// Signature (dummy signature in this example)
	signature := base64.RawURLEncoding.EncodeToString([]byte("dummy_signature"))

	// Combine to form token
	token := strings.Join([]string{header, payloadBase64, signature}, ".")
	return token, nil
}
