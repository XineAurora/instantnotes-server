package models

type Folder struct {
	ID      uint   `gorm:"primarykey"`
	Name    string `gorm:"size:255"`
	UserID  uint   `gorm:"not null"`
	User    User   `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	GroupID *uint
	Group   Group `gorm:"foreignKey:GroupID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
