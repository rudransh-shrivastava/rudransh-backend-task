package store

import (
	"errors"

	"github.com/rudransh-shrivastava/rudransh-backend-task/internal/schema"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CourseStore struct {
	db     *gorm.DB
	logger *logrus.Logger
}

func NewCourseStore(db *gorm.DB, logger *logrus.Logger) *CourseStore {
	return &CourseStore{
		db:     db,
		logger: logger,
	}
}

func (s *CourseStore) CreateCourse(course *schema.Course) error {
	if err := s.db.Create(course).Error; err != nil {
		s.logger.Error("Failed to create course", err)
		return errors.New("failed to create course")
	}
	return nil
}

func (s *CourseStore) ListCourses(limit, offset int) ([]schema.Course, error) {
	var courses []schema.Course

	if err := s.db.Order("created_at desc").Limit(limit).Offset(offset).Find(&courses).Error; err != nil {
		s.logger.Error("Failed to list courses", err)
		return nil, errors.New("database error")
	}

	return courses, nil
}
