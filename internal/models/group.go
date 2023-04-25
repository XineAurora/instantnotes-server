package models

type Group struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"size:255"`
}
