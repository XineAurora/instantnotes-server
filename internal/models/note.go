package models

import (
	"gorm.io/gorm"
)

type Note struct {
	gorm.Model
	Title    string `gorm:"size:255" json:"title"`
	Content  string `gorm:"size:65535" json:"content"`
	UserID   uint   `gorm:"not null" json:"userId"`
	User     User   `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	FolderID *uint  `json:"folderId"`
	Folder   Folder `gorm:"foreignKey:FolderID;default:NULL;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	GroupID  *uint  `json:"groupId"`
	Group    Group  `gorm:"foreignKey:GroupID;default:NULL;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}
