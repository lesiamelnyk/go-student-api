package handlers

import (
	"encoding/json"
	"errors"
	"go-student-api/internal/models"
	"go-student-api/internal/repository"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockCourseRepo struct {
	shouldFail bool
}

func (m *mockCourseRepo) GetAll(title string, teacherID int, limit int, offset int) ([]models.Course, error) {

	if m.shouldFail {
		return nil, errors.New("database error")
	}

	return []models.Course{
		{
			ID:    1,
			Title: "Go Basics",
		},
	}, nil
}

func (m *mockCourseRepo) GetByID(id int) (*models.Course, error) {

	if id != 1 {
		return nil, repository.ErrNotFound
	}

	return &models.Course{
		ID:    1,
		Title: "Go Basics",
	}, nil
}

func (m *mockCourseRepo) Create(course *models.Course) error {

	if m.shouldFail {
		return errors.New("database error")
	}

	return nil
}

func (m *mockCourseRepo) Delete(id int) error {

	if id != 1 {
		return repository.ErrNotFound
	}

	return nil
}

func (m *mockCourseRepo) Update(id int, course *models.Course) error {

	if id != 1 {
		return repository.ErrNotFound
	}
	return nil
}

func (m *mockCourseRepo) GetCount() (int, error) {
	if m.shouldFail {
		return 0, errors.New("database error")
	}
	return 1, nil
}

func TestGetAllCourses_Success(t *testing.T) {

	repo := &mockCourseRepo{}
	handler := CourseHandler{repo: repo}

	req := httptest.NewRequest(http.MethodGet, "/courses", nil)
	w := httptest.NewRecorder()

	handler.GetAll(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetCourses_ResponseJSON(t *testing.T) {

	repo := &mockCourseRepo{}
	handler := CourseHandler{repo: repo}

	req := httptest.NewRequest(http.MethodGet, "/courses", nil)
	w := httptest.NewRecorder()

	handler.GetAll(w, req)

	var courses []models.Course
	err := json.Unmarshal(w.Body.Bytes(), &courses)

	assert.NoError(t, err)
	assert.Equal(t, 1, courses[0].ID)
}

func TestGetCourses_Count(t *testing.T) {

	repo := &mockCourseRepo{}
	handler := CourseHandler{repo: repo}

	req := httptest.NewRequest(http.MethodGet, "/courses", nil)
	w := httptest.NewRecorder()

	handler.GetAll(w, req)

	var courses []models.Course
	json.Unmarshal(w.Body.Bytes(), &courses)

	assert.Equal(t, 1, len(courses))
}

func TestGetCourseByID_Success(t *testing.T) {

	repo := &mockCourseRepo{}
	handler := CourseHandler{repo: repo}

	req := httptest.NewRequest(http.MethodGet, "/courses?id=1", nil)
	w := httptest.NewRecorder()

	handler.GetByID(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetCourse_NotFound(t *testing.T) {

	repo := &mockCourseRepo{}
	handler := CourseHandler{repo: repo}

	req := httptest.NewRequest(http.MethodGet, "/courses?id=999", nil)
	w := httptest.NewRecorder()

	handler.GetByID(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCreateCourse_Success(t *testing.T) {

	repo := &mockCourseRepo{}
	handler := CourseHandler{repo: repo}

	body := `{"title":"Go Advanced"}`

	req := httptest.NewRequest(http.MethodPost, "/courses", strings.NewReader(body))
	w := httptest.NewRecorder()

	handler.CreateCourse(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestCreateCourse_InvalidJSON(t *testing.T) {

	repo := &mockCourseRepo{}
	handler := CourseHandler{repo: repo}

	body := `invalid-json`

	req := httptest.NewRequest(http.MethodPost, "/courses", strings.NewReader(body))
	w := httptest.NewRecorder()

	handler.CreateCourse(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateCourse_DatabaseError(t *testing.T) {

	repo := &mockCourseRepo{shouldFail: true}
	handler := CourseHandler{repo: repo}

	body := `{"title":"Go Advanced"}`

	req := httptest.NewRequest(http.MethodPost, "/courses", strings.NewReader(body))
	w := httptest.NewRecorder()

	handler.CreateCourse(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestDeleteCourse_Success(t *testing.T) {

	repo := &mockCourseRepo{}
	handler := CourseHandler{repo: repo}

	req := httptest.NewRequest(http.MethodDelete, "/courses?id=1", nil)
	w := httptest.NewRecorder()

	handler.DeleteCourse(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestDeleteCourse_NotFound(t *testing.T) {

	repo := &mockCourseRepo{}
	handler := CourseHandler{repo: repo}

	req := httptest.NewRequest(http.MethodDelete, "/courses?id=999", nil)
	w := httptest.NewRecorder()

	handler.DeleteCourse(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
