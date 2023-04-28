package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name         string `gorm:"not null;size:255"`
	PasswordHash string `gorm:"not null;size:255"`
	Email        string `gorm:"not null;unique;"`
}
