package main

import (
	"fmt"
	"go-student-api/internal/database"
	"go-student-api/internal/handlers"
	"go-student-api/internal/repository"
	"log"
	"net/http"
)

func main() {

	db := database.ConnectDB()
	defer db.Close()

	studentRepo := repository.NewStudentRepository(db)
	studentHandler := handlers.NewStudentHandler(studentRepo)

	teacherRepo := repository.NewTeacherRepository(db)
	teacherHandler := handlers.NewTeacherHandler(teacherRepo)

	courseRepo := repository.NewCourseRepository(db)
	courseHandler := handlers.NewCourseHandler(courseRepo)

	enrollmentRepo := repository.NewEnrollmentRepository(db)
	enrollmentHandler := handlers.NewEnrollmentHandler(enrollmentRepo)

	http.HandleFunc("/students", func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {
		case http.MethodGet:
			if id := r.URL.Query().Get("id"); id != "" {
				studentHandler.GetByID(w, r)
			} else {
				studentHandler.GetAll(w, r)
			}
		case http.MethodPost:
			studentHandler.CreateStudent(w, r)

		case http.MethodPut:
			studentHandler.UpdateStudent(w, r)
		case http.MethodDelete:
			studentHandler.DeleteStudent(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

	})

	http.HandleFunc("/teachers", func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {
		case http.MethodGet:
			if id := r.URL.Query().Get("id"); id != "" {
				teacherHandler.GetByID(w, r)
			} else {
				teacherHandler.GetAll(w, r)
			}
		case http.MethodPost:
			teacherHandler.CreateTeacher(w, r)

		case http.MethodPut:
			teacherHandler.UpdateTeacher(w, r)
		case http.MethodDelete:
			teacherHandler.DeleteTeacher(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

	})

	http.HandleFunc("/courses", func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {
		case http.MethodGet:
			if id := r.URL.Query().Get("id"); id != "" {
				courseHandler.GetByID(w, r)
			} else {
				courseHandler.GetAll(w, r)
			}
		case http.MethodPost:
			courseHandler.CreateCourse(w, r)

		case http.MethodPut:
			courseHandler.UpdateCourse(w, r)
		case http.MethodDelete:
			courseHandler.DeleteCourse(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

	})

	http.HandleFunc("/enrollments", func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {
		case http.MethodGet:
			if studentID := r.URL.Query().Get("student_id"); studentID != "" {
				enrollmentHandler.GetCoursesByStudent(w, r)
			} else if courseID := r.URL.Query().Get("course_id"); courseID != "" {
				enrollmentHandler.GetStudentsByCourse(w, r)
			} else {
				http.Error(w, "student_id or course_id is required", http.StatusBadRequest)
				return
			}

		case http.MethodPost:
			enrollmentHandler.CreateEnrollment(w, r)
		case http.MethodDelete:
			enrollmentHandler.DeleteEnrollment(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

	})

	http.HandleFunc("/students/count", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			studentHandler.GetCount(w, r)
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	http.HandleFunc("/teachers/count", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			teacherHandler.GetCount(w, r)
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	http.HandleFunc("/courses/count", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			courseHandler.GetCount(w, r)
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	http.HandleFunc("/enrollments/count", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			enrollmentHandler.GetCount(w, r)
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	fmt.Println("Server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}

}
