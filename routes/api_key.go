package routes

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"jobScheduler/handlers"
	"jobScheduler/logger"
	"jobScheduler/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// GenerateSecureKey creates a 32-byte (64 character) hex-encoded string
func GenerateSecureKey() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// GenerateAPIKey handles the request from the frontend
func GenerateAPIKey(db *gorm.DB) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		auth_ctx := ctx.Locals("auth_ctx").(handlers.AuthContext)

		key, err := GenerateSecureKey()
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error":   "Failed to generate API key",
			})
		}

		var user models.User

		if err := db.First(&user, auth_ctx.UserID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {

				errorMessage := fmt.Sprintf("Failed to find user with id %d: %v", auth_ctx.UserID, err)

				logger.L.Error(errorMessage)
				return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"success": false,
					"error":   "User not found",
				})
			}
			// For any other database error
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error":   "Database error: " + err.Error(),
			})
		}

		if err := db.Model(&user).Update("api_key", key).Error; err != nil {
			logger.L.Error("Failed to save API key to database", "user_id", auth_ctx.UserID, "error", err)
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error":   "Failed to save API key",
			})
		}

		return ctx.JSON(fiber.Map{
			"success": true,
			"api_key": key,
		})
	}
}
