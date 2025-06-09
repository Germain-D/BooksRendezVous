package controllers

import (
	"booksrendezvous-backend/database"
	"booksrendezvous-backend/services"
	"booksrendezvous-backend/utils"

	"github.com/gofiber/fiber/v2"
)

var passwordLogger = utils.SugaredLogger

// ForgetPassword handles password reset request
func ForgetPassword(c *fiber.Ctx) error {
	passwordLogger.Info("Received forget password request")

	// Parse request body
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		passwordLogger.Errorw("Failed to parse request body", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}

	email := data["email"]
	if email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email is required",
		})
	}

	// Create email service
	emailService := services.NewEmailService(database.DB)

	// Send password reset email
	if err := emailService.SendPasswordResetEmail(email); err != nil {
		passwordLogger.Errorw("Failed to send password reset email", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to send password reset email",
		})
	}

	// Always return success to prevent email enumeration attacks
	return c.JSON(fiber.Map{
		"message": "If the email exists, a password reset link has been sent",
	})
}

// VerifyResetToken verifies the password reset token
func VerifyResetToken(c *fiber.Ctx) error {
	passwordLogger.Info("Verifying reset token")

	// Parse request body to get the token
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		passwordLogger.Errorw("Failed to parse request body", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"valid": false,
			"error": "Failed to parse request body",
		})
	}

	token := data["token"]
	if token == "" {
		sugar.Error("Token is required")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"valid": false,
			"error": "Token is required",
		})
	}

	emailService := services.NewEmailService(database.DB)
	userID, err := emailService.VerifyPasswordResetToken(token)
	if err != nil {
		passwordLogger.Infow("Invalid or expired token", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"valid": false,
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"valid":   true,
		"user_id": userID,
	})
}

// ResetPassword resets user password using a valid token
func ResetPassword(c *fiber.Ctx) error {
	passwordLogger.Info("Processing password reset")

	// Parse request body
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		passwordLogger.Errorw("Failed to parse request body", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}

	token := data["token"]
	password := data["password"]

	if token == "" || password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Token and password are required",
		})
	}

	// Validate password
	if len(password) < 10 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Password must be at least 10 characters long",
		})
	}

	emailService := services.NewEmailService(database.DB)
	if err := emailService.ResetPassword(token, password); err != nil {
		passwordLogger.Errorw("Failed to reset password", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Password has been reset successfully",
	})
}
