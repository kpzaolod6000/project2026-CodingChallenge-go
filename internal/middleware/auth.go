package middleware

import (
	"crypto/subtle"
	"os"

	"github.com/gofiber/fiber/v2"
)

func KeyAuth() fiber.Handler {
	expectedKey := os.Getenv("API_KEY")
	if expectedKey == "" {
		expectedKey = "coding-challenge"
	}

	return func(c *fiber.Ctx) error {
		clientKey := c.Get("X-API-Key")
		if clientKey == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "falta la cabecera X-API-Key",
			})
		}

		if subtle.ConstantTimeCompare([]byte(clientKey), []byte(expectedKey)) != 1 {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "clave de API inválida",
			})
		}

		return c.Next()
	}
}