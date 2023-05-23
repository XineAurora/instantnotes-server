package models

type FolderLink struct {
	ParentFolderID uint   `gorm:"not null"`
	ParentFolder   Folder `gorm:"foreignKey:ParentFolderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ChildFolderID  uint   `gorm:"not null"`
	ChildFolder    Folder `gorm:"foreignKey:ChildFolderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
