package handlers

import (
	"encoding/json"
	"errors"
	"go-student-api/internal/models"
	"go-student-api/internal/repository"
	"net/http"
	"strconv"
)

type TeacherHandler struct {
	repo repository.TeacherRepository
}

func (h *TeacherHandler) DeleteTeacher(w http.ResponseWriter, r *http.Request) {

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
		http.Error(w, "failed to delete teacher", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *TeacherHandler) UpdateTeacher(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	id, err := getID(r)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var teacher models.Teacher

	err = json.NewDecoder(r.Body).Decode(&teacher)
	if err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	err = h.repo.Update(id, &teacher)
	if errors.Is(err, repository.ErrNotFound) {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "failed to update teacher", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *TeacherHandler) GetByID(w http.ResponseWriter, r *http.Request) {

	id, err := getID(r)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	teacher, err := h.repo.GetByID(id)
	if errors.Is(err, repository.ErrNotFound) {
		http.Error(w, "teacher not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "failed to get teacher", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teacher)
}

func (h *TeacherHandler) CreateTeacher(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var teacher models.Teacher

	err := json.NewDecoder(r.Body).Decode(&teacher)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err = h.repo.Create(&teacher)
	if err != nil {
		http.Error(w, "failed to create teacher", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(teacher)
}

func NewTeacherHandler(repo repository.TeacherRepository) *TeacherHandler {
	return &TeacherHandler{repo: repo}
}

func (h *TeacherHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	firstName := r.URL.Query().Get("first_name")
	lastName := r.URL.Query().Get("last_name")
	department := r.URL.Query().Get("department")

	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

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

	teachers, err := h.repo.GetAll(firstName, lastName, department, limit, offset)
	if err != nil {
		http.Error(w, "failed to get teachers", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teachers)
}

func (h *TeacherHandler) GetCount(w http.ResponseWriter, r *http.Request) {

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
