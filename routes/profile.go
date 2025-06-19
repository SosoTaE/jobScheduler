package routes

import (
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

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"data": fiber.Map{
				"username": authCtx.Username,
			},
		})
	}
}
