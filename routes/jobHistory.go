package routes

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"jobScheduler/models"
	"jobScheduler/structs"
	"math"
	"strconv"
)

// ListJobHistory now allows any authenticated user to see any job's history.
func ListJobHistory(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		jobID := c.Params("id")

		// 2. SIMPLIFIED LOGIC: Check only if the job exists, regardless of owner.
		var job models.Job
		if err := db.First(&job, jobID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"success": false,
					"error":   "Job not found",
				})
			}
			// Handle any other potential database errors
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error":   "Database error",
			})
		}

		page, _ := strconv.Atoi(c.Query("page", "1"))
		limit, _ := strconv.Atoi(c.Query("limit", "10"))
		if page < 1 {
			page = 1
		}
		if limit < 1 || limit > 100 {
			limit = 10
		}
		offset := (page - 1) * limit

		var executions []models.JobExecution
		var totalCount int64

		db.Model(&models.JobExecution{}).Where("job_id = ?", jobID).Count(&totalCount)
		db.Order("created_at desc").Where("job_id = ?", jobID).Offset(offset).Limit(limit).Find(&executions)

		totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"data":    executions,
			"meta": structs.PaginationMeta{
				TotalRecords: totalCount,
				TotalPages:   totalPages,
				CurrentPage:  page,
				PageSize:     limit,
			},
		})
	}
}
