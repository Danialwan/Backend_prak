package middlewares

import (
	"fmt"
	"project-crud/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(c *fiber.Ctx) error {
	// Ambil token dari header Authorization
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing or invalid token"})
	}

	// Ambil tokennya setelah "Bearer "

	// Decode token
	token := strings.TrimPrefix(authHeader, "Bearer ")
	payload, err := utils.DecodeJWT(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}

	// Debugging payload
	fmt.Println("Decoded payload:", payload)

	// Validasi payload
	username, ok := payload["username"].(string)
	if !ok || username == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token payload"})
	}
	// Validasi token
	// username, err := utils.ValidateJWT(token)
	// if err != nil {
	// 	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	// }

	// Set username ke context untuk dipakai di handler berikutnya
	c.Locals("username", username)
	return c.Next()
}
