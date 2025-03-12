package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rudransh-shrivastava/rudransh-backend-task/internal/db"
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
	// TODO: implement auth
	r.HandleFunc("/api/v1/courses", s.getCourses).Methods("GET")

	s.logger.Info("Server is running on port 8080")
	http.ListenAndServe(":8080", r)
}

func (s *Server) getCourses(w http.ResponseWriter, r *http.Request) {
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
