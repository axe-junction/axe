package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	GoogleID      string         `json:"google_id" gorm:"uniqueIndex"`
	Email         string         `json:"email" gorm:"uniqueIndex;not null"`
	Name          string         `json:"name"`
	GivenName     string         `json:"given_name"`
	FamilyName    string         `json:"family_name"`
	Picture       string         `json:"picture"`
	Locale        string         `json:"locale"`
	VerifiedEmail bool           `json:"verified_email"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

type UserService interface {
	CreateOrUpdateUser(user *User) error
	GetUserByGoogleID(googleID string) (*User, error)
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id uint) (*User, error)
}
