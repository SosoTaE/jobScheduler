package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// ScheduleTime represents a specific time of day.
type ScheduleTime struct {
	Hour   int `json:"hour"`
	Minute int `json:"minute"`
}

// Schedule defines the structured schedule for a job.
// Schedule defines the structured schedule for a job.
type Schedule struct {
	Years       []int          `json:"years,omitempty"`
	Months      []int          `json:"months,omitempty"`
	DaysOfMonth []int          `json:"daysOfMonth,omitempty"`
	Weekdays    []time.Weekday `json:"weekdays,omitempty"`
	Times       []ScheduleTime `json:"times"`
}

func (s *Schedule) Validate() error {
	// --- THIS IS THE NEW VALIDATION LOGIC FOR THE YEAR ---
	// Get the current year once for comparison.
	currentYear := time.Now().Year()

	// Loop through the years provided by the user.
	for _, year := range s.Years {
		// If any year is in the past, reject the entire schedule.
		if year < currentYear {
			return fmt.Errorf("invalid year: %d is in the past and cannot be scheduled", year)
		}
	}
	// --- END YEAR VALIDATION ---

	// A schedule must have at least one time to run.
	if len(s.Times) == 0 {
		return errors.New("schedule must contain at least one execution time")
	}

	// Validate the time values.
	for _, t := range s.Times {
		if t.Hour < 0 || t.Hour > 23 {
			return fmt.Errorf("invalid hour value: %d, must be between 0-23", t.Hour)
		}
		if t.Minute < 0 || t.Minute > 59 {
			return fmt.Errorf("invalid minute value: %d, must be between 0-59", t.Minute)
		}
	}

	// ... (other validations like for months can go here) ...

	return nil
}

// GORM requires these two methods to handle custom types like jsonb.
func (s Schedule) Value() (driver.Value, error) {
	return json.Marshal(s)
}
func (s *Schedule) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &s)
}

// Job represents a scheduled task in the system.

type Job struct {
	gorm.Model
	Name      string     `json:"name" gorm:"not null"`
	Command   string     `json:"command" gorm:"not null"`
	Schedule  Schedule   `json:"schedule" gorm:"type:jsonb"`
	Status    string     `json:"status" gorm:"default:'pending'"`
	LastRunAt *time.Time `json:"lastRunAt,omitempty"`
	UserID    uint       `json:"userId"`
}
