package controllers

import (
	"booksrendezvous-backend/database"
	"booksrendezvous-backend/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type PublicController struct {
	PublicID string `json:"publicid"`
}

func GetPublicUser(c *fiber.Ctx) error {
	sugar.Info("Get public user")

	var request PublicController
	if err := c.BodyParser(&request); err != nil {
		sugar.Errorw("Failed to parse request body", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}

	// Check if user exists
	var publicusers models.Publicusers
	database.DB.Where("public_id = ?", request.PublicID).First(&publicusers)
	if publicusers.PublicID == "" {
		sugar.Warnw("Login attempt failed: user not found")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid credentials",
		})
	}

	if !publicusers.IsPublic {
		sugar.Warnw("Login attempt failed: user not public")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid credentials",
		})
	}

	// retrieve books from public user
	var books []models.Book
	database.DB.Where("user_id = ?", publicusers.UserID).Find(&books)

	return c.JSON(fiber.Map{
		"books": books,
	})
}

func ChangePublicVisibility(c *fiber.Ctx) error {
	sugar.Info("Change public visibility")

	// Check if user is authorized
	uiidStr, ok := CheckAuth(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// convert string to uuid type
	userID, err := uuid.Parse(uiidStr)
	if err != nil {
		sugar.Errorw("Failed to parse uuid", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse uuid",
		})
	}

	// Check if user exists
	var publicusers models.Publicusers
	database.DB.Where("user_id = ?", userID).First(&publicusers)
	if publicusers.UserID == "" {
		sugar.Warnw("Login attempt failed: user not found")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid credentials",
		})
	}
	publicusers.IsPublic = !publicusers.IsPublic

	database.DB.Where("user_id = ?", userID).Save(&publicusers)

	return c.JSON(fiber.Map{
		"message": "Public visibility changed successfully",
	})
}

func GetPublicVisibility(c *fiber.Ctx) error {
	sugar.Info("Change public visibility")

	// Check if user is authorized
	uiidStr, ok := CheckAuth(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// convert string to uuid type
	userID, err := uuid.Parse(uiidStr)
	if err != nil {
		sugar.Errorw("Failed to parse uuid", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse uuid",
		})
	}

	// Check if user exists
	var publicusers models.Publicusers
	database.DB.Where("user_id = ?", userID).First(&publicusers)
	if publicusers.UserID == "" {
		sugar.Warnw("Login attempt failed: user not found")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid credentials",
		})
	}

	return c.JSON(fiber.Map{
		"public":    publicusers.IsPublic,
		"sharelink": publicusers.PublicID,
	})
}
