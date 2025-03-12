package api

import (
	"context"
	"net/http"
	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/gorilla/mux"
	"github.com/rudransh-shrivastava/rudransh-backend-task/internal/db"
	"github.com/rudransh-shrivastava/rudransh-backend-task/internal/schema"
	"github.com/rudransh-shrivastava/rudransh-backend-task/internal/store"
	"github.com/rudransh-shrivastava/rudransh-backend-task/internal/utils/logger"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/option"
	"gorm.io/gorm"
)

type Server struct {
	courseStore *store.CourseStore
	userStore   *store.UserStore
	quizStore   *store.QuizStore
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
	quizStore := store.NewQuizStore(db, logger)

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
		quizStore:   quizStore,
		logger:      logger,
		db:          db,
		authClient:  authClient,
	}
}

func (s *Server) Run() {
	r := mux.NewRouter()
	r.Use(corsMiddleware)
	r.HandleFunc("/api/v1/register", s.registerUser).Methods("POST")

	api := r.PathPrefix("/api/v1").Subrouter()

	api.Use(s.authMiddleware)

	api.Handle("/courses", RBACMiddleware(s.db, schema.Student, schema.Educator, schema.Admin)(http.HandlerFunc(s.getCourses))).Methods("GET")
	api.Handle("/courses/quiz", RBACMiddleware(s.db, schema.Educator, schema.Admin)(http.HandlerFunc(s.getQuiz))).Methods("GET")
	api.Handle("/courses", RBACMiddleware(s.db, schema.Educator, schema.Admin)(http.HandlerFunc(s.createCourse))).Methods("POST")
	api.Handle("/courses", RBACMiddleware(s.db, schema.Educator, schema.Admin)(http.HandlerFunc(s.deleteCourse))).Methods("DELETE")     // TODO: fix
	api.Handle("/quiz/generate", RBACMiddleware(s.db, schema.Educator, schema.Admin)(http.HandlerFunc(s.generateQuiz))).Methods("POST") // TODO: educators can make quizzes for their courses ONLY

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
