package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/rudransh-shrivastava/rudransh-backend-task/internal/schema"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Mock implementations
//
// MockCourseStore simulates the behavior of the CourseStore.
type MockCourseStore struct {
	Courses []schema.Course
	Err     error
}

func (m *MockCourseStore) ListCourses(limit, offset int) ([]schema.Course, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	if offset > len(m.Courses) {
		return []schema.Course{}, nil
	}
	end := offset + limit
	if end > len(m.Courses) {
		end = len(m.Courses)
	}
	return m.Courses[offset:end], nil
}

func (m *MockCourseStore) CreateCourse(course *schema.Course) error {
	if m.Err != nil {
		return m.Err
	}
	course.ID = uint(len(m.Courses) + 1)
	course.CreatedAt = time.Now()
	m.Courses = append(m.Courses, *course)
	return nil
}

func (m *MockCourseStore) GetCourseById(id uint) (*schema.Course, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	for _, c := range m.Courses {
		if c.ID == id {
			return &c, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *MockCourseStore) DeleteCourse(course *schema.Course) error {
	if m.Err != nil {
		return m.Err
	}
	for i, c := range m.Courses {
		if c.ID == course.ID {
			m.Courses = append(m.Courses[:i], m.Courses[i+1:]...)
			return nil
		}
	}
	return gorm.ErrRecordNotFound
}

// MockUserStore simulates the behavior of the UserStore
type MockUserStore struct {
	User schema.User
	Err  error
}

func (m *MockUserStore) GetUserByUID(uid string) (*schema.User, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	if m.User.UID == uid {
		return &m.User, nil
	}
	return nil, errors.New("user not found")
}

// Unused in the tests, but required to implement the interface
func (m *MockUserStore) GetUserFromContext(ctx context.Context) (*schema.User, error) {
	return &m.User, nil
}
func (m *MockUserStore) CreateUser(user *schema.User) error {
	return nil
}

// TestServer setup

// TestServer embeds Server and includes the mocks
type TestServer struct {
	*Server
	mockCourseStore *MockCourseStore
	mockUserStore   *MockUserStore
}

func newTestServer() *TestServer {
	logger := logrus.New()
	mockCourseStore := &MockCourseStore{
		Courses: []schema.Course{
			{ID: 1, Title: "Course 1", CreatedAt: time.Now()},
			{ID: 2, Title: "Course 2", CreatedAt: time.Now()},
			{ID: 3, Title: "Course 3", CreatedAt: time.Now()},
		},
	}
	mockUserStore := &MockUserStore{
		User: schema.User{
			ID:    1,
			UID:   "test-uid",
			Email: "test@example.com",
			Name:  "Test User",
			Role:  schema.Student,
		},
	}
	s := &Server{
		courseStore: mockCourseStore,
		userStore:   mockUserStore,
		logger:      logger,
	}
	return &TestServer{
		Server:          s,
		mockCourseStore: mockCourseStore,
		mockUserStore:   mockUserStore,
	}
}

// Tests for getCourses
func TestGetCourses_Success(t *testing.T) {
	ts := newTestServer()

	// Request with valid limit and offset.
	req := httptest.NewRequest("GET", "/api/v1/courses?limit=2&offset=1", nil)
	rr := httptest.NewRecorder()

	ts.getCourses(rr, req)
	res := rr.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200 OK, got %d", res.StatusCode)
	}

	var courses []schema.Course
	body, _ := io.ReadAll(res.Body)
	if err := json.Unmarshal(body, &courses); err != nil {
		t.Errorf("failed to unmarshal response: %v", err)
	}
	if len(courses) != 2 {
		t.Errorf("expected 2 courses, got %d", len(courses))
	}
}

func TestGetCourses_InvalidLimit(t *testing.T) {
	ts := newTestServer()

	req := httptest.NewRequest("GET", "/api/v1/courses?limit=abc&offset=0", nil)
	rr := httptest.NewRecorder()

	ts.getCourses(rr, req)
	res := rr.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400 Bad Request, got %d", res.StatusCode)
	}
}

// Tsts for createCourse
func TestCreateCourse_Success(t *testing.T) {
	ts := newTestServer()

	courseData := map[string]string{
		"title": "New Course",
	}
	courseJSON, _ := json.Marshal(courseData)
	req := httptest.NewRequest("POST", "/api/v1/courses", bytes.NewBuffer(courseJSON))
	req.Header.Set("Content-Type", "application/json")
	// Inject userID into the request context
	ctx := context.WithValue(req.Context(), "userID", "test-uid")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	ts.postCourse(rr, req)
	res := rr.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200 OK, got %d", res.StatusCode)
	}

	var course schema.Course
	body, _ := io.ReadAll(res.Body)
	if err := json.Unmarshal(body, &course); err != nil {
		t.Errorf("failed to unmarshal response: %v", err)
	}
	if course.Title != "New Course" {
		t.Errorf("expected course title 'New Course', got '%s'", course.Title)
	}
}

func TestCreateCourse_InvalidContentType(t *testing.T) {
	ts := newTestServer()

	req := httptest.NewRequest("POST", "/api/v1/courses", nil)
	req.Header.Set("Content-Type", "text/plain")
	rr := httptest.NewRecorder()

	ts.postCourse(rr, req)
	res := rr.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400 Bad Request, got %d", res.StatusCode)
	}
}

// Tests for deleteCourse

func TestDeleteCourse_Success(t *testing.T) {
	ts := newTestServer()

	// Prepare a delete request for an existing course (ID 1)
	course := schema.Course{ID: 1, Title: "Course 1"}
	courseJSON, _ := json.Marshal(course)
	req := httptest.NewRequest("DELETE", "/api/v1/courses", bytes.NewBuffer(courseJSON))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	ts.deleteCourse(rr, req)
	res := rr.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200 OK, got %d", res.StatusCode)
	}

	// Verify that the course was deleted
	courses, _ := ts.mockCourseStore.ListCourses(10, 0)
	for _, c := range courses {
		if c.ID == 1 {
			t.Errorf("course with ID 1 should have been deleted")
		}
	}
}

func TestDeleteCourse_CourseNotFound(t *testing.T) {
	ts := newTestServer()

	// Prepare a delete request for a non-existent course (ID 999)
	course := schema.Course{ID: 999}
	courseJSON, _ := json.Marshal(course)
	req := httptest.NewRequest("DELETE", "/api/v1/courses", bytes.NewBuffer(courseJSON))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	ts.deleteCourse(rr, req)
	res := rr.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusNotFound {
		t.Errorf("expected status 404 Not Found, got %d", res.StatusCode)
	}
}
