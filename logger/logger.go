package logger

import (
	"log/slog"
	"os"
)

var L *slog.Logger

// init runs automatically when the package is imported.
func init() {
	// Configure the logger to write to a file.
	// os.O_APPEND: Add to the file, don't overwrite.
	// os.O_CREATE: Create the file if it doesn't exist.
	logFile, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic("failed to open log file: " + err.Error())
	}

	handler := slog.NewJSONHandler(logFile, &slog.HandlerOptions{
		// Set the minimum log level.
		// Levels are: Debug, Info, Warn, Error
		Level: slog.LevelDebug,
	})

	// Create a new logger with our configured handler.
	L = slog.New(handler)

	slog.SetDefault(L)

	L.Info("Logger initialized successfully")
}
