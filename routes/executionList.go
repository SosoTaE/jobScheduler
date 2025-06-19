package routes

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"jobScheduler/models"
	"jobScheduler/structs"
	"math"
	"strconv"
)

func ListAllExecutions(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		page, _ := strconv.Atoi(c.Query("page", "1"))
		limit, _ := strconv.Atoi(c.Query("limit", "10"))
		if page < 1 {
			page = 1
		}
		if limit < 1 || limit > 100 {
			limit = 100
		}
		offset := (page - 1) * limit

		var executions []models.JobExecution
		var totalCount int64

		db.Model(&models.JobExecution{}).Count(&totalCount)
		db.Order("created_at desc").Offset(offset).Limit(limit).Find(&executions)

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
