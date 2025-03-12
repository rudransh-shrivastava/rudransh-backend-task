package store

import (
	"context"
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

func (s *UserStore) GetUserByUID(uid string) (schema.User, error) {
	var user schema.User

	if err := s.db.Where("uid = ?", uid).First(&user).Error; err != nil {
		s.logger.Error("Failed to get user", err)
		return user, errors.New("failed to get user")
	}

	return user, nil
}

func (s *UserStore) GetUserFromContext(ctx context.Context) (schema.User, error) {
	uid := ctx.Value("userID").(string)
	user, err := s.GetUserByUID(uid)
	if err != nil {
		return user, err
	}
	return user, nil
}
