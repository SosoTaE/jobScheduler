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

func CreateJob(db *gorm.DB) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		autx_ctx := ctx.Locals("auth_ctx").(handlers.AuthContext)
		newJob := new(models.Job)

		if err := ctx.BodyParser(newJob); err != nil {
			logger.L.Error(err.Error())
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error":   err.Error(),
			})
		}

		if newJob.Name == "" || newJob.Command == "" {
			logger.L.Error("Missing required fields: name, command")
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error":   "Missing required fields: name, command",
			})
		}

		if err := newJob.Schedule.Validate(); err != nil {
			logger.L.Error(err.Error())
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error":   err.Error(),
			})
		}

		newJob.Status = "pending"
		newJob.CreatedAt = time.Now()
		newJob.UserID = autx_ctx.UserID

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
}
