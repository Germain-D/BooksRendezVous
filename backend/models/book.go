package models

import "github.com/lib/pq"

type Book struct {
	ID            string         `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID        string         `gorm:"type:uuid;index;references:User" json:"userId"`   // Relation avec la table User
	GoogleBooksID string         `gorm:"size:512" json:"googleBooksId"`                   // Taille augmentée pour les ID complexes
	Status        string         `gorm:"size:128;default:'pending'" json:"status"`        // Valeur par défaut
	Rating        int            `gorm:"check:rating >= 0 AND rating <= 5" json:"rating"` // Validation de range
	Comment       string         `gorm:"type:text" json:"comment"`
	Title         string         `gorm:"size:255;not null" json:"title"`
	ImageUrl      string         `gorm:"type:text" json:"imageUrl"`
	Authors       pq.StringArray `gorm:"type:text[];not null" json:"authors"` // Utilisation de pq.StringArray
	Description   string         `gorm:"type:text" json:"description"`
	Favorite      bool           `gorm:"default:false" json:"favorite"` // Valeur par défaut
	PageCount     int            `gorm:"default:0" json:"pageCount"`    // Valeur par défaut
	Genres        pq.StringArray `gorm:"type:text[]" json:"genres"`     // Utilisation de pq.StringArray
	PublishedDate string         `gorm:"size:128" json:"publishedDate"`

	// Relation supplémentaire si nécessaire
	User User `gorm:"foreignKey:UserID" json:"-"`
}
