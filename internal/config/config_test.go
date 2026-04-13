package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetDSN(t *testing.T) {
	cfg := Config{
		Port:       "8080",
		DBHost:     "localhost",
		DBPort:     "5432",
		DBUser:     "postgres",
		DBPassword: "secret",
		DBName:     "testdb",
		DBSSLMode:  "disable",
	}

	dsn := cfg.GetDSN()

	require.Equal(t, "host=localhost port=5432 user=postgres password=secret dbname=testdb sslmode=disable", dsn)
}

func TestGetEnv(t *testing.T) {
	t.Run("with default when not set", func(t *testing.T) {
		os.Unsetenv("TEST_VAR")
		val := getEnv("TEST_VAR", "default")
		require.Equal(t, "default", val)
	})

	t.Run("returns value when set", func(t *testing.T) {
		os.Setenv("TEST_VAR", "custom")
		defer os.Unsetenv("TEST_VAR")
		val := getEnv("TEST_VAR", "default")
		require.Equal(t, "custom", val)
	})

	t.Run("empty string uses default", func(t *testing.T) {
		os.Setenv("TEST_VAR", "")
		defer os.Unsetenv("TEST_VAR")
		val := getEnv("TEST_VAR", "default")
		require.Equal(t, "default", val)
	})

	t.Run("panics when no default and not set", func(t *testing.T) {
		t.Skip("log.Fatalf calls os.Exit(1) which cannot be caught by test framework")
	})
}

func TestLoadConfig(t *testing.T) {
	t.Run("loads from env vars", func(t *testing.T) {
		t.Setenv("PORT", "3000")
		t.Setenv("DB_HOST", "db.example.com")
		t.Setenv("DB_PORT", "5433")
		t.Setenv("DB_USER", "admin")
		t.Setenv("DB_PASSWORD", "password123")
		t.Setenv("DB_NAME", "mydb")
		t.Setenv("DB_SSLMODE", "require")

		LoadConfig()

		require.Equal(t, "3000", AppConfig.Port)
		require.Equal(t, "db.example.com", AppConfig.DBHost)
		require.Equal(t, "5433", AppConfig.DBPort)
		require.Equal(t, "admin", AppConfig.DBUser)
		require.Equal(t, "password123", AppConfig.DBPassword)
		require.Equal(t, "mydb", AppConfig.DBName)
		require.Equal(t, "require", AppConfig.DBSSLMode)
	})

	t.Run("uses defaults when not set", func(t *testing.T) {
		t.Setenv("PORT", "9090")
		t.Setenv("DB_HOST", "localhost")
		t.Setenv("DB_PORT", "5432")
		t.Setenv("DB_USER", "postgres")
		t.Setenv("DB_PASSWORD", "secret")
		t.Setenv("DB_NAME", "testdb")
		t.Setenv("DB_SSLMODE", "")

		LoadConfig()

		require.Equal(t, "9090", AppConfig.Port)
		require.Equal(t, "disable", AppConfig.DBSSLMode)
	})
}
