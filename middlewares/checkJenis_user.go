package middlewares

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func CheckJenisUser(ju int) fiber.Handler {
	return func(c *fiber.Ctx) error {
		jenisUser := c.Locals("JenisUser") // Ambil jenis_user dari context
		if jenisUser != ju {
			return c.Status(http.StatusForbidden).JSON(fiber.Map{
				"error": fmt.Sprintf("Access denied: jenis user %d required", ju, "input %s", jenisUser),
			})
		}
		return c.Next()
	}
}
