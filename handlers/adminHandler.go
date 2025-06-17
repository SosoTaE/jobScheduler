package handlers

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"jobScheduler/logger"
	"jobScheduler/models"
	"os"
)

func SeedAdminUser(db *gorm.DB, password string) {
	var user models.User
	// Check if a user with username "admin" already exists
	err := db.First(&user, "username = ?", "admin").Error
	if err == nil {
		// Admin user already exists
		logger.L.Info("Admin user already exists. Skipping seed.")
		return
	}

	logger.L.Info("Admin user not found, creating one...")

	// Hash the password from the environment variable
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.L.Error("FATAL: Could not hash admin password for seeding: %v", err)
		os.Exit(1)
	}

	// Create the new admin user
	adminUser := models.User{
		Username:     "admin",
		PasswordHash: string(hashedPassword),
		IsAdmin:      true,
	}

	// Save the user to the database
	if err := db.Create(&adminUser).Error; err != nil {
		logger.L.Error("FATAL: Could not seed admin user: %v", err)
		os.Exit(1)
	}

	logger.L.Info("Admin user created successfully.")
}
