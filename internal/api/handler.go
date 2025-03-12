package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/rudransh-shrivastava/rudransh-backend-task/internal/schema"
	"github.com/rudransh-shrivastava/rudransh-backend-task/internal/utils"
)

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
	if course.Title == "" {
		utils.WriteErrorResponse(w, "title is required", http.StatusBadRequest)
		return
	}

	// Look up the user in your database using the UID.
	uid, _ := r.Context().Value("userID").(string)
	user, err := s.userStore.GetUserByUID(uid)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}
	course.User = user

	if err := s.courseStore.CreateCourse(&course); err != nil {
		s.logger.Error("Failed to create course", err)
		utils.WriteErrorResponse(w, "failed to create course", http.StatusInternalServerError)
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
