package model

import "time"

type Order struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	RestaurantID uint      `gorm:"not null" json:"restaurantId"`
	TotalSum     float64   `json:"totalSum"`
	UserID       uint      `json:"userId"`
	TableID      uint      `json:"tableId"`
	Date         time.Time `gorm:"not null" json:"date"`
	Status       string    `gorm:"not null" json:"status"`
	OrderFoods   []uint    `json:"foods"`
}

func (Order) TableName() string {
	return "orders"
}

type OrderRequest struct {
	RestaurantID uint      `gorm:"not null" json:"restaurantId"`
	TotalSum     float64   `json:"totalSum"`
	TableID      uint      `json:"tableId"`
	Date         time.Time `gorm:"not null" json:"date"`
	Status       string    `gorm:"not null" json:"status"`
	OrderFoods   []uint    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"foods"`
}

type OrderFood struct {
	ID      uint `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderID uint `gorm:"not null" json:"orderId"`
	FoodID  uint `gorm:"not null" json:"foodId"`
}

type OrderFoodResponse struct {
	FoodID uint `json:"foodId"`
	Food   Food `json:"food"`
}

type OrderResponse struct {
	ID           uint                `gorm:"primaryKey;autoIncrement" json:"id"`
	TotalSum     float64             `json:"totalSum"`
	Date         time.Time           `gorm:"not null" json:"date"`
	Status       string              `gorm:"not null" json:"status"`
	RestaurantID uint                `gorm:"not null" json:"restaurantId"`
	Restaurant   Restaurant          `gorm:"foreignKey:RestaurantID" json:"restaurant"`
	TableID      uint                `json:"tableId"`
	Table        Table               `gorm:"foreignKey:TableID" json:"table"`
	UserID       uint                `json:"userId"`
	User         UserResponse        `gorm:"foreignKey:UserID" json:"user"`
	OrderFoods   []OrderFoodResponse `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"order_foods"`
	Foods        []Food              `json:"foods"`
}
