package main

import (
	"os"

	"github.com/one2n-sre-bootcamp/student-api/internal/config"
	"github.com/one2n-sre-bootcamp/student-api/internal/db"
	"github.com/one2n-sre-bootcamp/student-api/pkg/logger"
)

func main() {
	// 1. Initialize Logger
	logger.InitLogger()
	defer logger.Sync()

	// 2. Load Config
	config.LoadConfig()

	// 3. Parse Subcommands
	if len(os.Args) < 2 {
		logger.Log.Error("expected 'up' or 'down' subcommands")
		os.Exit(1)
	}

	// 4. Initialize DB Connection
	db.InitDB()
	defer db.CloseDB()

	command := os.Args[1]
	switch command {
	case "up":
		logger.Log.Info("Running migrations UP...")
		if err := db.MigrateUp(); err != nil {
			logger.Log.Error("Migration UP failed", "error", err)
			os.Exit(1)
		}
		logger.Log.Info("Migrations UP completed successfully")

	case "down":
		logger.Log.Info("Running migrations DOWN...")
		if err := db.MigrateDown(); err != nil {
			logger.Log.Error("Migration DOWN failed", "error", err)
			os.Exit(1)
		}
		logger.Log.Info("Migrations DOWN completed successfully")

	default:
		logger.Log.Error("unknown command", "command", command)
		os.Exit(1)
	}
}
