package models

type GroupMember struct {
	ID      uint  `gorm:"primaryKey"`
	GroupID *uint `gorm:"not null"`
	Group   Group `gorm:"foreignKey:GroupID"`
	UserID  *uint `gorm:"not null"`
	User    User  `gorm:"foreignKey:UserID"`
}
