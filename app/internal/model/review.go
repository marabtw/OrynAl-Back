package model

import "time"

type RestaurantReview struct {
	ID           uint         `gorm:"primaryKey;autoIncrement" json:"id"`
	Stars        int          `gorm:"not null" json:"stars"`
	Description  string       `json:"description"`
	UserID       uint         `gorm:"not null" json:"user_id"`
	RestaurantID uint         `gorm:"not null" json:"restaurant_id"`
	Date         time.Time    `gorm:"default:CURRENT_TIMESTAMP" json:"date"`
	User         UserResponse `gorm:"foreignKey:UserID" json:"user"`
}
