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

	err = database.AutoMigrate(&schema.User{}, &schema.Course{})

	if err != nil {
		return nil, err
	}
	return database, nil
}
