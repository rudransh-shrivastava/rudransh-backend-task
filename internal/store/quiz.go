package store

import (
	"errors"

	"github.com/rudransh-shrivastava/rudransh-backend-task/internal/schema"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type QuizStore struct {
	db     *gorm.DB
	logger *logrus.Logger
}

func NewQuizStore(db *gorm.DB, logger *logrus.Logger) *QuizStore {
	return &QuizStore{
		db:     db,
		logger: logger,
	}
}

func (qs *QuizStore) CreateQuiz(quiz *schema.Quiz) error {
	if err := qs.db.Create(quiz).Error; err != nil {
		qs.logger.Error("Failed to create quiz", err)
		return errors.New("failed to create quiz")
	}
	return nil
}

func (qs *QuizStore) GetQuizById(id uint) (schema.Quiz, error) {
	var quiz schema.Quiz

	if err := qs.db.Where("id = ?", id).First(&quiz).Error; err != nil {
		qs.logger.Error("Failed to get quiz", err)
		return quiz, errors.New("failed to get quiz")
	}

	return quiz, nil
}

func (qs *QuizStore) RegisterQuizTaken(user *schema.User, quiz *schema.Quiz) error {
	if err := qs.db.Create(&schema.QuizzesTaken{User: *user, Quiz: *quiz}).Error; err != nil {
		qs.logger.Error("Failed to take quiz", err)
		return errors.New("failed to take quiz")
	}
	return nil
}
