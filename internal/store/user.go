package store

import (
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
