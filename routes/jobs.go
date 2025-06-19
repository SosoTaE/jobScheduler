package routes

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"jobScheduler/models"
	"jobScheduler/structs"
	"math"
	"strconv"
)

// ListJobs retrieves a paginated list of jobs with flexible filtering.
func ListJobs(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		page, _ := strconv.Atoi(c.Query("page", "1"))
		if page < 1 {
			page = 1
		}

		limit, _ := strconv.Atoi(c.Query("limit", "10"))
		if limit < 1 {
			limit = 10
		}
		if limit > 100 {
			limit = 100
		}

		offset := (page - 1) * limit

		// 3. Fetch data from the database based on user role and query params.
		var jobs []models.Job
		var totalCount int64

		// Create a base query chain.
		query := db.Model(&models.Job{})

		// --- NEW LOGIC: Check for a userID filter in the query ---
		filterUserID := c.Query("userID")

		if filterUserID != "" {
			query = query.Where("user_id = ?", filterUserID)
		}

		// Run the count query on the (potentially filtered) data.
		if err := query.Count(&totalCount).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error while counting jobs"})
		}

		// Run the find query with pagination and ordering.
		if err := query.Order("created_at desc").Offset(offset).Limit(limit).Find(&jobs).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error while fetching jobs"})
		}

		// 4. Assemble and return the paginated response.
		totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"data":    jobs,
			"meta": structs.PaginationMeta{
				TotalRecords: totalCount,
				TotalPages:   totalPages,
				CurrentPage:  page,
				PageSize:     limit,
			},
		})
	}
}
