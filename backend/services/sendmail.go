package services

import (
	"booksrendezvous-backend/models"
	"booksrendezvous-backend/utils"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/smtp"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// EmailService handles email-related operations
type EmailService struct {
	DB *gorm.DB
}

// Configuration variables
var emailConfig, _ = utils.LoadConfig()
var logger = utils.SugaredLogger

// NewEmailService creates a new email service instance
func NewEmailService(db *gorm.DB) *EmailService {
	return &EmailService{DB: db}
}

// SendPasswordResetEmail sends a password reset email to the provided address
func (s *EmailService) SendPasswordResetEmail(email string) error {
	// Check if the email exists in the database
	var user models.User
	result := s.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			logger.Warnw("Reset password attempt for non-existent email", "email", email)
			// Return success anyway to prevent email enumeration attacks
			return nil
		}
		logger.Errorw("Database error when finding user", "error", result.Error)
		return fmt.Errorf("database error: %w", result.Error)
	}

	// Generate a secure token
	token, err := generateSecureToken(32)
	if err != nil {
		logger.Errorw("Failed to generate secure token", "error", err)
		return fmt.Errorf("token generation failed: %w", err)
	}

	// Store the token in the database with an expiration time
	resetToken := models.PasswordResetToken{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour), // Token valid for 24 hours
		CreatedAt: time.Now(),
	}

	// Create or update token in database
	if err := s.DB.Transaction(func(tx *gorm.DB) error {
		// Delete any existing tokens for this user
		if err := tx.Where("user_id = ?", user.ID).Delete(&models.PasswordResetToken{}).Error; err != nil {
			return err
		}

		// Create new token
		if err := tx.Create(&resetToken).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		logger.Errorw("Failed to save reset token", "error", err)
		return fmt.Errorf("failed to save reset token: %w", err)
	}

	// Send the email
	if err := s.sendEmail(user.Email, user.Name, token); err != nil {
		logger.Errorw("Failed to send password reset email", "error", err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	logger.Infow("Password reset email sent successfully", "email", email)
	return nil
}

// VerifyPasswordResetToken verifies if a token is valid and not expired
func (s *EmailService) VerifyPasswordResetToken(token string) (string, error) {
	var resetToken models.PasswordResetToken
	result := s.DB.Where("token = ?", token).First(&resetToken)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return "", errors.New("invalid or expired token")
		}
		return "", fmt.Errorf("database error: %w", result.Error)
	}

	// Check if token has expired
	if time.Now().After(resetToken.ExpiresAt) {
		return "", errors.New("token has expired")
	}

	return resetToken.UserID, nil
}

// ResetPassword updates the user's password if the token is valid
func (s *EmailService) ResetPassword(token, newPassword string) error {
	// Verify token and get user ID
	userID, err := s.VerifyPasswordResetToken(token)
	if err != nil {
		return err
	}

	// Update password in the database
	return s.DB.Transaction(func(tx *gorm.DB) error {
		// Hash the new password
		hashedPassword, err := utils.HashPassword(newPassword)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}

		// Update user's password
		if err := tx.Model(&models.User{}).Where("id = ?", userID).Update("password", hashedPassword).Error; err != nil {
			return fmt.Errorf("failed to update password: %w", err)
		}

		// Delete the used token
		if err := tx.Where("token = ?", token).Delete(&models.PasswordResetToken{}).Error; err != nil {
			return fmt.Errorf("failed to delete token: %w", err)
		}

		return nil
	})
}

// Helper functions

// generateSecureToken creates a cryptographically secure token
func generateSecureToken(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// sendEmail sends the password reset email
func (s *EmailService) sendEmail(to, name, token string) error {
	// SMTP server configuration
	smtpHost := emailConfig.SMTPHost
	smtpPort := emailConfig.SMTPPort
	smtpUsername := emailConfig.SMTPUsername
	smtpPassword := emailConfig.SMTPPassword
	fromEmail := emailConfig.FromEmail

	// Reset link configuration
	frontendURL := emailConfig.FrontendURL
	resetLink := fmt.Sprintf("%s/password/reset?token=%s", frontendURL, token)

	// Construct email content
	subject := "Réinitialisation de votre mot de passe"
	body := fmt.Sprintf(`Bonjour %s,

Nous avons reçu une demande de réinitialisation de mot de passe pour votre compte. 
Pour réinitialiser votre mot de passe, cliquez sur le lien suivant :

%s

Ce lien expirera dans 24 heures.

Si vous n'avez pas demandé cette réinitialisation, vous pouvez ignorer cet email.

Cordialement,
L'équipe Bibliothèque
`, name, resetLink)

	msg := "From: " + fromEmail + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n\r\n" +
		body

	// Connect to SMTP server and send email
	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, fromEmail, []string{to}, []byte(msg))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
