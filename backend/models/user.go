// models/user.go
package models

type User struct {
	ID       string `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name     string `gorm:"not null" json:"name"`
	Email    string `gorm:"unique; not null" json:"email"`
	Password []byte `gorm:"not null" json:"-"`
}

type Publicusers struct {
	UserID   string `gorm:"type:uuid;references:User" json:"userId"`
	IsPublic bool   `gorm:"type:boolean;not null" json:"isPublic"`
	PublicID string `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"publicID"`

	User User `gorm:"foreignKey:UserID" json:"-"`
}
