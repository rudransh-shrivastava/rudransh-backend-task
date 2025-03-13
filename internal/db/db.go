package db

import (
	"github.com/glebarez/sqlite"
	"github.com/rudransh-shrivastava/rudransh-backend-task/internal/schema"
	"gorm.io/gorm"
)

func NewDB() (*gorm.DB, error) {
	dbName := "db.sqlite3"
	database, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{
		PrepareStmt: true,
	})
	if err != nil {
		return nil, err
	}
	sqlDB, err := database.DB()
	sqlDB.Exec("PRAGMA foreign_keys = ON")

	// Migrate our schemas
	err = database.AutoMigrate(&schema.User{}, &schema.Course{}, &schema.Quiz{}, &schema.QuizzesTaken{})

	if err != nil {
		return nil, err
	}
	return database, nil
}
