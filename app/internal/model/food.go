package model

type Food struct {
	ID           uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name         string  `gorm:"not null" json:"name"`
	Type         string  `gorm:"not null" json:"type"`
	Description  string  `json:"description"`
	Price        float64 `gorm:"not null" json:"price"`
	Available    bool    `gorm:"not null" json:"available"`
	PhotoID      uint    `json:"photo_id,omitempty"`
	Photo        Photo   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"photo,omitempty"`
	RestaurantID uint    `gorm:"not null" json:"restaurantId"`
}
