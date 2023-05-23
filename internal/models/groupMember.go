package models

type GroupMember struct {
	ID      uint  `gorm:"primaryKey"`
	GroupID uint  `gorm:"not null"`
	Group   Group `gorm:"foreignKey:GroupID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserID  uint  `gorm:"not null"`
	User    User  `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
