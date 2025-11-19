package models

import "gorm.io/gorm"

// User represents an application user
type User struct {
	gorm.Model
	Email    string    `gorm:"unique;not null"`
	Password string    `gorm:"not null"`
	Projects []Project `gorm:"foreignKey:OwnerID"`
	Tasks    []Task    `gorm:"foreignKey:AssigneeID"`
}
