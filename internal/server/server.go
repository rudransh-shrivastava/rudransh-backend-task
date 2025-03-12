package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type Server struct {
	db *gorm.DB
}

func NewServer() *Server {
	// init db
	return &Server{}
}

func (s *Server) Run() {
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/test", s.test).Methods("GET")

	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", r)
}

func (s *Server) test(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}
