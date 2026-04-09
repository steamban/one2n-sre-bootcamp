package db

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/one2n-sre-bootcamp/student-api/internal/config"
	"github.com/one2n-sre-bootcamp/student-api/pkg/logger"
)

// DB is the global database connection
var DB *sql.DB

// InitDB initializes the PostgreSQL database connection
func InitDB() {
	var err error
	dsn := config.AppConfig.GetDSN()

	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		logger.Log.Error("Failed to open database", "error", err)
		panic(err)
	}

	if err = DB.Ping(); err != nil {
		logger.Log.Error("Failed to ping database", "error", err)
		panic(err)
	}

	logger.Log.Info("Database connection established")
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
	driver, err := postgres.WithInstance(DB, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to create migration driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
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
