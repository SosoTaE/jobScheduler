package models

import (
	"time"

	"gorm.io/gorm"
)

// Job represents a scheduled task in the system.

type Job struct {
	gorm.Model
	Name      string     `json:"name" gorm:"not null"`
	Command   string     `json:"command" gorm:"not null"`
	Schedule  string     `json:"schedule" gorm:"not null"`
	Status    string     `json:"status" gorm:"default:'pending'"`
	LastRunAt *time.Time `json:"lastRunAt,omitempty"`
	UserID    uint       `json:"userId"`
}
