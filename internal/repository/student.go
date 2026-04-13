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
	GetStudents(limit, offset int) ([]model.Student, int, error)
	GetStudentByID(id int64) (*model.Student, error)
	UpdateStudent(id int64, updates model.UpdateStudent) (*model.Student, error)
	DeleteStudent(id int64) (*model.Student, error)
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

// GetStudents retrieves students with pagination
func (r *studentRepo) GetStudents(limit, offset int) ([]model.Student, int, error) {
	var total int
	countQuery := `SELECT COUNT(*) FROM students WHERE deleted_at IS NULL`
	if err := r.db.QueryRow(countQuery).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT id, first_name, last_name, age, gender, email, phone, class, rank,
		       address_line1, address_line2, city, state, pincode, created_at, updated_at
		FROM students
		WHERE deleted_at IS NULL
		ORDER BY id
		LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	students := make([]model.Student, 0, limit)
	for rows.Next() {
		var s model.Student
		err := rows.Scan(
			&s.ID, &s.FirstName, &s.LastName, &s.Age, &s.Gender, &s.Email, &s.Phone, &s.Class, &s.Rank,
			&s.AddressLine1, &s.AddressLine2, &s.City, &s.State, &s.Pincode, &s.CreatedAt, &s.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		students = append(students, s)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return students, total, nil
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
func (r *studentRepo) UpdateStudent(id int64, updates model.UpdateStudent) (*model.Student, error) {
	var setClauses []string
	var args []any
	argID := 1

	if updates.FirstName != nil {
		setClauses = append(setClauses, fmt.Sprintf("first_name = $%d", argID))
		args = append(args, *updates.FirstName)
		argID++
	}
	if updates.LastName != nil {
		setClauses = append(setClauses, fmt.Sprintf("last_name = $%d", argID))
		args = append(args, *updates.LastName)
		argID++
	}
	if updates.Age != nil {
		setClauses = append(setClauses, fmt.Sprintf("age = $%d", argID))
		args = append(args, *updates.Age)
		argID++
	}
	if updates.Gender != nil {
		setClauses = append(setClauses, fmt.Sprintf("gender = $%d", argID))
		args = append(args, *updates.Gender)
		argID++
	}
	if updates.Email != nil {
		setClauses = append(setClauses, fmt.Sprintf("email = $%d", argID))
		args = append(args, *updates.Email)
		argID++
	}
	if updates.Phone != nil {
		setClauses = append(setClauses, fmt.Sprintf("phone = $%d", argID))
		args = append(args, *updates.Phone)
		argID++
	}
	if updates.Class != nil {
		setClauses = append(setClauses, fmt.Sprintf("class = $%d", argID))
		args = append(args, *updates.Class)
		argID++
	}
	if updates.Rank != nil {
		setClauses = append(setClauses, fmt.Sprintf("rank = $%d", argID))
		args = append(args, *updates.Rank)
		argID++
	}
	if updates.AddressLine1 != nil {
		setClauses = append(setClauses, fmt.Sprintf("address_line1 = $%d", argID))
		args = append(args, *updates.AddressLine1)
		argID++
	}
	if updates.AddressLine2 != nil {
		setClauses = append(setClauses, fmt.Sprintf("address_line2 = $%d", argID))
		args = append(args, *updates.AddressLine2)
		argID++
	}
	if updates.City != nil {
		setClauses = append(setClauses, fmt.Sprintf("city = $%d", argID))
		args = append(args, *updates.City)
		argID++
	}
	if updates.State != nil {
		setClauses = append(setClauses, fmt.Sprintf("state = $%d", argID))
		args = append(args, *updates.State)
		argID++
	}
	if updates.Pincode != nil {
		setClauses = append(setClauses, fmt.Sprintf("pincode = $%d", argID))
		args = append(args, *updates.Pincode)
		argID++
	}

	if len(setClauses) == 0 {
		return nil, nil
	}

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
func (r *studentRepo) DeleteStudent(id int64) (*model.Student, error) {
	query := `
		UPDATE students SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL 
		RETURNING id, first_name, last_name, age, gender, email, phone, class, rank, address_line1, address_line2, city, state, pincode, created_at, updated_at, deleted_at`

	row := r.db.QueryRow(query, time.Now(), id)

	var student model.Student
	err := row.Scan(
		&student.ID, &student.FirstName, &student.LastName, &student.Age,
		&student.Gender, &student.Email, &student.Phone, &student.Class,
		&student.Rank, &student.AddressLine1, &student.AddressLine2,
		&student.City, &student.State, &student.Pincode,
		&student.CreatedAt, &student.UpdatedAt, &student.DeletedAt,
	)
	if err != nil {
		return nil, err
	}

	return &student, nil
}
