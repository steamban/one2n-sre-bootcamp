package repository

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/one2n-sre-bootcamp/student-api/internal/model"
	"github.com/stretchr/testify/require"
)

func TestCreateStudent(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewStudentRepository(db)

		student := &model.Student{
			FirstName:    "John",
			LastName:     "Doe",
			Age:          15,
			Gender:       "Male",
			Email:        "john@example.com",
			Phone:        "1234567890",
			Class:        "10th",
			Rank:         stringPtr("A"),
			AddressLine1: "123 Main St",
			City:         "NYC",
			State:        "NY",
			Pincode:      "123456",
		}

		rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
		mock.ExpectQuery("INSERT INTO students").
			WithArgs(student.FirstName, student.LastName, student.Age, student.Gender,
				student.Email, student.Phone, student.Class, student.Rank,
				student.AddressLine1, student.AddressLine2, student.City, student.State,
				student.Pincode, sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(rows)

		err = repo.CreateStudent(student)

		require.NoError(t, err)
		require.Equal(t, int64(1), student.ID)
		require.False(t, student.CreatedAt.IsZero())
		require.False(t, student.UpdatedAt.IsZero())
	})

	t.Run("db error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewStudentRepository(db)

		student := &model.Student{
			FirstName:    "John",
			LastName:     "Doe",
			Age:          15,
			Gender:       "Male",
			Email:        "john@example.com",
			Phone:        "1234567890",
			Class:        "10th",
			Rank:         stringPtr("A"),
			AddressLine1: "123 Main St",
			City:         "NYC",
			State:        "NY",
			Pincode:      "123456",
		}

		mock.ExpectQuery("INSERT INTO students").
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
				sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
				sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
				sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnError(sql.ErrConnDone)

		err = repo.CreateStudent(student)

		require.Error(t, err)
		require.Equal(t, sql.ErrConnDone, err)
	})
}

func TestGetStudents(t *testing.T) {
	t.Run("success with empty list", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewStudentRepository(db)

		countRows := sqlmock.NewRows([]string{"count"}).AddRow(0)
		mock.ExpectQuery("SELECT COUNT").WillReturnRows(countRows)

		rows := sqlmock.NewRows([]string{
			"id", "first_name", "last_name", "age", "gender", "email", "phone",
			"class", "rank", "address_line1", "address_line2", "city", "state",
			"pincode", "created_at", "updated_at",
		})
		mock.ExpectQuery("SELECT id, first_name").WillReturnRows(rows)

		students, total, err := repo.GetStudents(10, 0)

		require.NoError(t, err)
		require.Equal(t, 0, total)
		require.Len(t, students, 0)
	})

	t.Run("success with students", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewStudentRepository(db)

		now := time.Now()
		countRows := sqlmock.NewRows([]string{"count"}).AddRow(2)
		mock.ExpectQuery("SELECT COUNT").WillReturnRows(countRows)

		rows := sqlmock.NewRows([]string{
			"id", "first_name", "last_name", "age", "gender", "email", "phone",
			"class", "rank", "address_line1", "address_line2", "city", "state",
			"pincode", "created_at", "updated_at",
		}).AddRow(1, "John", "Doe", 15, "Male", "john@example.com", "1234567890",
			"10th", "A", "123 Main St", "", "NYC", "NY", "123456", now, now).
			AddRow(2, "Jane", "Doe", 16, "Female", "jane@example.com", "0987654321",
				"11th", "B", "456 Oak Ave", "Apt 2", "LA", "CA", "654321", now, now)
		mock.ExpectQuery("SELECT id, first_name").WillReturnRows(rows)

		students, total, err := repo.GetStudents(10, 0)

		require.NoError(t, err)
		require.Equal(t, 2, total)
		require.Len(t, students, 2)
		require.Equal(t, "John", students[0].FirstName)
		require.Equal(t, "Jane", students[1].FirstName)
	})

	t.Run("count query error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewStudentRepository(db)

		mock.ExpectQuery("SELECT COUNT").WillReturnError(sql.ErrConnDone)

		students, total, err := repo.GetStudents(10, 0)

		require.Error(t, err)
		require.Equal(t, sql.ErrConnDone, err)
		require.Nil(t, students)
		require.Equal(t, 0, total)
	})

	t.Run("query error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewStudentRepository(db)

		countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
		mock.ExpectQuery("SELECT COUNT").WillReturnRows(countRows)

		mock.ExpectQuery("SELECT id, first_name").WillReturnError(sql.ErrConnDone)

		students, total, err := repo.GetStudents(10, 0)

		require.Error(t, err)
		require.Equal(t, sql.ErrConnDone, err)
		require.Nil(t, students)
		require.Equal(t, 0, total)
	})

	t.Run("scan error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewStudentRepository(db)

		countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
		mock.ExpectQuery("SELECT COUNT").WillReturnRows(countRows)

		rows := sqlmock.NewRows([]string{
			"id", "first_name", "last_name", "age", "gender", "email", "phone",
			"class", "rank", "address_line1", "address_line2", "city", "state",
			"pincode", "created_at", "updated_at",
		}).AddRow("not-an-int", "John", "Doe", 15, "Male", "john@example.com", "1234567890",
			"10th", "A", "123 Main St", "", "NYC", "NY", "123456", time.Now(), time.Now())
		mock.ExpectQuery("SELECT id, first_name").WillReturnRows(rows)

		students, total, err := repo.GetStudents(10, 0)

		require.Error(t, err)
		require.Nil(t, students)
		require.Equal(t, 0, total)
	})

	t.Run("rows error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewStudentRepository(db)

		countRows := sqlmock.NewRows([]string{"count"}).AddRow(2)
		mock.ExpectQuery("SELECT COUNT").WillReturnRows(countRows)

		rows := sqlmock.NewRows([]string{
			"id", "first_name", "last_name", "age", "gender", "email", "phone",
			"class", "rank", "address_line1", "address_line2", "city", "state",
			"pincode", "created_at", "updated_at",
		}).AddRow(1, "John", "Doe", 15, "Male", "john@example.com", "1234567890",
			"10th", "A", "123 Main St", "", "NYC", "NY", "123456", time.Now(), time.Now()).
			RowError(0, sql.ErrConnDone)
		mock.ExpectQuery("SELECT id, first_name").WillReturnRows(rows)

		students, total, err := repo.GetStudents(10, 0)

		require.Error(t, err)
		require.Equal(t, sql.ErrConnDone, err)
		require.Nil(t, students)
		require.Equal(t, 0, total)
	})
}

func TestGetStudentByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewStudentRepository(db)

		now := time.Now()
		rows := sqlmock.NewRows([]string{
			"id", "first_name", "last_name", "age", "gender", "email", "phone",
			"class", "rank", "address_line1", "address_line2", "city", "state",
			"pincode", "created_at", "updated_at",
		}).AddRow(1, "John", "Doe", 15, "Male", "john@example.com", "1234567890",
			"10th", "A", "123 Main St", "", "NYC", "NY", "123456", now, now)
		mock.ExpectQuery("SELECT id, first_name").
			WithArgs(int64(1)).
			WillReturnRows(rows)

		student, err := repo.GetStudentByID(1)

		require.NoError(t, err)
		require.NotNil(t, student)
		require.Equal(t, int64(1), student.ID)
		require.Equal(t, "John", student.FirstName)
	})

	t.Run("not found", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewStudentRepository(db)

		mock.ExpectQuery("SELECT id, first_name").
			WithArgs(int64(999)).
			WillReturnError(sql.ErrNoRows)

		student, err := repo.GetStudentByID(999)

		require.NoError(t, err)
		require.Nil(t, student)
	})

	t.Run("db error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewStudentRepository(db)

		mock.ExpectQuery("SELECT id, first_name").
			WithArgs(int64(1)).
			WillReturnError(sql.ErrConnDone)

		student, err := repo.GetStudentByID(1)

		require.Error(t, err)
		require.Equal(t, sql.ErrConnDone, err)
		require.Nil(t, student)
	})
}

func TestUpdateStudent(t *testing.T) {
	t.Run("success single field", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewStudentRepository(db)

		name := "Jane"
		now := time.Now()
		rows := sqlmock.NewRows([]string{
			"id", "first_name", "last_name", "age", "gender", "email", "phone",
			"class", "rank", "address_line1", "address_line2", "city", "state",
			"pincode", "created_at", "updated_at",
		}).AddRow(1, "Jane", "Doe", 15, "Male", "john@example.com", "1234567890",
			"10th", "A", "123 Main St", "", "NYC", "NY", "123456", now, now)
		mock.ExpectQuery("UPDATE students SET").
			WithArgs("Jane", sqlmock.AnyArg(), int64(1)).
			WillReturnRows(rows)

		updates := model.UpdateStudent{FirstName: &name}
		student, err := repo.UpdateStudent(1, updates)

		require.NoError(t, err)
		require.NotNil(t, student)
		require.Equal(t, "Jane", student.FirstName)
	})

	t.Run("success multiple fields", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewStudentRepository(db)

		name := "Jane"
		age := 16
		now := time.Now()
		rows := sqlmock.NewRows([]string{
			"id", "first_name", "last_name", "age", "gender", "email", "phone",
			"class", "rank", "address_line1", "address_line2", "city", "state",
			"pincode", "created_at", "updated_at",
		}).AddRow(1, "Jane", "Doe", 16, "Male", "john@example.com", "1234567890",
			"10th", "A", "123 Main St", "", "NYC", "NY", "123456", now, now)
		mock.ExpectQuery("UPDATE students SET").
			WithArgs("Jane", 16, sqlmock.AnyArg(), int64(1)).
			WillReturnRows(rows)

		updates := model.UpdateStudent{FirstName: &name, Age: &age}
		student, err := repo.UpdateStudent(1, updates)

		require.NoError(t, err)
		require.NotNil(t, student)
	})

	t.Run("empty updates", func(t *testing.T) {
		db, _, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewStudentRepository(db)

		updates := model.UpdateStudent{}
		student, err := repo.UpdateStudent(1, updates)

		require.NoError(t, err)
		require.Nil(t, student)
	})

	t.Run("db error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewStudentRepository(db)

		name := "Jane"
		mock.ExpectQuery("UPDATE students SET").
			WithArgs("Jane", sqlmock.AnyArg(), int64(1)).
			WillReturnError(sql.ErrConnDone)

		updates := model.UpdateStudent{FirstName: &name}
		student, err := repo.UpdateStudent(1, updates)

		require.Error(t, err)
		require.Equal(t, sql.ErrConnDone, err)
		require.Nil(t, student)
	})
}

func TestDeleteStudent(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewStudentRepository(db)

		now := time.Now()
		rows := sqlmock.NewRows([]string{
			"id", "first_name", "last_name", "age", "gender", "email", "phone",
			"class", "rank", "address_line1", "address_line2", "city", "state",
			"pincode", "created_at", "updated_at", "deleted_at",
		}).AddRow(1, "John", "Doe", 15, "Male", "john@example.com", "1234567890",
			"10th", "A", "123 Main St", "", "NYC", "NY", "123456", now, now, now)
		mock.ExpectQuery("UPDATE students SET deleted_at").
			WithArgs(sqlmock.AnyArg(), int64(1)).
			WillReturnRows(rows)

		student, err := repo.DeleteStudent(1)

		require.NoError(t, err)
		require.NotNil(t, student)
		require.Equal(t, int64(1), student.ID)
		require.NotNil(t, student.DeletedAt)
	})

	t.Run("not found", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewStudentRepository(db)

		mock.ExpectQuery("UPDATE students SET deleted_at").
			WithArgs(sqlmock.AnyArg(), int64(999)).
			WillReturnError(sql.ErrNoRows)

		student, err := repo.DeleteStudent(999)

		require.Error(t, err)
		require.Nil(t, student)
	})

	t.Run("db error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewStudentRepository(db)

		mock.ExpectQuery("UPDATE students SET deleted_at").
			WithArgs(sqlmock.AnyArg(), int64(1)).
			WillReturnError(sql.ErrConnDone)

		student, err := repo.DeleteStudent(1)

		require.Error(t, err)
		require.Equal(t, sql.ErrConnDone, err)
		require.Nil(t, student)
	})
}

func stringPtr(s string) *string {
	return &s
}
