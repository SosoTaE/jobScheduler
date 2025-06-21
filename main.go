package main

import (
	"github.com/glebarez/sqlite" // Pure Go SQLite driver
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/session"
	"gorm.io/gorm"
	"io"
	"jobScheduler/config"
	"jobScheduler/handlers"
	"jobScheduler/logger"
	"jobScheduler/models"
	"jobScheduler/routes"
	"jobScheduler/worker"
	"log"
	"os"
	"time"
)

func main() {

	file, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		logger.L.Error(err.Error())
		os.Exit(1)
	}
	defer func() {
		if err := file.Close(); err != nil {
			logger.L.Error(err.Error())
		}
	}()

	multiWriter := io.MultiWriter(os.Stderr, file)

	log.SetOutput(multiWriter)

	db, err := gorm.Open(sqlite.Open("dispatch.db"), &gorm.Config{})
	if err != nil {
		logger.L.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}

	logger.L.Info("Database connection successful using SQLite.")

	err = db.AutoMigrate(&models.Job{}, &models.User{}, &models.JobExecution{})
	if err != nil {
		logger.L.Error("Failed to migrate tables", err)
		os.Exit(1)
	}

	adminCredential, err := config.GetAdminCredential()
	if err != nil {
		logger.L.Error("Failed to get admin credential", err)
		os.Exit(1)
	}

	handlers.SeedAdminUser(db, adminCredential)

	store := session.New(session.Config{
		Expiration:     24 * time.Hour,
		CookieHTTPOnly: true,
		CookieSameSite: "Lax",
	})

	workerConfig, err := config.NewWorkerConfig()

	if err != nil {
		logger.L.Error("Failed to create worker config:", err)
		os.Exit(1)
	}

	worker.StartWorkerPool(workerConfig.Workers, workerConfig.QueueSize, db)

	app := fiber.New()

	// Add this CORS middleware to allow credentials from the frontend origin
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3001", // Or your frontend's actual origin
		AllowCredentials: true,
	}))

	app.Static("/", "./public")

	app.Post("/api/login", handlers.Login(db, store))
	app.Post("/api/logout", handlers.Logout(store))

	app.Use(handlers.AuthRequired(store))

	app.Post("/api/register", handlers.Register(db))

	app.Post("/api/create/job", routes.CreateJob(db))
	app.Put("/api/update/job", routes.UpdateJob(db))
	app.Delete("/api/delete/job", routes.DeleteJob(db))

	app.Get("/api/jobs", routes.ListJobs(db))
	app.Get("/api/job/:id", routes.GetJobDetails(db))
	app.Get("/api/job/:id/history", routes.ListJobHistory(db))
	app.Get("/api/executions", routes.ListAllExecutions(db))

	app.Get("/api/profile", routes.Profile())
	app.Get("/api/users", routes.ListUsers(db))

	app.Listen("0.0.0.0:3000")

}
