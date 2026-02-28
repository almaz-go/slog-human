package main

import (
	"log/slog"
	"os"
	"github.com/almaz-go/slog-human" // Убедись, что путь совпадает с go.mod
)

func main() {
	// Initialize our pretty handler
	handler := sloghuman.NewHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	// Set it as default
	logger := slog.New(handler)
	slog.SetDefault(logger)

	// Let's see the magic!
	slog.Info("Server started", "port", 8080, "env", "production")
	slog.Warn("High memory usage detected", "usage_percent", 85)
	slog.Error("Database connection failed", "retry_count", 3, "error", "timeout")
	slog.Debug("Debugging complex logic", "step", "auth_check")
}