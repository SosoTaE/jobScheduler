package routes

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"jobScheduler/models"
)

func GetJobDetails(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		jobID := c.Params("id")

		var job models.Job
		if err := db.First(&job, jobID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"success": true,
					"error":   "Job not found",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error":   "Database error",
			})
		}

		// 4. Return the Job object directly.
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"data":    job,
		})
	}
}
