package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/one2n-sre-bootcamp/student-api/internal/config"
	"github.com/one2n-sre-bootcamp/student-api/internal/db"
	"github.com/one2n-sre-bootcamp/student-api/pkg/logger"
)

func main() {
	// Initialize logger
	logger.InitLogger()
	defer logger.Sync()

	// Load configuration
	config.LoadConfig()

	// Initialize Database connection only (no migrations)
	db.InitDB()
	defer db.CloseDB()

	// Initialize Gin router
	r := gin.Default()

	// Healthcheck endpoint
	r.GET("/healthcheck", func(c *gin.Context) {
		logger.Log.Info("Healthcheck hit")
		c.JSON(200, gin.H{
			"status": "UP",
		})
	})

	// Start server
	logger.Log.Info("Starting server", "port", config.AppConfig.Port)
	if err := r.Run(":" + config.AppConfig.Port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
