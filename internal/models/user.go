package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name         string `gorm:"not null;size:255" json:"name"`
	PasswordHash string `gorm:"not null;size:255" json:"-"`
	Email        string `gorm:"not null;unique;" json:"email"`
}
