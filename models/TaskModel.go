package models

import (
	"time"

	"gorm.io/gorm"
)

// Task represents a task in a project. A Task may be assigned to a User.
type Task struct {
	gorm.Model
	Title       string     `gorm:"not null"`
	Description string
	Status      string     `gorm:"default:'todo'"`
	Priority    int
	DueDate     *time.Time
	ProjectID   uint
	Project     Project
	AssigneeID  *uint
	Assignee    *User
}
