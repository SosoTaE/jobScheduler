package models

import "gorm.io/gorm"

// User represents a user account in the system.
type User struct {
	gorm.Model
	Username     string `json:"username" gorm:"unique;not null"`
	PasswordHash string `json:"-" gorm:"not null"`
	IsAdmin      bool   `json:"isAdmin" gorm:"not null"`
	APIKey       string `gorm:"unique"`
}
