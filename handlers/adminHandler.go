package handlers

import (
	"jobScheduler/config"
	"jobScheduler/logger"
	"jobScheduler/models"
	"os"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedAdminUser(db *gorm.DB, adminCredential config.AdminCredential) {
	var user models.User
	err := db.First(&user, "username = ?", adminCredential.Username).Error
	if err == nil {
		// Admin user already exists
		logger.L.Info("Admin user already exists. Skipping seed.")
		return
	}

	logger.L.Info("Admin user not found, creating one...")

	// Hash the password from the environment variable
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminCredential.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.L.Error("FATAL: Could not hash admin password for seeding: %v", err)
		os.Exit(1)
	}

	// Create the new admin user
	adminUser := models.User{
		Username:     adminCredential.Username,
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
