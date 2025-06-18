package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"jobScheduler/logger"
	"os"
	"strconv"
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

// WorkerConfig holds the configuration for the concurrent worker pool.
// The name is changed from 'Threading' to be more idiomatic to Go's concurrency model.
type WorkerConfig struct {
	QueueSize int
	Workers   int
}

// NewWorkerConfig creates a new configuration object by reading from environment variables.
// It provides sensible defaults if variables are not set.
func NewWorkerConfig() (*WorkerConfig, error) {
	// Attempt to load .env file but don't treat a "not found" error as fatal.
	if err := godotenv.Load(); err != nil {
		logger.L.Warn("Could not load .env file, reading from OS environment. This is normal in production.")
	}

	config := &WorkerConfig{}
	var err error

	// --- Get Worker Pool Size ---
	workersStr := os.Getenv("WORKERS")
	if workersStr == "" {
		config.Workers = 5 // Default value
	} else {
		config.Workers, err = strconv.Atoi(workersStr)
		if err != nil {
			return nil, fmt.Errorf("invalid WORKERS value: must be an integer")
		}
	}

	// --- Get Queue Size ---
	queueSizeStr := os.Getenv("QUEUE_SIZE")
	if queueSizeStr == "" {
		config.QueueSize = 100 // Default value
	} else {
		config.QueueSize, err = strconv.Atoi(queueSizeStr)
		if err != nil {
			return nil, fmt.Errorf("invalid QUEUE_SIZE value: must be an integer")
		}
	}

	// Validate the values
	if config.Workers <= 0 {
		return nil, fmt.Errorf("WORKERS count must be positive")
	}
	if config.QueueSize <= 0 {
		return nil, fmt.Errorf("QUEUE_SIZE must be positive")
	}

	// Return the populated config and a nil error on success
	return config, nil
}
