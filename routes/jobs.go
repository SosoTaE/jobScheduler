package routes

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"jobScheduler/logger"
	"jobScheduler/models"
	"jobScheduler/structs"
	"math"
)

func ListJobs(ctx *fiber.Ctx, db *gorm.DB) error {
	// --- 1. Set up pagination parameters ---
	page := ctx.QueryInt("page", 1)
	if page < 1 {
		page = 1
	}

	limit := ctx.QueryInt("limit", 10)
	if limit < 1 {
		limit = 10
	}
	// Set a max limit to prevent abuse
	if limit > 100 {
		limit = 100
	}

	// Calculate the offset for the database query
	offset := (page - 1) * limit

	// --- 2. Fetch data from the database ---
	var jobs []models.Job
	var totalCount int64

	// First query: Get the total number of jobs
	if err := db.Model(&models.Job{}).Count(&totalCount).Error; err != nil {
		logger.L.Error("Failed to count jobs")
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Database error",
		})
	}

	if err := db.Order("created_at desc").Offset(offset).Limit(limit).Find(&jobs).Error; err != nil {
		errorMessage := fmt.Sprintf("Failed to get jobs error: %s", err.Error())
		logger.L.Error(errorMessage)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   errorMessage,
		})
	}

	// --- 3. Assemble the response ---
	totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))

	response := structs.PaginatedJobsResponse{
		Data: jobs,
		Meta: structs.PaginationMeta{
			TotalRecords: totalCount,
			TotalPages:   totalPages,
			CurrentPage:  page,
			PageSize:     limit,
		},
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    response,
	})
}
