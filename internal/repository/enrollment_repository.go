package repository

import (
	"database/sql"
	"fmt"
	"go-student-api/internal/models"
)

type EnrollmentRepository interface {
	Enroll(studentID int, courseID int) error
	Unenroll(studentID int, courseID int) error
	GetCoursesByStudent(studentID int) ([]models.Course, error)
	GetStudentsByCourse(courseID int) ([]models.Student, error)
	GetCount() (int, error)
}

type EnrollmentRepositoryImpl struct {
	DB *sql.DB
}

func NewEnrollmentRepository(db *sql.DB) EnrollmentRepository {
	return &EnrollmentRepositoryImpl{DB: db}
}

func (r *EnrollmentRepositoryImpl) Enroll(studentID int, courseID int) error {
	query := "INSERT INTO enrollments (student_id, course_id) VALUES ($1,$2)"

	_, err := r.DB.Exec(query, studentID, courseID)
	if err != nil {
		return fmt.Errorf("failed to enroll student: %w", err)
	}

	return nil
}

func (r *EnrollmentRepositoryImpl) Unenroll(studentID int, courseID int) error {

	query := "DELETE FROM enrollments WHERE student_id=$1 AND course_id=$2"

	result, err := r.DB.Exec(query, studentID, courseID)
	if err != nil {
		return fmt.Errorf("failed to delete enrollment: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *EnrollmentRepositoryImpl) GetCoursesByStudent(studentID int) ([]models.Course, error) {
	query := ` SELECT c.id, c.title, c.description, c.teacher_id
               FROM courses c
               JOIN enrollments e ON e.course_id = c.id
               WHERE e.student_id = $1`
	var courses []models.Course

	rows, err := r.DB.Query(query, studentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get courses: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var course models.Course
		err := rows.Scan(&course.ID, &course.Title, &course.Description, &course.TeacherID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan course: %w", err)
		}
		courses = append(courses, course)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return courses, nil
}

func (r *EnrollmentRepositoryImpl) GetStudentsByCourse(courseID int) ([]models.Student, error) {
	query := ` SELECT s.id, s.first_name, s.last_name, s.email
			   FROM students s
			   JOIN enrollments e ON e.student_id = s.id
			   WHERE e.course_id = $1`
	var students []models.Student

	rows, err := r.DB.Query(query, courseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get students: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var student models.Student
		err := rows.Scan(&student.ID, &student.FirstName, &student.LastName, &student.Email)
		if err != nil {
			return nil, fmt.Errorf("failed to scan student: %w", err)
		}
		students = append(students, student)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return students, nil

}

func (r *EnrollmentRepositoryImpl) GetCount() (int, error) {

	query := "SELECT COUNT(*) FROM enrollments"

	count := 0

	err := r.DB.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get enrollments count: %w", err)
	}

	return count, nil
}
