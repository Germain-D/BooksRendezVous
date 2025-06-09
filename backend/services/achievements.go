package services

import (
	"booksrendezvous-backend/models"
	"booksrendezvous-backend/utils"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

var config, _ = utils.LoadConfig()

// Initialize logger with the log level from environment variables
var sugar = utils.SugaredLogger

type AchievementService struct {
	DB *gorm.DB
}

func NewAchievementService(db *gorm.DB) *AchievementService {
	return &AchievementService{DB: db}
}

// Déclencheur principal pour vérifier les succès
func (s *AchievementService) CheckAchievements(userID string) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		// Récupération des stats mises à jour
		var stat models.UserStat
		if err := tx.First(&stat, "user_id = ?", userID).Error; err != nil {
			sugar.Errorw("User stats not found", "userID", userID)
			return fmt.Errorf("user stats not found: %w", err)
		}

		// Récupération de tous les succès activables par les stats
		var achievements []models.Achievement
		if err := tx.Find(&achievements).Error; err != nil {
			sugar.Errorw("Failed to get achievements", "error", err)
			return fmt.Errorf("failed to get achievements: %w", err)
		}

		for _, achievement := range achievements {
			sugar.Infow("Processing achievement", "name", achievement.Name)
			currentValue, ok := getStatValue(&stat, achievement.TargetStat)
			if !ok {
				sugar.Infow("Skipping achievement", "name", achievement.Name)
				continue
			}

			sugar.Infow("Achievement progress", "currentValue", currentValue, "targetValue", achievement.TargetValue)

			err := s.processAchievement(tx, userID, achievement, currentValue)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// Méthode interne de traitement par succès
func (s *AchievementService) processAchievement(tx *gorm.DB, userID string, achievement models.Achievement, currentValue int) error {
	var userAch models.UserAchievement
	err := tx.Where("user_id = ? AND achievement_id = ?", userID, achievement.ID).
		First(&userAch).
		Error

	progress := calculateProgress(achievement, currentValue)
	unlocked := progress >= achievement.TargetValue

	sugar.Infow("Achievement progress", "progress", progress, "unlocked", unlocked)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		sugar.Infow("New achievement unlocked", "userID", userID, "achievement", achievement.Name)
		return s.createNewAchievement(tx, userID, achievement, progress, unlocked)
	}
	if err != nil {
		sugar.Errorw("Database error", "error", err)
		return fmt.Errorf("database error: %w", err)
	}

	return s.updateExistingAchievement(tx, &userAch, progress, unlocked)
}

// Gestion de la création des nouveaux UserAchievement
func (s *AchievementService) createNewAchievement(tx *gorm.DB, userID string, achievement models.Achievement, progress int, unlocked bool) error {
	userAch := models.UserAchievement{
		UserID:        userID,
		AchievementID: achievement.ID,
		Progress:      progress,
	}

	if unlocked {
		userAch.UnlockedAt = time.Now()
	}

	if err := tx.Create(&userAch).Error; err != nil {
		sugar.Errorw("Failed to create achievement", "error", err)
		return fmt.Errorf("failed to create achievement: %w", err)
	}

	if unlocked {
		go s.sendUnlockedNotification(userID, achievement.ID)
	}
	return nil
}

// Gestion des mises à jour des UserAchievement existants
func (s *AchievementService) updateExistingAchievement(tx *gorm.DB, userAch *models.UserAchievement, progress int, unlocked bool) error {
	updates := make(map[string]interface{})
	needsUpdate := false

	// Mise à jour progression si amélioration
	if progress > userAch.Progress && !userAch.IsUnlocked() {
		updates["progress"] = progress
		needsUpdate = true
	}

	// Déblocage si atteint
	if unlocked && !userAch.IsUnlocked() {
		updates["unlocked_at"] = time.Now()
		updates["notified"] = false
		needsUpdate = true
	}

	if needsUpdate {
		if err := tx.Model(userAch).Updates(updates).Error; err != nil {
			return fmt.Errorf("failed to update achievement: %w", err)
		}

		if unlocked {
			go s.sendUnlockedNotification(userAch.UserID, userAch.AchievementID)
		}
	}
	return nil
}

// Helper pour extraire les valeurs des stats
func getStatValue(stat *models.UserStat, field string) (int, bool) {
	sugar.Infow("Getting stat value", "field", field)
	switch field {
	case "TotalBooks":
		return stat.TotalBooks, true
	case "CompletedBooks":
		return stat.CompletedBooks, true
	case "ReadingBooks":
		return stat.ReadingBooks, true
	case "TotalPages":
		return stat.TotalPages, true
	case "FavoriteBooks":
		return stat.FavoriteBooks, true
	default:
		return 0, false // Ignore les succès non liés aux stats
	}
}

// Logique de calcul de progression selon le type de succès
func calculateProgress(a models.Achievement, currentValue int) int {
	switch a.Type {
	case models.TypeCounter:
		return currentValue // Compteur incrémental
	case models.TypeMilestone, models.TypeBadge:
		if currentValue >= a.TargetValue {
			return a.TargetValue
		}
		return currentValue // Progression partielle
	default:
		return 0
	}
}

// Système de notification (exemple)
func (s *AchievementService) sendUnlockedNotification(userID, achievementID string) {
	// Implémentation de la notification (webhook, email, websocket, etc)
	fmt.Printf("Achievement unlocked! User: %s, Achievement: %s\n", userID, achievementID)
}
