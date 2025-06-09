package models

import "time"

type PasswordResetToken struct {
	ID        string    `gorm:"primarykey"`
	UserID    string    `gorm:"not null"`
	Token     string    `gorm:"not null;uniqueIndex"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time
	User      User `gorm:"foreignKey:UserID"`
}
