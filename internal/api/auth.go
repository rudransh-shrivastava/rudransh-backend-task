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
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
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
		Role:  schema.Student,
	}
	err = s.userStore.CreateUser(dbUser)
	if err != nil {
		http.Error(w, "Error creating user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(userRecord)
}
