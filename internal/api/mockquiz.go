package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/rudransh-shrivastava/rudransh-backend-task/internal/schema"
	"github.com/rudransh-shrivastava/rudransh-backend-task/internal/utils"
)

type Question struct {
	Question string   `json:"question"`
	Options  []string `json:"options"`
	Answer   string   `json:"answer"`
}

// list of available questions
var quizPool = []Question{
	{"What is the capital of France?", []string{"Berlin", "Madrid", "Paris", "Rome"}, "Paris"},
	{"Which planet is known as the Red Planet?", []string{"Earth", "Mars", "Jupiter", "Venus"}, "Mars"},
	{"Who wrote 'Romeo and Juliet'?", []string{"Shakespeare", "Dickens", "Hemingway", "Fitzgerald"}, "Shakespeare"},
	{"What is the largest ocean on Earth?", []string{"Atlantic", "Indian", "Arctic", "Pacific"}, "Pacific"},
	{"What is the tallest mountain in the world?", []string{"K2", "Mount Everest", "Kangchenjunga", "Lhotse"}, "Mount Everest"},
	{"What is the smallest country in the world?", []string{"Vatican City", "Monaco", "San Marino", "Liechtenstein"}, "Vatican City"},
	{"What gas do plants absorb from the atmosphere?", []string{"Oxygen", "Carbon Dioxide", "Nitrogen", "Hydrogen"}, "Carbon Dioxide"},
	{"How many continents are there?", []string{"5", "6", "7", "8"}, "7"},
	{"What is the fastest land animal?", []string{"Cheetah", "Lion", "Horse", "Kangaroo"}, "Cheetah"},
	{"Which element has the chemical symbol 'O'?", []string{"Oxygen", "Gold", "Silver", "Osmium"}, "Oxygen"},
}

// mock: returns first n questions from the quiz pool
func generateQuiz(number int) ([]Question, error) {
	if number > len(quizPool) || number <= 0 {
		return nil, fmt.Errorf("invalid quiz size: must be between 1 and %d", len(quizPool))
	}
	return quizPool[:number], nil
}

func (s *Server) generateQuiz(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		utils.WriteErrorResponse(w, "Invalid Content-Type", http.StatusBadRequest)
		return
	}
	if r.Body == nil {
		utils.WriteErrorResponse(w, "Request body is empty", http.StatusBadRequest)
		return
	}
	var req struct {
		CourseID string `json:"course_id"`
		Number   string `json:"number"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if req.CourseID == "" || req.Number == "" {
		http.Error(w, "Course ID and number are required", http.StatusBadRequest)
		return
	}
	number, err := strconv.Atoi(req.Number)
	if err != nil {
		http.Error(w, "Invalid number, must be a number", http.StatusBadRequest)
		return
	}

	courseId, err := strconv.Atoi(req.CourseID)
	if err != nil {
		utils.WriteErrorResponse(w, "Invalid course id, must be a number", http.StatusBadRequest)
		return
	}

	course, err := s.courseStore.GetCourseById(uint(courseId))
	if err != nil {
		utils.WriteErrorResponse(w, "course not found", http.StatusNotFound)
		return
	}

	questions, err := generateQuiz(number)
	if err != nil {
		utils.WriteErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonQuestions, err := json.Marshal(questions)
	if err != nil {
		utils.WriteErrorResponse(w, "failed to generate quiz", http.StatusInternalServerError)
		return
	}
	schemaQuiz := schema.Quiz{
		Course:    course,
		Questions: string(jsonQuestions),
	}
	err = s.quizStore.CreateQuiz(&schemaQuiz)
	if err != nil {
		utils.WriteErrorResponse(w, "failed to create quiz", http.StatusInternalServerError)
		return
	}
	s.logger.Debugf("Generated a quiz for course %d successfully", courseId)
}

func (s *Server) getQuiz(w http.ResponseWriter, r *http.Request) {
	queryCourseId := r.URL.Query().Get("course_id")
	queryQuizId := r.URL.Query().Get("quiz_id")
	if queryCourseId == "" {
		http.Error(w, "Course ID is required", http.StatusBadRequest)
		return
	}
	if queryQuizId == "" {
		http.Error(w, "Quiz ID is required", http.StatusBadRequest)
		return
	}

	courseId, err := strconv.Atoi(queryCourseId)
	if err != nil {
		utils.WriteErrorResponse(w, "Invalid course id, must be a number", http.StatusBadRequest)
		return
	}

	quizId, err := strconv.Atoi(queryQuizId)
	if err != nil {
		utils.WriteErrorResponse(w, "Invalid quiz id, must be a number", http.StatusBadRequest)
		return
	}

	_, err = s.courseStore.GetCourseById(uint(courseId))
	if err != nil {
		utils.WriteErrorResponse(w, "course not found", http.StatusNotFound)
		return
	}

	quiz, err := s.quizStore.GetQuizById(uint(quizId))
	if err != nil {
		utils.WriteErrorResponse(w, "quiz not found", http.StatusNotFound)
		return
	}

	user, err := s.userStore.GetUserFromContext(r.Context())
	if err != nil {
		utils.WriteErrorResponse(w, "user not found", http.StatusUnauthorized)
		return
	}

	err = s.quizStore.RegisterQuizTaken(&user, &quiz)
	utils.WriteJSONResponse(w, quiz)
}
