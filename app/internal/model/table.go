package model

type Table struct {
	ID           uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name         string `gorm:"not null" json:"name"`
	Type         string `gorm:"not null" json:"type"`
	Description  string `json:"description"`
	Capacity     int    `gorm:"not null" json:"capacity"`
	PhotoID      uint   `json:"photo_id,omitempty"`
	Photo        Photo  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"photo,omitempty"`
	RestaurantID uint   `gorm:"not null" json:"restaurant_id"`
}
