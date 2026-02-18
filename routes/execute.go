package routes

import (
	"fmt"
	"jobScheduler/handlers"
	"jobScheduler/logger"
	"jobScheduler/models"
	"jobScheduler/worker"

	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Execute(db *gorm.DB) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		auth_ctx := ctx.Locals("auth_ctx").(handlers.AuthContext)

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

		newJob.Status = "running"
		newJob.UserID = auth_ctx.UserID

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

		output, err := worker.ExecuteCommand(newJob.Command)
		executionStatus := "succeeded"
		if err != nil {
			executionStatus = "failed"
			logger.L.Error("Job execution failed", "job_id", newJob.ID, "error", err, "output", output)
		} else {
			logger.L.Info("Job execution succeeded", "job_id", newJob.ID, "output", output)
		}

		db.Model(&newJob).Update("status", executionStatus)

		// Create the detailed execution record
		executionRecord := models.JobExecution{
			JobID:      newJob.ID,
			Status:     executionStatus,
			Output:     output,
			FinishedAt: time.Now(),
		}
		if result := db.Create(&executionRecord); result.Error != nil {
			logger.L.Error("Failed to save job execution history", "job_id", newJob.ID, "error", result.Error)
		}

		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"job":     newJob,
		})

	}
}
