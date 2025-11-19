package models

import "gorm.io/gorm"

// Project represents a project that contains many tasks and belongs to an owner (User)
type Project struct {
	gorm.Model
	Name        string `gorm:"not null"`
	Description string
	OwnerID     uint
	Owner       User
	Tasks       []Task `gorm:"foreignKey:ProjectID"`
}
