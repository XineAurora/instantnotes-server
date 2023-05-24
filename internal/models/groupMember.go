package models

type GroupMember struct {
	ID      uint  `gorm:"primaryKey" json:"id"`
	GroupID uint  `gorm:"not null" json:"groupId"`
	Group   Group `gorm:"foreignKey:GroupID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	UserID  uint  `gorm:"not null" json:"userId"`
	User    User  `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}
