package models

type Folder struct {
	ID      uint   `gorm:"primarykey" json:"id"`
	Name    string `gorm:"size:255" json:"name"`
	UserID  uint   `gorm:"not null" json:"userId"`
	User    User   `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	GroupID *uint  `json:"groupId"`
	Group   Group  `gorm:"foreignKey:GroupID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}
