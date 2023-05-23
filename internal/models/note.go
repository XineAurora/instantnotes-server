package models

import (
	"gorm.io/gorm"
)

type Note struct {
	gorm.Model
	Title    string `gorm:"size:255"`
	Content  string `gorm:"size:65535"`
	UserID   uint   `gorm:"not null"`
	User     User   `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	FolderID *uint
	Folder   Folder `gorm:"foreignKey:FolderID;default:NULL;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	GroupID  *uint
	Group    Group `gorm:"foreignKey:GroupID;default:NULL;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
