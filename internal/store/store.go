package store

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Store struct {
	db     *gorm.DB
	logger *logrus.Logger
}

func NewStore(db *gorm.DB, logger *logrus.Logger) *Store {
	return &Store{
		db:     db,
		logger: logger,
	}
}
