package models

import "time"

type AchievementType string

const (
	TypeCounter   AchievementType = "counter"   // Ex: "X livres lus"
	TypeMilestone AchievementType = "milestone" // Ex: "1000 pages"
	TypeBadge     AchievementType = "badge"     // Ex: "Premier livre"
)

type Achievement struct {
	ID          string          `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name        string          `gorm:"type:varchar(255);unique;not null"`
	Description string          `gorm:"type:text"`
	Type        AchievementType `gorm:"type:varchar(30);index;not null"`
	TargetValue int
	TargetStat  string `gorm:"type:varchar(50)"` // Add this field
	IsHidden    bool   `gorm:"default:false"`
	Category    string `gorm:"type:varchar(50);index"`
	Image       string `gorm:"type:varchar(255)"`

	Users []UserAchievement `gorm:"foreignKey:AchievementID"`
}

type UserAchievement struct {
	UserID        string    `gorm:"type:uuid;primaryKey"`
	AchievementID string    `gorm:"type:uuid;primaryKey"`
	Progress      int       `gorm:"default:0;not null"`
	UnlockedAt    time.Time `gorm:"index;not null"`
	Notified      bool      `gorm:"default:false"` // Si notification envoyée

	Achievement Achievement `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	User        User        `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// Add IsUnlocked method
func (ua *UserAchievement) IsUnlocked() bool {
	return !ua.UnlockedAt.IsZero()
}
