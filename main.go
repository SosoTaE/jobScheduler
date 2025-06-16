package main

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io"
	"jobScheduler/logger"
	"jobScheduler/models"
	"jobScheduler/routes"
	"log"
	"os"
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
	}

	app := fiber.New()

	app.Post("/api/create/job", func(ctx *fiber.Ctx) error {
		return routes.CreateJob(ctx, db)
	})

	app.Put("/api/update/job", func(ctx *fiber.Ctx) error {
		return routes.UpdateJob(ctx, db)
	})

	app.Delete("/api/delete/job", func(ctx *fiber.Ctx) error {
		return routes.DeleteJob(ctx, db)
	})

	app.Listen("0.0.0.0:3000")

}
