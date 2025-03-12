package utils

import (
	"encoding/json"
	"net/http"
)

func WriteJSONResponse(w http.ResponseWriter, msg any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(msg)
}

func WriteErrorResponse(w http.ResponseWriter, msg string, status int) {
	w.WriteHeader(status)
	w.Write([]byte(msg))
}
