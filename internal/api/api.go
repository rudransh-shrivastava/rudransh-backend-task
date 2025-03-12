package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rudransh-shrivastava/rudransh-backend-task/internal/db"
	"github.com/rudransh-shrivastava/rudransh-backend-task/internal/middleware"
	"github.com/rudransh-shrivastava/rudransh-backend-task/internal/schema"
	"github.com/rudransh-shrivastava/rudransh-backend-task/internal/store"
	"github.com/rudransh-shrivastava/rudransh-backend-task/internal/utils"
	"github.com/rudransh-shrivastava/rudransh-backend-task/internal/utils/logger"
	"github.com/sirupsen/logrus"
)

type Server struct {
	courseStore *store.CourseStore
	logger      *logrus.Logger
}

func NewServer() *Server {
	logger := logger.NewLogger()
	db, err := db.NewDB()
	if err != nil {
		logger.Fatalf("failed to connect to database: %v", err)
	}
	courseStore := store.NewCourseStore(db, logger)
	return &Server{
		courseStore: courseStore,
		logger:      logger,
	}
}

func (s *Server) Run() {
	r := mux.NewRouter()

	rateLimitMiddleware := middleware.NewRateLimitMiddleware()
	// TODO: implement auth
	r.HandleFunc("/api/v1/courses", s.getCourses).Methods("GET")
	r.HandleFunc("/api/v1/courses", s.createCourse).Methods("POST")
	// r.HandleFunc("/api/v1/courses", s.updateCourse).Methods("UPDATE")
	// r.HandleFunc("/api/v1/courses", s.deleteCourse).Methods("DELETE")

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
	// TODO: get the limit and offset from serach query parameters
	limit := 10
	offset := 0

	courses, err := s.courseStore.ListCourses(limit, offset)
	if err != nil {
		s.logger.Error("Failed to list courses", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	utils.WriteJSONResponse(w, courses)
}

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
