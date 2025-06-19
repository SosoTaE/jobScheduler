package routes

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"jobScheduler/models"
)

func ListUsers(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var users []models.User
		if err := db.Find(&users).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error":   "Database error while fetching users",
			})
		}

		// 4. Return the list of users.
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"data":    users,
		})
	}
}
