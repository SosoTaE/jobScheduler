package routes

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"jobScheduler/handlers"
)

// Profile is a simple handler to demonstrate an authenticated route.
func Profile() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authCtx, ok := c.Locals("auth_ctx").(handlers.AuthContext)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error":   "Unauthorized",
			})
		}

		message := fmt.Sprintf("Hello, %s! Your user ID is %d.", authCtx.Username, authCtx.UserID)
		if authCtx.IsAdmin {
			message += " You have admin privileges."
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"message": message,
		})
	}
}
