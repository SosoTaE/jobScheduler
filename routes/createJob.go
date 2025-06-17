package routes

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"jobScheduler/handlers"
	"jobScheduler/logger"
	"jobScheduler/models"
	"time"
)

func CreateJob(ctx *fiber.Ctx, db *gorm.DB) error {
	user, ok := ctx.Locals("auth_ctx").(handlers.AuthContext)

	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to parse auth context",
		})
	}

	if user.IsAdmin != true {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "User is not admin",
		})
	}

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
	newJob.CreatedAt = time.Now()
	newJob.UserID = user.UserID

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
