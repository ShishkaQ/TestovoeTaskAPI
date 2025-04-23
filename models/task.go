package models

import (
	"time"
)

// Task модель задачи
// ID, CreatedAt, Title, Description, Status, Attempts
type Task struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Attempts    int       `json:"attempts"`
}