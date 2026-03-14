package handlers

import (
	"encoding/json"
	"errors"
	"go-student-api/internal/models"
	"go-student-api/internal/repository"
	"net/http"
	"strconv"
)

type CourseHandler struct {
	repo repository.CourseRepository
}

func (h *CourseHandler) DeleteCourse(w http.ResponseWriter, r *http.Request) {

	id, err := getID(r)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	err = h.repo.Delete(id)

	if errors.Is(err, repository.ErrNotFound) {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "failed to delete course", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *CourseHandler) UpdateCourse(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	id, err := getID(r)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var course models.Course

	err = json.NewDecoder(r.Body).Decode(&course)
	if err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	err = h.repo.Update(id, &course)
	if errors.Is(err, repository.ErrNotFound) {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "failed to update course", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *CourseHandler) GetByID(w http.ResponseWriter, r *http.Request) {

	id, err := getID(r)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	course, err := h.repo.GetByID(id)
	if errors.Is(err, repository.ErrNotFound) {
		http.Error(w, "course not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "failed to get course", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(course)
}

func (h *CourseHandler) CreateCourse(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var course models.Course

	err := json.NewDecoder(r.Body).Decode(&course)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err = h.repo.Create(&course)
	if err != nil {
		http.Error(w, "failed to create course", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(course)
}

func NewCourseHandler(repo repository.CourseRepository) *CourseHandler {
	return &CourseHandler{repo: repo}
}

func (h *CourseHandler) GetAll(w http.ResponseWriter, r *http.Request) {

	title := r.URL.Query().Get("title")

	teacherIDStr := r.URL.Query().Get("teacher_id")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	teacherID := 0
	if teacherIDStr != "" {
		id, err := strconv.Atoi(teacherIDStr)
		if err != nil {
			http.Error(w, "invalid teacher_id", http.StatusBadRequest)
			return
		}
		teacherID = id
	}

	limit := 0
	if limitStr != "" {
		limitTemp, err := strconv.Atoi(limitStr)
		if err != nil {
			http.Error(w, "invalid limit", http.StatusBadRequest)
			return
		}
		limit = limitTemp
	}

	offset := 0
	if offsetStr != "" {
		offsetTemp, err := strconv.Atoi(offsetStr)
		if err != nil {
			http.Error(w, "invalid offset", http.StatusBadRequest)
			return
		}
		offset = offsetTemp
	}

	courses, err := h.repo.GetAll(title, teacherID, limit, offset)
	if err != nil {
		http.Error(w, "failed to get courses", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courses)
}

func (h *CourseHandler) GetCount(w http.ResponseWriter, r *http.Request) {

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
