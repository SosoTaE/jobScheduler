package main

import (
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

	"github.com/glebarez/sqlite" // Pure Go SQLite driver
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/session"
	"gorm.io/gorm"
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

	// Updated Session Configuration
	store := session.New(session.Config{
		Expiration:     24 * time.Hour,
		CookieHTTPOnly: true,
		CookieSameSite: "Lax", // Better for local dev
		CookieSecure:   false,
	})

	workerConfig, err := config.NewWorkerConfig()

	if err != nil {
		logger.L.Error("Failed to create worker config:", err)
		os.Exit(1)
	}

	worker.StartWorkerPool(workerConfig.Workers, workerConfig.QueueSize, db)

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOriginsFunc: func(origin string) bool {
			// This allows both localhost and 127.0.0.1 dynamically
			return origin == "http://localhost:3000" || origin == "http://127.0.0.1:3000"
		},
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowCredentials: true,
	}))
	app.Static("/", "./public")

	api := app.Group("/api")

	api.Post("/login", handlers.Login(db, store))

	api.Use(handlers.AuthRequired(store, db))

	api.Post("/logout", handlers.Logout(store))

	api.Post("/register", handlers.Register(db))

	api.Post("/execute", routes.Execute(db))

	api.Post("/create/job", routes.CreateJob(db))
	api.Put("/update/job", routes.UpdateJob(db))
	api.Delete("/delete/job", routes.DeleteJob(db))

	api.Get("/jobs", routes.ListJobs(db))
	api.Get("/job/:id", routes.GetJobDetails(db))
	api.Get("/job/:id/history", routes.ListJobHistory(db))
	api.Get("/executions", routes.ListAllExecutions(db))

	api.Get("/profile", routes.Profile())
	api.Get("/users", routes.ListUsers(db))

	api.Post("/generate-api-key", routes.GenerateAPIKey(db))

	app.Listen("0.0.0.0:3000")
}
