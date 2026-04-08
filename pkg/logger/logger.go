package logger

import (
	"log/slog"
	"os"
)

// Log is the global logger instance
var Log *slog.Logger

// InitLogger initializes a new slog logger
func InitLogger() {
	// For development, we use a text-based handler.
	// This can be swapped for slog.NewJSONHandler for production.
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	Log = slog.New(handler)
	slog.SetDefault(Log)
}

// Sync is a no-op for slog as it doesn't typically require flushing
// like zap does, but we keep the signature for main.go compatibility.
func Sync() {}
