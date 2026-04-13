package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/one2n-sre-bootcamp/student-api/internal/model"
	"github.com/one2n-sre-bootcamp/student-api/internal/repository"
	"github.com/one2n-sre-bootcamp/student-api/pkg/logger"
)

// StudentHandler handles student-related requests
type StudentHandler struct {
	repo repository.StudentRepository
}

// NewStudentHandler creates a new StudentHandler
func NewStudentHandler(repo repository.StudentRepository) *StudentHandler {
	return &StudentHandler{repo: repo}
}

// CreateStudent handles the POST /api/v1/students request
func (h *StudentHandler) CreateStudent(c *gin.Context) {
	var student model.Student
	if err := c.ShouldBindJSON(&student); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.CreateStudent(&student); err != nil {
		logger.Log.Error("failed to create student in database", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create student"})
		return
	}

	c.JSON(http.StatusCreated, student)
}

// GetStudents handles the GET /api/v1/students request
func (h *StudentHandler) GetStudents(c *gin.Context) {
	limit, offset := 10, 0

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}
	if o := c.Query("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	students, total, err := h.repo.GetStudents(limit, offset)
	if err != nil {
		logger.Log.Error("failed to retrieve students", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve students"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   students,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

// GetStudentByID handles the GET /api/v1/students/:id request
func (h *StudentHandler) GetStudentByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid student ID"})
		return
	}

	student, err := h.repo.GetStudentByID(id)
	if err != nil {
		logger.Log.Error("failed to get student by ID", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get student"})
		return
	}

	if student == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "student not found"})
		return
	}

	c.JSON(http.StatusOK, student)
}

// PatchStudent handles the PATCH /api/v1/students/:id request
func (h *StudentHandler) UpdateStudent(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid student ID"})
		return
	}

	var updates model.UpdateStudent
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedStudent, err := h.repo.UpdateStudent(id, updates)
	if err != nil {
		logger.Log.Error("failed to patch student", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update student"})
		return
	}

	if updatedStudent == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "student not found"})
		return
	}

	c.JSON(http.StatusOK, updatedStudent)
}

// DeleteStudent handles the DELETE /api/v1/students/:id request
func (h *StudentHandler) DeleteStudent(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid student ID"})
		return
	}

	deletedStudent, err := h.repo.DeleteStudent(id)
	if err != nil {
		logger.Log.Error("failed to delete student", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete student"})
		return
	}

	if deletedStudent == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "student not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "student deleted successfully", "data": deletedStudent})
}
