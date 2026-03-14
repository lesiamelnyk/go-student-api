package handlers

import (
	"encoding/json"
	"errors"
	"go-student-api/internal/models"
	"go-student-api/internal/repository"
	"net/http"
	"strconv"
)

type EnrollmentHandler struct {
	repo repository.EnrollmentRepository
}

func (h *EnrollmentHandler) DeleteEnrollment(w http.ResponseWriter, r *http.Request) {

	studentStr := r.URL.Query().Get("student_id")
	courseStr := r.URL.Query().Get("course_id")

	studentID, err := strconv.Atoi(studentStr)
	if err != nil {
		http.Error(w, "invalid student id", http.StatusBadRequest)
		return
	}

	courseID, err := strconv.Atoi(courseStr)
	if err != nil {
		http.Error(w, "invalid course id", http.StatusBadRequest)
		return
	}

	err = h.repo.Unenroll(studentID, courseID)
	if errors.Is(err, repository.ErrNotFound) {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "failed to delete enrollment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *EnrollmentHandler) CreateEnrollment(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var enrollment models.Enrollment

	err := json.NewDecoder(r.Body).Decode(&enrollment)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err = h.repo.Enroll(enrollment.StudentID, enrollment.CourseID)
	if err != nil {
		http.Error(w, "failed to enroll student", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(enrollment)
}

func (h *EnrollmentHandler) GetCoursesByStudent(w http.ResponseWriter, r *http.Request) {

	idStr := r.URL.Query().Get("student_id")
	if idStr == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid student id", http.StatusBadRequest)
		return
	}

	courses, err := h.repo.GetCoursesByStudent(id)
	if errors.Is(err, repository.ErrNotFound) {
		http.Error(w, "course not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "failed to get course", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courses)
}

func (h *EnrollmentHandler) GetStudentsByCourse(w http.ResponseWriter, r *http.Request) {

	courseStr := r.URL.Query().Get("course_id")

	courseID, err := strconv.Atoi(courseStr)
	if err != nil {
		http.Error(w, "invalid course id", http.StatusBadRequest)
		return
	}

	students, err := h.repo.GetStudentsByCourse(courseID)
	if errors.Is(err, repository.ErrNotFound) {
		http.Error(w, "student not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "failed to get student", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(students)
}

func NewEnrollmentHandler(repo repository.EnrollmentRepository) *EnrollmentHandler {
	return &EnrollmentHandler{repo: repo}
}

func (h *EnrollmentHandler) GetCount(w http.ResponseWriter, r *http.Request) {

	count, err := h.repo.GetCount()
	if err != nil {
		http.Error(w, "failed to get count", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]int{
		"count": count,
	})
}
