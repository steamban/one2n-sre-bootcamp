package db

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/one2n-sre-bootcamp/student-api/internal/config"
	"github.com/one2n-sre-bootcamp/student-api/pkg/logger"
	_ "modernc.org/sqlite"
)

// DB is the global database connection
var DB *sql.DB

// InitDB initializes the SQLite database connection
func InitDB() {
	var err error
	dbPath := config.AppConfig.DBPath

	dir := filepath.Dir(dbPath)
	if dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			logger.Log.Error("Failed to create database directory", "error", err)
			panic(err)
		}
	}

	DB, err = sql.Open("sqlite", dbPath)
	if err != nil {
		logger.Log.Error("Failed to open database", "error", err)
		panic(err)
	}

	if err = DB.Ping(); err != nil {
		logger.Log.Error("Failed to ping database", "error", err)
		panic(err)
	}

	logger.Log.Info("Database connection established", "path", dbPath)
}

// MigrateUp runs all pending migrations
func MigrateUp() error {
	m, err := getMigrateInstance()
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}

// MigrateDown rolls back the last migration
func MigrateDown() error {
	m, err := getMigrateInstance()
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Steps(-1); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}

func getMigrateInstance() (*migrate.Migrate, error) {
	driver, err := sqlite.WithInstance(DB, &sqlite.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to create migration driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"sqlite", driver)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize migration: %w", err)
	}
	return m, nil
}

// CloseDB closes the database connection
func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}
