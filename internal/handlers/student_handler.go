package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-student-api/internal/models"
	"go-student-api/internal/repository"
	"net/http"
	"strconv"
)

func getID(r *http.Request) (int, error) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		return 0, fmt.Errorf("missing id")
	}
	return strconv.Atoi(idStr)
}

type StudentHandler struct {
	repo repository.StudentRepository
}

func (h *StudentHandler) DeleteStudent(w http.ResponseWriter, r *http.Request) {

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
		http.Error(w, "failed to delete student", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *StudentHandler) UpdateStudent(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	id, err := getID(r)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var student models.Student

	err = json.NewDecoder(r.Body).Decode(&student)
	if err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	err = h.repo.Update(id, &student)
	if errors.Is(err, repository.ErrNotFound) {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "failed to update student", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *StudentHandler) GetByID(w http.ResponseWriter, r *http.Request) {

	id, err := getID(r)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	student, err := h.repo.GetByID(id)
	if errors.Is(err, repository.ErrNotFound) {
		http.Error(w, "student not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "failed to get student", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(student)
}

func (h *StudentHandler) CreateStudent(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var student models.Student

	err := json.NewDecoder(r.Body).Decode(&student)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err = h.repo.Create(&student)
	if err != nil {
		http.Error(w, "failed to create student", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(student)
}

func NewStudentHandler(repo repository.StudentRepository) *StudentHandler {
	return &StudentHandler{repo: repo}
}

func (h *StudentHandler) GetAll(w http.ResponseWriter, r *http.Request) {

	firstName := r.URL.Query().Get("first_name")
	lastName := r.URL.Query().Get("last_name")
	email := r.URL.Query().Get("email")

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

	students, err := h.repo.GetAll(firstName, lastName, email, limit, offset)
	if err != nil {
		http.Error(w, "failed to get students", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(students)
}

func (h *StudentHandler) GetCount(w http.ResponseWriter, r *http.Request) {

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
