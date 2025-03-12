package api

import (
	"context"
	"encoding/json"
	"net/http"

	"firebase.google.com/go/auth"
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

	// TODO: Create the user in the db
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(userRecord)
}
