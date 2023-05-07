package models

type FoldersLink struct {
	ParentFolderID uint   `gorm:"not null"`
	ParentFolder   Folder `gorm:"foreignKey:ParentFolderID"`
	ChildFolderID  uint   `gorm:"not null"`
	ChildFolder    Folder `gorm:"foreignKey:ChildFolderID"`
}
