package models

type FolderLink struct {
	ParentFolderID uint   `gorm:"not null" json:"parentFolderId"`
	ParentFolder   Folder `gorm:"foreignKey:ParentFolderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	ChildFolderID  uint   `gorm:"not null" json:"childFolderId"`
	ChildFolder    Folder `gorm:"foreignKey:ChildFolderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}
