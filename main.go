package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io"
	"jobScheduler/handlers"
	"jobScheduler/logger"
	"jobScheduler/models"
	"jobScheduler/routes"
	"log"
	"os"
	"time"
)

func main() {

	file, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	multiWriter := io.MultiWriter(os.Stderr, file)

	log.SetOutput(multiWriter)

	databaseCredential := DatabaseCredential{}
	err = databaseCredential.readENV()
	if err != nil {
		logger.L.Error("Failed to read database credential", err)
	}

	dsn := databaseCredential.ConnectionString()

	// Open a connection to the database.
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.L.Error("Failed to connect to database:", err)
		os.Exit(1)
	}

	err = db.AutoMigrate(&models.Job{})
	if err != nil {
		logger.L.Error("Failed to migrate database:", err)
		os.Exit(1)
	}

	err = db.AutoMigrate(&models.User{})
	if err != nil {
		logger.L.Error("Failed to migrate database:", err)
		os.Exit(1)
	}

	handlers.SeedAdminUser(db, os.Getenv("ADMIN_PASSWORD"))

	// --- 5. Create Session Store ---
	store := session.New(session.Config{
		Expiration:     24 * time.Hour,
		CookieHTTPOnly: true,
		CookieSameSite: "Lax",
	})

	app := fiber.New()

	app.Post("/api/login", handlers.Login(db, store))
	app.Post("/api/logout", handlers.Logout(store))

	app.Use(handlers.AuthRequired(store))

	app.Post("/api/register", handlers.Register(db))

	app.Post("/api/create/job", func(ctx *fiber.Ctx) error {
		return routes.CreateJob(ctx, db)
	})

	app.Put("/api/update/job", func(ctx *fiber.Ctx) error {
		return routes.UpdateJob(ctx, db)
	})

	app.Delete("/api/delete/job", func(ctx *fiber.Ctx) error {
		return routes.DeleteJob(ctx, db)
	})

	app.Get("/api/jobs", func(ctx *fiber.Ctx) error {
		return routes.ListJobs(ctx, db)
	})

	app.Get("/api/profile", routes.Profile())

	app.Listen("0.0.0.0:3000")

}
