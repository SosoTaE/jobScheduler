package handlers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"jobScheduler/models"
	"jobScheduler/structs"
)

func Login(db *gorm.DB, store *session.Store) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		req := new(structs.LoginRequest)
		if err := ctx.BodyParser(req); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
		}

		// --- 1. Find the user in the database ---
		var user models.User
		err := db.First(&user, "username = ?", req.Username).Error

		// Check if user was found
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid username or password"})
		} else if err != nil {
			// Handle other potential database errors
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
		if err != nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid username or password"})
		}

		sess, err := store.Get(ctx)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create session"})
		}

		sess.Set("is_authenticated", true)
		sess.Set("username", user.Username)
		sess.Set("user_id", user.ID) // Store the user's ID for future use
		sess.Set("is_admin", user.IsAdmin)

		if err := sess.Save(); err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not save session"})
		}

		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"message": "Logged in successfully",
		})
	}
}

type AuthContext struct {
	UserID   uint
	Username string
	IsAdmin  bool
}

// AuthRequired now only needs the session store.
func AuthRequired(store *session.Store) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		sess, err := store.Get(ctx)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error":   "Session error",
			})
		}

		// --- Check for authentication and extract all data from the session ---
		if auth, ok := sess.Get("is_authenticated").(bool); !ok || !auth {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error":   "Unauthorized",
			})
		}

		userID, ok1 := sess.Get("user_id").(uint)
		username, ok2 := sess.Get("username").(string)
		isAdmin, _ := sess.Get("is_admin").(bool)

		if !ok1 || !ok2 {
			err = sess.Destroy()
			if err != nil {
				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"success": false,
					"error":   "Session error",
				})
			}
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized: Invalid session data"})
		}

		ctx.Locals("auth_ctx", AuthContext{
			UserID:   userID,
			Username: username,
			IsAdmin:  isAdmin,
		})

		return ctx.Next()
	}
}

// Logout destroys the session.
func Logout(store *session.Store) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		sess, err := store.Get(ctx)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error":   "Could not get session",
			})
		}
		if err := sess.Destroy(); err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error":   "Could not destroy session",
			})
		}
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"message": "Logged out successfully",
		})
	}
}

// RegistrationRequest defines the structure for the registration request body.
type RegistrationRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Register is the handler for creating a new user account.
func Register(db *gorm.DB) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		auth_ctx := ctx.Locals("auth_ctx").(*AuthContext)
		if auth_ctx.IsAdmin != true {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error":   "Only admin can register users",
			})
		}

		req := new(RegistrationRequest)
		if err := ctx.BodyParser(req); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error":   "Cannot parse JSON",
			})
		}

		// 2. Validate input fields
		if len(req.Username) < 3 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error":   "Username must be at least 3 characters long",
			})
		}
		if len(req.Password) < 8 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error":   "Password must be at least 8 characters long",
			})
		}

		var existingUser models.User

		err := db.First(&existingUser, "username = ?", req.Username).Error
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			if err == nil {
				return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{
					"success": false,
					"error":   "Username already exists",
				})
			}
			// For any other unexpected database error
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error":   "Database error",
			})
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error":   "Failed to hash password",
			})
		}

		newUser := models.User{
			Username:     req.Username,
			PasswordHash: string(hashedPassword),
		}

		if result := db.Create(&newUser); result.Error != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error":   "Could not create user",
			})
		}

		return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
			"success": true,
			"message": "User registered successfully. Please log in.",
		})
	}
}
