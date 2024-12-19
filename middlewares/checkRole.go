package middlewares

import (
	"fmt"
	"strings"

	"project-crud/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CheckRole(requiredRole int) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Ambil token dari header Authorization
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Authorization header missing"})
		}

		// Pastikan format token valid
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token format"})
		}

		// Decode token
		token := tokenParts[1]
		payload, err := utils.DecodeJWT(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
		}

		// Log isi payload untuk debugging
		fmt.Printf("Payload: %+v\n", payload)

		// Validasi role
		role, ok := payload["role"].(float64)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid role in token"})
		}
		if int(role) != requiredRole {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Access denied"})
		}

		// Validasi id
		idStr, ok := payload["id"].(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid id in token"})
		}

		// Convert string id back to ObjectID
		id, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid ObjectID format"})
		}

		// Validasi id_jenis_user
		idJenisUser, ok := payload["id_jenis_user"].(float64)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid id_jenis_user in token"})
		}

		// Log id dan id_jenis_user
		fmt.Printf("ID: %s, ID Jenis User: %d\n", id.Hex(), int(idJenisUser))

		return c.Next()
	}
}
