package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/one2n-sre-bootcamp/student-api/internal/api/handler"
	"github.com/one2n-sre-bootcamp/student-api/internal/api/router"
	"github.com/one2n-sre-bootcamp/student-api/internal/config"
	"github.com/one2n-sre-bootcamp/student-api/internal/db"
	"github.com/one2n-sre-bootcamp/student-api/internal/repository"
	"github.com/one2n-sre-bootcamp/student-api/pkg/logger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

	// Initialize Repository and Handler
	studentRepo := repository.NewStudentRepository(db.DB)
	studentHandler := handler.NewStudentHandler(studentRepo)

	// Initialize Gin router
	r := router.SetupRouter(studentHandler)

	// Metrics endpoint
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

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
