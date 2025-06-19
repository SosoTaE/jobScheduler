package config

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"jobScheduler/logger"
	"os"
	"strconv"
)

type AdminCredential struct {
	Username string
	Password string
}

func GetAdminCredential() (AdminCredential, error) {
	if err := godotenv.Load(); err != nil {
		return AdminCredential{}, err
	}

	adminCredential := AdminCredential{}

	adminCredential.Username = os.Getenv("ADMIN_PASSWORD")
	if adminCredential.Username == "" {
		return adminCredential, errors.New("ADMIN_PASSWORD is not set")
	}

	adminCredential.Password = os.Getenv("ADMIN_PASSWORD")
	if adminCredential.Password == "" {
		return adminCredential, errors.New("ADMIN_PASSWORD is not set")
	}

	return adminCredential, nil
}

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
