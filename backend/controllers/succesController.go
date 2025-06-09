package controllers

import (
	"booksrendezvous-backend/database"
	"booksrendezvous-backend/models"
	"booksrendezvous-backend/services"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SeedBaseAchievements(db *gorm.DB, jsonPath string) error {
	// Lecture du fichier JSON
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return fmt.Errorf("erreur lecture fichier JSON: %w", err)
	}

	// Structure temporaire pour le décodage
	var payload struct {
		Achievements []struct {
			Name        string `json:"name"`
			Description string `json:"description"`
			Type        string `json:"type"`
			TargetValue int    `json:"targetValue"`
			TargetSat   string `json:"targetStat"`
			IsHidden    bool   `json:"isHidden"`
			Category    string `json:"category"`
		} `json:"achievements"`
	}

	if err := json.Unmarshal(data, &payload); err != nil {
		return fmt.Errorf("erreur décodage JSON: %w", err)
	}

	// Transaction batch
	return db.Transaction(func(tx *gorm.DB) error {
		for _, a := range payload.Achievements {
			achievement := models.Achievement{
				Name:        a.Name,
				Description: a.Description,
				Type:        models.AchievementType(a.Type),
				TargetValue: a.TargetValue,
				TargetStat:  a.TargetSat,
				IsHidden:    a.IsHidden,
				Category:    a.Category,
			}

			// Insertion avec conflit contrôlé
			result := tx.Clauses(
				clause.OnConflict{
					Columns: []clause.Column{{Name: "name"}},
					DoUpdates: clause.Assignments(map[string]interface{}{
						"description":  achievement.Description,
						"type":         achievement.Type,
						"target_value": achievement.TargetValue,
						"target_stat":  achievement.TargetStat,
						"is_hidden":    achievement.IsHidden,
						"category":     achievement.Category,
					}),
				},
			).Create(&achievement)

			if result.Error != nil {
				return fmt.Errorf("erreur sur le succès '%s': %w", a.Name, result.Error)
			}

			if result.RowsAffected > 0 {
				fmt.Printf("Succès créé: %s\n", a.Name)
			} else {
				fmt.Printf("Succès existant mis à jour: %s\n", a.Name)
			}
		}
		return nil
	})
}

func GetAllAchievements(c *fiber.Ctx) error {
	var achievements []models.Achievement
	if err := database.DB.Find(&achievements).Error; err != nil {
		sugar.Errorw("Failed to get achievements", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get achievements",
		})
	}
	return c.JSON(achievements)
}

func GetAchievements(c *fiber.Ctx) error {
	sugar.Info("Received a user achievements request")

	// Get and validate UUID from auth
	uuidStr, ok := CheckAuth(c)
	if !ok {
		sugar.Warn("Unauthorized access attempt")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized access",
		})
	}

	userID, err := uuid.Parse(uuidStr)
	if err != nil {
		sugar.Errorw("Invalid UUID format", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID format",
		})
	}

	// Start transaction
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Check achievements in transaction
	service := services.NewAchievementService(tx)
	if err := service.CheckAchievements(uuidStr); err != nil {
		tx.Rollback()
		sugar.Errorw("Achievement check failed", "error", err, "userID", userID)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to process achievements",
		})
	}

	// Get all achievements and user achievements in single query
	var achievements []models.Achievement
	if err := tx.Find(&achievements).Error; err != nil {
		tx.Rollback()
		sugar.Errorw("Failed to fetch achievements", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Unable to retrieve achievements",
		})
	}

	var userAchievements []models.UserAchievement
	if err := tx.Where("user_id = ?", userID).Find(&userAchievements).Error; err != nil {
		tx.Rollback()
		sugar.Errorw("Failed to fetch user achievements", "error", err, "userID", userID)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Unable to retrieve user achievements",
		})
	}

	// Create achievement lookup map
	unlockedAchievements := make(map[string]time.Time)
	for _, ua := range userAchievements {
		unlockedAchievements[ua.AchievementID] = ua.UnlockedAt
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		sugar.Errorw("Transaction commit failed", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	// Create response
	type SimplifiedAchievement struct {
		Name        string     `json:"name"`
		Image       string     `json:"image"`
		Description string     `json:"description"`
		UnlockedAt  *time.Time `json:"unlockedAt,omitempty"`
	}

	response := make([]SimplifiedAchievement, 0, len(achievements))
	for _, achievement := range achievements {
		// Safely check if achievement is unlocked
		unlockedTime, isUnlocked := unlockedAchievements[achievement.ID]

		// Skip hidden achievements that aren't unlocked
		if achievement.IsHidden && (!isUnlocked || unlockedTime.IsZero()) {
			continue
		}

		sa := SimplifiedAchievement{
			Name:        achievement.Name,
			Image:       achievement.Image,
			Description: achievement.Description,
		}

		// Set unlock time if achievement is unlocked
		if isUnlocked && !unlockedTime.IsZero() {
			sa.UnlockedAt = &unlockedTime
		}

		response = append(response, sa)
	}

	sugar.Infow("User achievements retrieved successfully", "userID", userID)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"achievements": response,
	})
}
