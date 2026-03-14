package repository

import (
	"database/sql"
	"fmt"
	"go-student-api/internal/models"
)

type StudentRepository interface {
	GetAll(firstName string, lastName string, email string, limit int, offset int) ([]models.Student, error)
	GetByID(id int) (*models.Student, error)
	Create(student *models.Student) error
	Update(id int, student *models.Student) error
	Delete(id int) error
	GetCount() (int, error)
}

type StudentRepositoryImpl struct {
	DB *sql.DB
}

func NewStudentRepository(db *sql.DB) StudentRepository {
	return &StudentRepositoryImpl{DB: db}
}

func (r *StudentRepositoryImpl) Create(student *models.Student) error {
	query := "INSERT INTO students (first_name, last_name, email) VALUES ($1, $2, $3) RETURNING id"
	err := r.DB.QueryRow(query, student.FirstName, student.LastName, student.Email).Scan(&student.ID)
	if err != nil {
		return fmt.Errorf("failed to create student: %w", err)
	}
	return nil
}

func (r *StudentRepositoryImpl) GetByID(studentId int) (*models.Student, error) {
	query := "SELECT id, first_name, last_name, email FROM students WHERE id = $1"

	var student models.Student

	err := r.DB.QueryRow(query, studentId).Scan(&student.ID, &student.FirstName, &student.LastName, &student.Email)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get student: %w", err)
	}

	return &student, nil
}

func (r *StudentRepositoryImpl) GetAll(firstName string, lastName string, email string, limit int, offset int) ([]models.Student, error) {
	query := "SELECT id, first_name, last_name, email FROM students WHERE 1=1"

	args := []interface{}{}
	argID := 1

	if firstName != "" {
		query += fmt.Sprintf(" AND first_name ILIKE $%d", argID)
		args = append(args, "%"+firstName+"%")
		argID++
	}

	if lastName != "" {
		query += fmt.Sprintf(" AND last_name ILIKE $%d", argID)
		args = append(args, "%"+lastName+"%")
		argID++
	}

	if email != "" {
		query += fmt.Sprintf(" AND email ILIKE $%d", argID)
		args = append(args, "%"+email+"%")
		argID++
	}

	if limit == 0 {
		limit = 10
	}
	query += fmt.Sprintf(" LIMIT $%d", argID)
	args = append(args, limit)
	argID++

	if offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argID)
		args = append(args, offset)
		argID++
	}

	var students []models.Student

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get students: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var student models.Student
		err := rows.Scan(&student.ID, &student.FirstName, &student.LastName, &student.Email)
		if err != nil {
			return nil, fmt.Errorf("failed to get student: %w", err)
		}
		students = append(students, student)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return students, nil
}

func (r *StudentRepositoryImpl) Update(studentId int, student *models.Student) error {
	query := "UPDATE students SET first_name=$1, last_name=$2, email=$3 WHERE id=$4"

	result, err := r.DB.Exec(query, student.FirstName, student.LastName, student.Email, studentId)
	if err != nil {
		return fmt.Errorf("failed to update student: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *StudentRepositoryImpl) Delete(studentId int) error {
	query := "DELETE FROM students WHERE id=$1"

	result, err := r.DB.Exec(query, studentId)
	if err != nil {
		return fmt.Errorf("failed to delete student: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *StudentRepositoryImpl) GetCount() (int, error) {

	query := "SELECT COUNT(*) FROM students"

	count := 0

	err := r.DB.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get student count: %w", err)
	}

	return count, nil
}
