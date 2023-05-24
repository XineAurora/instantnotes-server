package models

type Group struct {
	ID      uint   `gorm:"primaryKey" json:"id"`
	Name    string `gorm:"size:255" json:"name"`
	OwnerID uint   `gorm:"not null" json:"ownerId"`
	User    User   `gorm:"foreignKey:OwnerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}
