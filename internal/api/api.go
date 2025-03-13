package api

import (
	"context"
	"net/http"
	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/gorilla/mux"
	"github.com/rudransh-shrivastava/rudransh-backend-task/internal/config"
	"github.com/rudransh-shrivastava/rudransh-backend-task/internal/db"
	"github.com/rudransh-shrivastava/rudransh-backend-task/internal/schema"
	"github.com/rudransh-shrivastava/rudransh-backend-task/internal/store"
	"github.com/rudransh-shrivastava/rudransh-backend-task/internal/utils/logger"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/option"
	"gorm.io/gorm"
)

type Server struct {
	courseStore store.CourseStoreInterface
	userStore   store.UserStoreInterface
	quizStore   store.QuizStoreInterface
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
	s := store.NewStore(db, logger)
	courseStore := store.NewCourseStore(s)
	userStore := store.NewUserStore(s)
	quizStore := store.NewQuizStore(s)

	// Initialize the Firebase SDK
	// the name key.json is used but we can also get it from the env vars if needed
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
	r.Use(corsMiddleware) // Use cors middleware to prevent CORS errors

	r.HandleFunc("/api/v1/register", s.registerUser).Methods("POST") // the auth endpoint
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
		Addr:         config.Envs.PublicHost + ":" + config.Envs.Port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	s.logger.Fatal(server.ListenAndServe())
}
