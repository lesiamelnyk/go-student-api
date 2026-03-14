package repository

import (
	"database/sql"
	"fmt"
	"go-student-api/internal/models"
)

type TeacherRepository interface {
	GetAll(firstName string, lastName string, department string, limit int, offset int) ([]models.Teacher, error)
	GetByID(id int) (*models.Teacher, error)
	Create(teacher *models.Teacher) error
	Update(id int, teacher *models.Teacher) error
	Delete(id int) error
	GetCount() (int, error)
}

type TeacherRepositoryImpl struct {
	DB *sql.DB
}

func NewTeacherRepository(db *sql.DB) TeacherRepository {
	return &TeacherRepositoryImpl{DB: db}
}

func (r *TeacherRepositoryImpl) Create(teacher *models.Teacher) error {
	query := "INSERT INTO teachers (first_name, last_name, department) VALUES ($1, $2, $3) RETURNING id"
	err := r.DB.QueryRow(query, teacher.FirstName, teacher.LastName, teacher.Department).Scan(&teacher.ID)
	if err != nil {
		return fmt.Errorf("failed to create teacher: %w", err)
	}
	return nil
}

func (r *TeacherRepositoryImpl) GetByID(teacherId int) (*models.Teacher, error) {
	query := "SELECT id, first_name, last_name, department FROM teachers WHERE id = $1"

	var teacher models.Teacher

	err := r.DB.QueryRow(query, teacherId).Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Department)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get teacher: %w", err)
	}

	return &teacher, nil
}

func (r *TeacherRepositoryImpl) GetAll(firstName string, lastName string, department string, limit int, offset int) ([]models.Teacher, error) {
	query := "SELECT id, first_name, last_name, department FROM teachers WHERE 1=1"

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

	if department != "" {
		query += fmt.Sprintf(" AND department ILIKE $%d", argID)
		args = append(args, "%"+department+"%")
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

	var teachers []models.Teacher

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get teachers: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var teacher models.Teacher
		err := rows.Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Department)
		if err != nil {
			return nil, fmt.Errorf("failed to scan teacher: %w", err)
		}
		teachers = append(teachers, teacher)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return teachers, nil
}

func (r *TeacherRepositoryImpl) Update(teacherId int, teacher *models.Teacher) error {
	query := "UPDATE teachers SET first_name=$1, last_name=$2, department=$3 WHERE id=$4"

	result, err := r.DB.Exec(query, teacher.FirstName, teacher.LastName, teacher.Department, teacherId)
	if err != nil {
		return fmt.Errorf("failed to update teacher: %w", err)
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

func (r *TeacherRepositoryImpl) Delete(teacherId int) error {
	query := "DELETE FROM teachers WHERE id=$1"

	result, err := r.DB.Exec(query, teacherId)
	if err != nil {
		return fmt.Errorf("failed to delete teacher: %w", err)
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

func (r *TeacherRepositoryImpl) GetCount() (int, error) {

	query := "SELECT COUNT(*) FROM teachers"

	count := 0

	err := r.DB.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get teacher count: %w", err)
	}

	return count, nil
}
