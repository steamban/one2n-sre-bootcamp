package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/one2n-sre-bootcamp/student-api/internal/model"
	"github.com/one2n-sre-bootcamp/student-api/pkg/logger"
	"github.com/stretchr/testify/require"
)

type mockStudentRepo struct {
	createStudentFn  func(student *model.Student) error
	getStudentsFn    func(limit, offset int) ([]model.Student, int, error)
	getStudentByIDFn func(id int64) (*model.Student, error)
	updateStudentFn  func(id int64, updates model.UpdateStudent) (*model.Student, error)
	deleteStudentFn  func(id int64) (*model.Student, error)
}

func (m *mockStudentRepo) CreateStudent(student *model.Student) error {
	if m.createStudentFn != nil {
		return m.createStudentFn(student)
	}
	return nil
}

func (m *mockStudentRepo) GetStudents(limit, offset int) ([]model.Student, int, error) {
	if m.getStudentsFn != nil {
		return m.getStudentsFn(limit, offset)
	}
	return []model.Student{}, 0, nil
}

func (m *mockStudentRepo) GetStudentByID(id int64) (*model.Student, error) {
	if m.getStudentByIDFn != nil {
		return m.getStudentByIDFn(id)
	}
	return nil, nil
}

func (m *mockStudentRepo) UpdateStudent(id int64, updates model.UpdateStudent) (*model.Student, error) {
	if m.updateStudentFn != nil {
		return m.updateStudentFn(id, updates)
	}
	return nil, nil
}

func (m *mockStudentRepo) DeleteStudent(id int64) (*model.Student, error) {
	if m.deleteStudentFn != nil {
		return m.deleteStudentFn(id)
	}
	return nil, nil
}

func setupRouter(repo *mockStudentRepo) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewStudentHandler(repo)
	r.POST("/api/v1/students", h.CreateStudent)
	r.GET("/api/v1/students", h.GetStudents)
	r.GET("/api/v1/students/:id", h.GetStudentByID)
	r.PATCH("/api/v1/students/:id", h.UpdateStudent)
	r.DELETE("/api/v1/students/:id", h.DeleteStudent)
	return r
}

func init() {
	logger.InitLogger()
}

func TestCreateStudent(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := &mockStudentRepo{
			createStudentFn: func(student *model.Student) error {
				student.ID = 1
				student.CreatedAt = time.Now()
				student.UpdatedAt = time.Now()
				return nil
			},
		}
		router := setupRouter(repo)

		payload := map[string]interface{}{
			"first_name":    "John",
			"last_name":     "Doe",
			"age":           15,
			"gender":        "Male",
			"email":         "john@example.com",
			"phone":         "1234567890",
			"class":         "10th",
			"rank":          "A",
			"address_line1": "123 Main St",
			"city":          "NYC",
			"state":         "NY",
			"pincode":       "123456",
		}
		body, _ := json.Marshal(payload)

		req, _ := http.NewRequest("POST", "/api/v1/students", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusCreated, w.Code)

		var resp model.Student
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		require.Equal(t, "John", resp.FirstName)
		require.Equal(t, int64(1), resp.ID)
	})

	t.Run("invalid JSON", func(t *testing.T) {
		repo := &mockStudentRepo{}
		router := setupRouter(repo)

		req, _ := http.NewRequest("POST", "/api/v1/students", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("validation failure - missing required field", func(t *testing.T) {
		repo := &mockStudentRepo{}
		router := setupRouter(repo)

		payload := map[string]interface{}{
			"first_name": "John",
		}
		body, _ := json.Marshal(payload)

		req, _ := http.NewRequest("POST", "/api/v1/students", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestGetStudents(t *testing.T) {
	t.Run("success with defaults", func(t *testing.T) {
		repo := &mockStudentRepo{
			getStudentsFn: func(limit, offset int) ([]model.Student, int, error) {
				return []model.Student{
					{ID: 1, FirstName: "John", LastName: "Doe"},
				}, 1, nil
			},
		}
		router := setupRouter(repo)

		req, _ := http.NewRequest("GET", "/api/v1/students", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)

		var resp map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		require.Equal(t, float64(1), resp["total"])
		require.Equal(t, float64(10), resp["limit"])
		require.Equal(t, float64(0), resp["offset"])
	})

	t.Run("custom limit and offset", func(t *testing.T) {
		var capturedLimit, capturedOffset int
		repo := &mockStudentRepo{
			getStudentsFn: func(limit, offset int) ([]model.Student, int, error) {
				capturedLimit = limit
				capturedOffset = offset
				return []model.Student{}, 0, nil
			},
		}
		router := setupRouter(repo)

		req, _ := http.NewRequest("GET", "/api/v1/students?limit=20&offset=5", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		require.Equal(t, 20, capturedLimit)
		require.Equal(t, 5, capturedOffset)
	})

	t.Run("invalid limit uses default", func(t *testing.T) {
		repo := &mockStudentRepo{
			getStudentsFn: func(limit, offset int) ([]model.Student, int, error) {
				return []model.Student{}, 0, nil
			},
		}
		router := setupRouter(repo)

		req, _ := http.NewRequest("GET", "/api/v1/students?limit=-5", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
	})
}

func TestGetStudentByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := &mockStudentRepo{
			getStudentByIDFn: func(id int64) (*model.Student, error) {
				return &model.Student{ID: id, FirstName: "John", LastName: "Doe"}, nil
			},
		}
		router := setupRouter(repo)

		req, _ := http.NewRequest("GET", "/api/v1/students/1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)

		var resp model.Student
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		require.Equal(t, int64(1), resp.ID)
		require.Equal(t, "John", resp.FirstName)
	})

	t.Run("invalid ID - non-numeric", func(t *testing.T) {
		repo := &mockStudentRepo{}
		router := setupRouter(repo)

		req, _ := http.NewRequest("GET", "/api/v1/students/abc", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("not found", func(t *testing.T) {
		repo := &mockStudentRepo{
			getStudentByIDFn: func(id int64) (*model.Student, error) {
				return nil, nil
			},
		}
		router := setupRouter(repo)

		req, _ := http.NewRequest("GET", "/api/v1/students/999", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestUpdateStudent(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		name := "Jane"
		repo := &mockStudentRepo{
			updateStudentFn: func(id int64, updates model.UpdateStudent) (*model.Student, error) {
				return &model.Student{ID: id, FirstName: name, LastName: "Doe"}, nil
			},
		}
		router := setupRouter(repo)

		payload := map[string]interface{}{
			"first_name": "Jane",
		}
		body, _ := json.Marshal(payload)

		req, _ := http.NewRequest("PATCH", "/api/v1/students/1", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)

		var resp model.Student
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		require.Equal(t, "Jane", resp.FirstName)
	})

	t.Run("invalid ID - non-numeric", func(t *testing.T) {
		repo := &mockStudentRepo{}
		router := setupRouter(repo)

		payload := map[string]interface{}{"first_name": "Jane"}
		body, _ := json.Marshal(payload)

		req, _ := http.NewRequest("PATCH", "/api/v1/students/abc", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("not found", func(t *testing.T) {
		repo := &mockStudentRepo{
			updateStudentFn: func(id int64, updates model.UpdateStudent) (*model.Student, error) {
				return nil, nil
			},
		}
		router := setupRouter(repo)

		payload := map[string]interface{}{"first_name": "Jane"}
		body, _ := json.Marshal(payload)

		req, _ := http.NewRequest("PATCH", "/api/v1/students/999", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("empty updates", func(t *testing.T) {
		repo := &mockStudentRepo{
			updateStudentFn: func(id int64, updates model.UpdateStudent) (*model.Student, error) {
				return nil, nil
			},
		}
		router := setupRouter(repo)

		payload := map[string]interface{}{}
		body, _ := json.Marshal(payload)

		req, _ := http.NewRequest("PATCH", "/api/v1/students/1", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestDeleteStudent(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := &mockStudentRepo{
			deleteStudentFn: func(id int64) (*model.Student, error) {
				return &model.Student{ID: id, FirstName: "John"}, nil
			},
		}
		router := setupRouter(repo)

		req, _ := http.NewRequest("DELETE", "/api/v1/students/1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)

		var resp map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		require.Equal(t, "student deleted successfully", resp["message"])
	})

	t.Run("invalid ID - non-numeric", func(t *testing.T) {
		repo := &mockStudentRepo{}
		router := setupRouter(repo)

		req, _ := http.NewRequest("DELETE", "/api/v1/students/abc", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("not found", func(t *testing.T) {
		repo := &mockStudentRepo{
			deleteStudentFn: func(id int64) (*model.Student, error) {
				return nil, nil
			},
		}
		router := setupRouter(repo)

		req, _ := http.NewRequest("DELETE", "/api/v1/students/999", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("repo error", func(t *testing.T) {
		repo := &mockStudentRepo{
			deleteStudentFn: func(id int64) (*model.Student, error) {
				return nil, errors.New("database error")
			},
		}
		router := setupRouter(repo)

		req, _ := http.NewRequest("DELETE", "/api/v1/students/1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestCreateStudent_RepoError(t *testing.T) {
	t.Run("repo error", func(t *testing.T) {
		repo := &mockStudentRepo{
			createStudentFn: func(student *model.Student) error {
				return errors.New("database error")
			},
		}
		router := setupRouter(repo)

		payload := map[string]interface{}{
			"first_name":    "John",
			"last_name":     "Doe",
			"age":           15,
			"gender":        "Male",
			"email":         "john@example.com",
			"phone":         "1234567890",
			"class":         "10th",
			"rank":          "A",
			"address_line1": "123 Main St",
			"city":          "NYC",
			"state":         "NY",
			"pincode":       "123456",
		}
		body, _ := json.Marshal(payload)

		req, _ := http.NewRequest("POST", "/api/v1/students", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestGetStudents_RepoError(t *testing.T) {
	t.Run("repo error", func(t *testing.T) {
		repo := &mockStudentRepo{
			getStudentsFn: func(limit, offset int) ([]model.Student, int, error) {
				return nil, 0, errors.New("database error")
			},
		}
		router := setupRouter(repo)

		req, _ := http.NewRequest("GET", "/api/v1/students", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestGetStudentByID_RepoError(t *testing.T) {
	t.Run("repo error", func(t *testing.T) {
		repo := &mockStudentRepo{
			getStudentByIDFn: func(id int64) (*model.Student, error) {
				return nil, errors.New("database error")
			},
		}
		router := setupRouter(repo)

		req, _ := http.NewRequest("GET", "/api/v1/students/1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestUpdateStudent_RepoError(t *testing.T) {
	t.Run("repo error", func(t *testing.T) {
		repo := &mockStudentRepo{
			updateStudentFn: func(id int64, updates model.UpdateStudent) (*model.Student, error) {
				return nil, errors.New("database error")
			},
		}
		router := setupRouter(repo)

		payload := map[string]interface{}{"first_name": "Jane"}
		body, _ := json.Marshal(payload)

		req, _ := http.NewRequest("PATCH", "/api/v1/students/1", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
