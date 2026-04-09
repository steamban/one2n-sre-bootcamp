package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds the application configuration
type Config struct {
	Port       string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
}

// AppConfig is a global variable to access the configuration
var AppConfig Config

// LoadConfig loads the configuration from environment variables or a .env file
func LoadConfig() {
	// Load .env file if it exists, ignore error if it doesn't
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, reading from system environment variables")
	}

	AppConfig = Config{
		Port:       getEnv("PORT", "8080"),
		DBHost:     getEnv("DB_HOST"),
		DBPort:     getEnv("DB_PORT"),
		DBUser:     getEnv("DB_USER"),
		DBPassword: getEnv("DB_PASSWORD"),
		DBName:     getEnv("DB_NAME"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
	}
}

// GetDSN returns the PostgreSQL connection string
func (c *Config) GetDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode)
}

// getEnv reads an environment variable.
// If the variable is not set and no defaultValue is provided, it panics.
// If the variable is not set and a defaultValue is provided, it returns the defaultValue.
func getEnv(key string, defaultValue ...string) string {
	value, exists := os.LookupEnv(key)
	if exists && value != "" {
		return value
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}

	log.Fatalf("Environment variable %s is required but not set", key)
	return ""
}
