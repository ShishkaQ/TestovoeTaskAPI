package storage

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"go-redis-task-service/models"
)

// InitDB открывает SQLite и выполняет миграцию
func InitDB(path string) (*gorm.DB, error) {
	if path == "" {
		path = "tasks.db"
	}
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&models.Task{}); err != nil {
		return nil, err
	}
	return db, nil
}