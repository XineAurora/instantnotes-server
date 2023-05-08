package models

type Group struct {
	ID      uint   `gorm:"primaryKey"`
	Name    string `gorm:"size:255"`
	OwnerID uint   `gorm:"not null"`
	User    User   `gorm:"foreignKey:OwnerID"`
}
