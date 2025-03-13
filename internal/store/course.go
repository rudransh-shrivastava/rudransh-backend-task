package store

import (
	"errors"

	"github.com/rudransh-shrivastava/rudransh-backend-task/internal/schema"
)

type CourseStoreInterface interface {
	ListCourses(limit, offset int) ([]schema.Course, error)
	CreateCourse(course *schema.Course) error
	GetCourseById(id uint) (*schema.Course, error)
	DeleteCourse(course *schema.Course) error
}

type CourseStore struct {
	*Store
}

func NewCourseStore(store *Store) *CourseStore {
	return &CourseStore{Store: store}
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

func (s *CourseStore) GetCourseById(id uint) (*schema.Course, error) {
	var course schema.Course

	if err := s.db.Where("id = ?", id).First(&course).Error; err != nil {
		s.logger.Error("Failed to get course", err)
		return &course, errors.New("failed to get course")
	}

	return &course, nil
}

// shouldnt allow to delete someone else course
func (s *CourseStore) DeleteCourse(course *schema.Course) error {
	if err := s.db.Delete(course).Error; err != nil {
		s.logger.Error("Failed to delete course", err)
		return errors.New("failed to delete course")
	}
	return nil
}
