package store

import (
	"errors"

	"github.com/rudransh-shrivastava/rudransh-backend-task/internal/schema"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserStore struct {
	db     *gorm.DB
	logger *logrus.Logger
}

func NewUserStore(db *gorm.DB, logger *logrus.Logger) *UserStore {
	return &UserStore{
		db:     db,
		logger: logger,
	}
}

func (us *UserStore) CreateUser(user *schema.User) error {
	if err := us.db.Create(user).Error; err != nil {
		us.logger.Error("Failed to create user", err)
		return errors.New("failed to create user")
	}
	return nil
}
