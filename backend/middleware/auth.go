package middleware

import (
	"booksrendezvous-backend/utils"
	"log"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

// Protected protect routes
func Protected() fiber.Handler {
	config, err := utils.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	return jwtware.New(jwtware.Config{
		SigningKey:   jwtware.SigningKey{Key: []byte(config.SecretKey)},
		ErrorHandler: jwtError,
	})
}

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"status": "error", "message": "Missing or malformed JWT", "data": nil})
	}
	return c.Status(fiber.StatusUnauthorized).
		JSON(fiber.Map{"status": "error", "message": "Invalid or expired JWT", "data": nil})
}
