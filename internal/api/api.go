package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/gorilla/mux"
	"github.com/rudransh-shrivastava/rudransh-backend-task/internal/db"
	"github.com/rudransh-shrivastava/rudransh-backend-task/internal/schema"
	"github.com/rudransh-shrivastava/rudransh-backend-task/internal/store"
	"github.com/rudransh-shrivastava/rudransh-backend-task/internal/utils"
	"github.com/rudransh-shrivastava/rudransh-backend-task/internal/utils/logger"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/option"
	"gorm.io/gorm"
)

type Server struct {
	courseStore *store.CourseStore
	userStore   *store.UserStore
	logger      *logrus.Logger
	db          *gorm.DB
	authClient  *auth.Client
}

func NewServer() *Server {
	logger := logger.NewLogger()
	db, err := db.NewDB()
	if err != nil {
		logger.Fatalf("failed to connect to database: %v", err)
	}
	courseStore := store.NewCourseStore(db, logger)
	userStore := store.NewUserStore(db, logger)

	// Initialize Firebase SDK
	opt := option.WithCredentialsFile("key.json")
	fbApp, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		logger.Fatalf("error initializing firebase app: %v", err)
	}
	authClient, err := fbApp.Auth(context.Background())
	if err != nil {
		logger.Fatalf("error initializing firebase auth: %v", err)
	}

	return &Server{
		courseStore: courseStore,
		userStore:   userStore,
		logger:      logger,
		db:          db,
		authClient:  authClient,
	}
}

func (s *Server) Run() {
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/register", s.registerUser).Methods("POST")

	api := r.PathPrefix("/api/v1").Subrouter()
	api.Use(s.authMiddleware)

	api.Handle("/courses", RBACMiddleware(s.db, schema.Student, schema.Educator, schema.Admin)(http.HandlerFunc(s.getCourses))).Methods("GET")
	api.Handle("/courses", RBACMiddleware(s.db, schema.Educator, schema.Admin)(http.HandlerFunc(s.createCourse))).Methods("POST")
	api.Handle("/courses", RBACMiddleware(s.db, schema.Educator, schema.Admin)(http.HandlerFunc(s.updateCourse))).Methods("UPDATE") // TODO: fix
	api.Handle("/courses", RBACMiddleware(s.db, schema.Educator, schema.Admin)(http.HandlerFunc(s.deleteCourse))).Methods("DELETE") // TODO: fix

	rateLimitMiddleware := s.newRateLimitMiddleware()
	r.Use(rateLimitMiddleware)

	s.logger.Info("Server is running on port 8080")
	server := &http.Server{
		Handler:      r,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	s.logger.Fatal(server.ListenAndServe())
}

func (s *Server) getCourses(w http.ResponseWriter, r *http.Request) {
	s.logger.Info("/GET /api/v1/courses")

	limitQuery := r.URL.Query().Get("limit")
	offsetQuery := r.URL.Query().Get("offset")
	if limitQuery == "" {
		limitQuery = "10"
	}
	if offsetQuery == "" {
		offsetQuery = "0"
	}

	limit, err := strconv.Atoi(limitQuery)
	if err != nil {
		utils.WriteErrorResponse(w, "Invalid limit, limit must be a number", http.StatusBadRequest)
		return
	}
	offset, err := strconv.Atoi(offsetQuery)
	if err != nil {
		utils.WriteErrorResponse(w, "Invalid offset, offset must be a number", http.StatusBadRequest)
		return
	}
	courses, err := s.courseStore.ListCourses(limit, offset)
	if err != nil {
		utils.WriteErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	utils.WriteJSONResponse(w, courses)
}

// TODO: handle bad input
func (s *Server) createCourse(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		utils.WriteErrorResponse(w, "Invalid Content-Type", http.StatusBadRequest)
		return
	}
	if r.Body == nil {
		utils.WriteErrorResponse(w, "Request body is empty", http.StatusBadRequest)
		return
	}

	var course schema.Course

	err := json.NewDecoder(r.Body).Decode(&course)
	if err != nil {

		utils.WriteErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if err := s.courseStore.CreateCourse(&course); err != nil {
		s.logger.Error("Failed to create course", err)
		utils.WriteErrorResponse(w, "failed to create course", http.StatusInternalServerError)
		return
	}

	utils.WriteJSONResponse(w, course)
}

func (s *Server) updateCourse(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		utils.WriteErrorResponse(w, "Invalid Content-Type", http.StatusBadRequest)
		return
	}
	if r.Body == nil {
		utils.WriteErrorResponse(w, "Request body is empty", http.StatusBadRequest)
		return
	}

	var course schema.Course

	err := json.NewDecoder(r.Body).Decode(&course)
	if err != nil {
		utils.WriteErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if course.ID == 0 {
		utils.WriteErrorResponse(w, "id is required", http.StatusBadRequest)
		return
	}
	if course.Title == "" {
		utils.WriteErrorResponse(w, "title is required", http.StatusBadRequest)
		return
	}

	_, err = s.courseStore.GetCourseById(course.ID)
	if err != nil {
		utils.WriteErrorResponse(w, "course not found", http.StatusNotFound)
		return
	}

	if err := s.courseStore.UpdateCourse(&course); err != nil {
		s.logger.Error("Failed to update course", err)
		utils.WriteErrorResponse(w, "failed to update course", http.StatusInternalServerError)
		return
	}

	utils.WriteJSONResponse(w, course)
}

func (s *Server) deleteCourse(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		utils.WriteErrorResponse(w, "Invalid Content-Type", http.StatusBadRequest)
		return
	}
	if r.Body == nil {
		utils.WriteErrorResponse(w, "Request body is empty", http.StatusBadRequest)
		return
	}

	var course schema.Course

	err := json.NewDecoder(r.Body).Decode(&course)
	if err != nil {
		utils.WriteErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if course.ID == 0 {
		utils.WriteErrorResponse(w, "id is required", http.StatusBadRequest)
		return
	}

	_, err = s.courseStore.GetCourseById(course.ID)
	if err != nil {
		utils.WriteErrorResponse(w, "course not found", http.StatusNotFound)
		return
	}

	if err := s.courseStore.DeleteCourse(&course); err != nil {
		s.logger.Error("Failed to delete course", err)
		utils.WriteErrorResponse(w, "failed to delete course", http.StatusInternalServerError)
		return
	}

	utils.WriteJSONResponse(w, course)
}
