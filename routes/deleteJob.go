package routes

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"jobScheduler/handlers"
	"jobScheduler/logger"
	"jobScheduler/models"
)

func DeleteJob(ctx *fiber.Ctx, db *gorm.DB) error {
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

	id := ctx.QueryInt("id")

	if id == 0 {
		logger.L.Error("DeleteJob: id is invalid")
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "id is required",
		})
	}

	result := db.Delete(&models.Job{}, id)

	// Check for any database errors during the delete operation.
	if result.Error != nil {
		errorMessage := fmt.Sprintf("Delete job with id: %d", result.Error.Error())
		logger.L.Error(errorMessage)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   fmt.Sprintf("Failed to delete job: with id: %d %s", id, result.Error.Error()),
		})
	}

	// Check if a record was actually found and deleted.
	if result.RowsAffected == 0 {
		errorMessage := fmt.Sprintf("job not found with id: %d", result.Error)
		logger.L.Error(errorMessage)
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"sucess": false,
			"error":  "Job not found",
		})
	}

	message := fmt.Sprintf("deleted the job with id: %d", id)
	logger.L.Info(message)

	// Respond with a success message.
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Job successfully deleted",
	})
}
