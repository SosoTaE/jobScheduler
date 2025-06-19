package routes

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"jobScheduler/logger"
	"jobScheduler/models"
)

func UpdateJob(db *gorm.DB) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		id := ctx.QueryInt("id")

		var existingJob models.Job

		if err := db.First(&existingJob, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {

				errorMessage := fmt.Sprintf("Failed to find job with id %s: %v", id, err)

				logger.L.Error(errorMessage)
				return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"success": false,
					"error":   "Job not found",
				})
			}
			// For any other database error
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error":   "Database error: " + err.Error(),
			})
		}

		var updatedData models.Job
		if err := ctx.BodyParser(&updatedData); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error":   "Cannot parse JSON: " + err.Error(),
			})
		}

		result := db.Model(&existingJob).Updates(updatedData)
		if result.Error != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error":   "Failed to update job: " + result.Error.Error(),
			})
		}

		message := fmt.Sprintf("updated the job with id: %d", existingJob.ID)
		logger.L.Info(message)

		return ctx.Status(fiber.StatusOK).JSON(
			fiber.Map{
				"success": true,
				"data":    existingJob,
			})
	}
}
