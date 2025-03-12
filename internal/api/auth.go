package api

import (
	"context"
	"encoding/json"
	"net/http"

	"firebase.google.com/go/auth"
	"github.com/rudransh-shrivastava/rudransh-backend-task/internal/schema"
)

// Handler to register new users
func (s *Server) registerUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}
	if schema.Role(req.Role) != schema.Student && schema.Role(req.Role) != schema.Educator {
		http.Error(w, "Invalid role", http.StatusBadRequest)
		return
	}

	params := (&auth.UserToCreate{}).Email(req.Email).Password(req.Password)
	userRecord, err := s.authClient.CreateUser(context.Background(), params)
	if err != nil {
		http.Error(w, "Error creating user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	s.logger.Infof("Registered a new user %+v", userRecord.UserInfo)

	dbUser := &schema.User{
		UID:   userRecord.UserInfo.UID,
		Email: userRecord.UserInfo.Email,
		Name:  userRecord.UserInfo.DisplayName,
		Role:  schema.Role(req.Role),
	}
	err = s.userStore.CreateUser(dbUser)
	if err != nil {
		http.Error(w, "Error creating user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(userRecord)
}
