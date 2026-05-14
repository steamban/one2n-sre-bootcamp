package router

import (
	"github.com/gin-gonic/gin"
	"github.com/one2n-sre-bootcamp/student-api/internal/api/handler"
	"github.com/one2n-sre-bootcamp/student-api/internal/api/middleware"
)

// SetupRouter initializes the Gin router with all routes
func SetupRouter(studentHandler *handler.StudentHandler) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.MetricsMiddleware())

	// API version 1
	v1 := r.Group("/api/v1")
	students := v1.Group("/students")
	students.POST("", studentHandler.CreateStudent)
	students.GET("", studentHandler.GetStudents)
	students.GET("/:id", studentHandler.GetStudentByID)
	students.PATCH("/:id", studentHandler.UpdateStudent)
	students.DELETE("/:id", studentHandler.DeleteStudent)

	return r
}
