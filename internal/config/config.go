package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds the application configuration
type Config struct {
	Port   string
	DBPath string
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
		Port:   getEnv("PORT", "8080"),
		DBPath: getEnv("DATABASE_PATH", "students.db"),
	}
}

// getEnv is a helper function to read an environment variable or return a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
