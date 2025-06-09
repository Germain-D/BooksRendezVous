package models

type UserStat struct {
	UserID         string  `gorm:"type:uuid;primaryKey"`
	TotalBooks     int     `gorm:"default:0;not null"` // Tous statuts
	CompletedBooks int     `gorm:"default:0;not null"` // Status completed
	ToReadBooks    int     `gorm:"default:0;not null"` // Status toread
	ReadingBooks   int     `gorm:"default:0;not null"` // Status reading
	FavoriteBooks  int     `gorm:"default:0;not null"`
	TotalPages     int     `gorm:"default:0;not null"`
	AverageRating  float64 `gorm:"type:decimal(3,2)"`

	User User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
