package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"jobScheduler/logger"
	"os"
)

type DatabaseCredential struct {
	DB_HOST string
	DB_USER string
	DB_PASS string
	DB_NAME string
	DB_PORT string
}

func (c *DatabaseCredential) readENV() error {
	if err := godotenv.Load(); err != nil {
		return err
	}
	c.DB_HOST = os.Getenv("DB_HOST")
	c.DB_USER = os.Getenv("DB_USER")
	c.DB_PASS = os.Getenv("DB_PASSWORD") // Reads from DB_PASSWORD
	c.DB_NAME = os.Getenv("DB_NAME")
	c.DB_PORT = os.Getenv("DB_PORT")

	logger.L.Info("Successfully read database credentials",
		"db_host", c.DB_HOST,
		"db_user", c.DB_USER,
		"db_name", c.DB_NAME,
		"db_port", c.DB_PORT,
		"db_pass", "***MASKED***", // NEVER LOG THE REAL PASSWORD
	)

	return nil
}

func (c *DatabaseCredential) ConnectionString() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
		c.DB_HOST, c.DB_USER, c.DB_PASS, c.DB_NAME, c.DB_PORT)
}
