package routes

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"jobScheduler/logger"
	"jobScheduler/models"
)

func CreateJob(ctx *fiber.Ctx, db *gorm.DB) error {
	newJob := new(models.Job)

	if err := ctx.BodyParser(newJob); err != nil {
		logger.L.Error(err.Error())
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	if newJob.Name == "" || newJob.Command == "" || newJob.Schedule == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Missing required fields: name, command, and schedule",
		})
	}

	newJob.Status = "pending"

	result := db.Create(&newJob)
	if result.Error != nil {
		logger.L.Error(result.Error.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to save job: " + result.Error.Error(),
		})
	}

	message := fmt.Sprintf("created new job with id: %d", newJob.ID)
	logger.L.Info(message)

	return ctx.Status(fiber.StatusCreated).JSON(
		fiber.Map{
			"success": true,
			"data":    newJob,
		})
}
