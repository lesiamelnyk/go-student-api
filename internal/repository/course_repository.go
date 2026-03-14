package repository

import (
	"database/sql"
	"fmt"
	"go-student-api/internal/models"
)

type CourseRepository interface {
	GetAll(title string, teacherID int, limit int, offset int) ([]models.Course, error)
	GetByID(id int) (*models.Course, error)
	Create(course *models.Course) error
	Update(id int, course *models.Course) error
	Delete(id int) error
	GetCount() (int, error)
}

type CourseRepositoryImpl struct {
	DB *sql.DB
}

func NewCourseRepository(db *sql.DB) CourseRepository {
	return &CourseRepositoryImpl{DB: db}
}

func (r *CourseRepositoryImpl) Create(course *models.Course) error {
	query := "INSERT INTO courses (title, description, teacher_id) VALUES ($1, $2, $3) RETURNING id"
	err := r.DB.QueryRow(query, course.Title, course.Description, course.TeacherID).Scan(&course.ID)
	if err != nil {
		return fmt.Errorf("failed to create course: %w", err)
	}
	return nil
}

func (r *CourseRepositoryImpl) GetByID(courseId int) (*models.Course, error) {
	query := "SELECT id, title, description, teacher_id FROM courses WHERE id = $1"

	var course models.Course

	err := r.DB.QueryRow(query, courseId).Scan(&course.ID, &course.Title, &course.Description, &course.TeacherID)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get course: %w", err)
	}

	return &course, nil
}

func (r *CourseRepositoryImpl) GetAll(title string, teacherID int, limit int, offset int) ([]models.Course, error) {

	query := "SELECT id, title, description, teacher_id FROM courses WHERE 1=1"

	args := []interface{}{}
	argID := 1

	if title != "" {
		query += fmt.Sprintf(" AND title ILIKE $%d", argID)
		args = append(args, "%"+title+"%")
		argID++
	}

	if teacherID != 0 {
		query += fmt.Sprintf(" AND teacher_id = $%d", argID)
		args = append(args, teacherID)
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

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get courses: %w", err)
	}
	defer rows.Close()

	var courses []models.Course

	for rows.Next() {
		var course models.Course

		err := rows.Scan(
			&course.ID,
			&course.Title,
			&course.Description,
			&course.TeacherID,
		)

		if err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}

		courses = append(courses, course)
	}

	return courses, nil
}

func (r *CourseRepositoryImpl) Update(courseId int, course *models.Course) error {
	query := "UPDATE courses SET title=$1, description=$2, teacher_id=$3 WHERE id=$4"

	result, err := r.DB.Exec(query, course.Title, course.Description, course.TeacherID, courseId)
	if err != nil {
		return fmt.Errorf("failed to update course: %w", err)
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

func (r *CourseRepositoryImpl) Delete(courseId int) error {
	query := "DELETE FROM courses WHERE id=$1"

	result, err := r.DB.Exec(query, courseId)
	if err != nil {
		return fmt.Errorf("failed to delete course: %w", err)
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

func (r *CourseRepositoryImpl) GetCount() (int, error) {

	query := "SELECT COUNT(*) FROM courses"

	count := 0

	err := r.DB.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get courses count: %w", err)
	}

	return count, nil
}
