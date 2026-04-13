package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/one2n-sre-bootcamp/student-api/internal/model"
)

// StudentRepository defines the interface for student data operations
type StudentRepository interface {
	CreateStudent(student *model.Student) error
	GetStudents() ([]model.Student, error)
	GetStudentByID(id int64) (*model.Student, error)
	UpdateStudent(id int64, updates map[string]any) (*model.Student, error)
	DeleteStudent(id int64) (int64, error)
}

// studentRepo implements StudentRepository
type studentRepo struct {
	db *sql.DB
}

// NewStudentRepository creates a new StudentRepository
func NewStudentRepository(db *sql.DB) StudentRepository {
	return &studentRepo{db: db}
}

// CreateStudent inserts a new student into the database
func (r *studentRepo) CreateStudent(student *model.Student) error {
	student.CreatedAt = time.Now()
	student.UpdatedAt = time.Now()

	query := `
		INSERT INTO students (
			first_name, last_name, age, gender, email, phone, class, rank,
			address_line1, address_line2, city, state, pincode, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		RETURNING id`

	err := r.db.QueryRow(
		query,
		student.FirstName, student.LastName, student.Age, student.Gender,
		student.Email, student.Phone, student.Class, student.Rank,
		student.AddressLine1, student.AddressLine2, student.City, student.State,
		student.Pincode, student.CreatedAt, student.UpdatedAt,
	).Scan(&student.ID)

	if err != nil {
		return err
	}

	return nil
}

// GetStudents retrieves all students from the database
func (r *studentRepo) GetStudents() ([]model.Student, error) {
	query := `
		SELECT id, first_name, last_name, age, gender, email, phone, class, rank,
		       address_line1, address_line2, city, state, pincode, created_at, updated_at
		FROM students
		WHERE deleted_at IS NULL`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []model.Student
	for rows.Next() {
		var s model.Student
		err := rows.Scan(
			&s.ID, &s.FirstName, &s.LastName, &s.Age, &s.Gender, &s.Email, &s.Phone, &s.Class, &s.Rank,
			&s.AddressLine1, &s.AddressLine2, &s.City, &s.State, &s.Pincode, &s.CreatedAt, &s.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		students = append(students, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return students, nil
}

// GetStudentByID retrieves a single student by ID
func (r *studentRepo) GetStudentByID(id int64) (*model.Student, error) {
	query := `
		SELECT id, first_name, last_name, age, gender, email, phone, class, rank,
		       address_line1, address_line2, city, state, pincode, created_at, updated_at
		FROM students
		WHERE id = $1 AND deleted_at IS NULL`

	var s model.Student
	err := r.db.QueryRow(query, id).Scan(
		&s.ID, &s.FirstName, &s.LastName, &s.Age, &s.Gender, &s.Email, &s.Phone, &s.Class, &s.Rank,
		&s.AddressLine1, &s.AddressLine2, &s.City, &s.State, &s.Pincode, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &s, nil
}

// UpdateStudent updates an existing student's information partially
func (r *studentRepo) UpdateStudent(id int64, updates map[string]any) (*model.Student, error) {
	if len(updates) == 0 {
		return nil, nil
	}

	var setClauses []string
	var args []interface{}
	argID := 1

	for column, value := range updates {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", column, argID))
		args = append(args, value)
		argID++
	}

	// Always update updated_at
	setClauses = append(setClauses, fmt.Sprintf("updated_at = $%d", argID))
	args = append(args, time.Now())
	argID++

	query := fmt.Sprintf(
		"UPDATE students SET %s WHERE id = $%d AND deleted_at IS NULL RETURNING id, first_name, last_name, age, gender, email, phone, class, rank, address_line1, address_line2, city, state, pincode, created_at, updated_at",
		strings.Join(setClauses, ", "),
		argID,
	)
	args = append(args, id)

	row := r.db.QueryRow(query, args...)

	var student model.Student
	err := row.Scan(
		&student.ID, &student.FirstName, &student.LastName, &student.Age,
		&student.Gender, &student.Email, &student.Phone, &student.Class,
		&student.Rank, &student.AddressLine1, &student.AddressLine2,
		&student.City, &student.State, &student.Pincode,
		&student.CreatedAt, &student.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &student, nil
}

// DeleteStudent performs a soft delete by setting the deleted_at timestamp
func (r *studentRepo) DeleteStudent(id int64) (int64, error) {
	query := `UPDATE students SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL`
	result, err := r.db.Exec(query, time.Now(), id)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}
