package postgre

import (
	"context"
	"github.com/alibekabdrakhman1/orynal/internal/model"
	"gorm.io/gorm"
)

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{
		DB: db,
	}
}

type OrderRepository struct {
	DB *gorm.DB
}

func (r *OrderRepository) CreateOrder(ctx context.Context, order *model.Order) (*model.OrderResponse, error) {
	var orderResponse model.OrderResponse

	tx := r.DB.WithContext(ctx).Begin()

	if err := tx.Table("orders").Create(order).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Table("orders").Preload("Restaurant").Preload("Table").Where("id = ?", order.ID).First(&orderResponse).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	for _, foodID := range order.OrderFoods {
		orderFood := model.OrderFood{
			OrderID: order.ID,
			FoodID:  foodID,
		}
		if err := tx.Table("order_foods").Create(&orderFood).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	var rest model.Restaurant
	if err := r.DB.WithContext(ctx).Table("restaurants").Where("id = ?", order.RestaurantID).First(&rest).Error; err != nil {
		return nil, err
	}
	orderResponse.Restaurant = rest

	var table model.Table
	if err := r.DB.WithContext(ctx).Table("tables").Where("id = ?", order.TableID).First(&table).Error; err != nil {
		return nil, err
	}
	orderResponse.Table = table

	var orderFoods []model.OrderFoodResponse
	for _, foodID := range order.OrderFoods {
		var food model.Food
		if err := r.DB.WithContext(ctx).Table("foods").Where("id = ?", foodID).First(&food).Error; err != nil {
			continue
		}
		orderFoods = append(orderFoods, model.OrderFoodResponse{
			FoodID: foodID,
			Food:   food,
		})
	}
	orderResponse.OrderFoods = orderFoods

	return &orderResponse, nil
}

func (r *OrderRepository) DeleteOrder(ctx context.Context, id uint) error {
	if err := r.DB.WithContext(ctx).Table("orders").Delete(&model.Order{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (r *OrderRepository) UpdateOrder(ctx context.Context, order *model.Order) (*model.OrderResponse, error) {
	var or model.Order
	if err := r.DB.WithContext(ctx).Table("orders").First(&or, order.ID).Error; err != nil {
		return nil, err
	}

	if err := r.DB.WithContext(ctx).Model(&or).Table("orders").Updates(order).Error; err != nil {
		return nil, err
	}

	return r.GetOrder(ctx, order.ID)
}

func (r *OrderRepository) GetOrder(ctx context.Context, id uint) (*model.OrderResponse, error) {
	var order model.OrderResponse

	if err := r.DB.WithContext(ctx).
		Table("orders").
		Preload("Restaurant").
		Preload("Table").
		Preload("User").
		Where("id = ?", id).
		First(&order).Error; err != nil {
		return nil, err
	}

	var rest model.Restaurant
	if err := r.DB.WithContext(ctx).Table("restaurants").Where("id = ?", order.RestaurantID).First(&rest).Error; err != nil {
		return nil, err
	}

	var icon model.Photo
	if err := r.DB.Table("photos").Where("id = ?", rest.IconID).First(&icon).Error; err != nil {
		return nil, err
	}

	rest.Icon = icon

	var services []model.Service
	if err := r.DB.Raw("SELECT services.* FROM services JOIN restaurant_service ON services.id = restaurant_service.service_id WHERE restaurant_service.restaurant_id = ?", rest.ID).Scan(&services).Error; err != nil {
		return nil, err
	}
	rest.Services = services

	order.Restaurant = rest

	var table model.Table
	if err := r.DB.WithContext(ctx).Table("tables").Where("id = ?", order.TableID).First(&table).Error; err != nil {
		return nil, err
	}
	order.Table = table

	var orderFoods []model.OrderFood

	if err := r.DB.WithContext(ctx).
		Table("order_foods").
		Where("order_id = ?", order.ID).
		Find(&orderFoods).Error; err != nil {
		return nil, err
	}

	var foods []model.Food
	for _, orderFood := range orderFoods {
		var food model.Food
		if err := r.DB.WithContext(ctx).Table("foods").Where("id = ?", orderFood.FoodID).First(&food).Error; err != nil {
			return nil, err
		}

		var photo model.Photo
		if err := r.DB.Table("photos").Where("id = ?", food.PhotoID).First(&photo).Error; err != nil {
			return nil, err
		}

		food.Photo = photo

		foods = append(foods, food)
	}

	var user model.UserResponse
	if err := r.DB.Table("users").Where("id = ?", order.UserID).First(&user).Error; err != nil {
		return nil, err
	}

	order.User = user
	order.Foods = foods

	return &order, nil
}

func (r *OrderRepository) GetAllOrders(ctx context.Context, userID uint, params *model.Params) (*model.ListResponse, error) {
	var orders []model.OrderResponse
	var totalItems int64

	query := r.DB.WithContext(ctx).Table("orders").Where("user_id = ?", userID)

	if err := query.Model(&model.OrderResponse{}).Count(&totalItems).Error; err != nil {
		return nil, err
	}

	query = query.Limit(params.Limit).Offset(params.Offset)

	if params.Order != nil && params.SortVector != nil {
		query.Order(params.Order.(string) + " " + params.SortVector.(string))
	}

	if err := query.Find(&orders).Error; err != nil {
		return nil, err
	}

	for i := 0; i < len(orders); i++ {
		var rest model.Restaurant
		if err := r.DB.WithContext(ctx).Table("restaurants").Where("id = ?", orders[i].RestaurantID).First(&rest).Error; err != nil {
			return nil, err
		}
		var services []model.Service
		if err := r.DB.Raw("SELECT services.* FROM services JOIN restaurant_service ON services.id = restaurant_service.service_id WHERE restaurant_service.restaurant_id = ?", rest.ID).Scan(&services).Error; err != nil {
			return nil, err
		}
		rest.Services = services

		orders[i].Restaurant = rest

		var table model.Table
		if err := r.DB.WithContext(ctx).Table("tables").Where("id = ?", orders[i].TableID).First(&table).Error; err != nil {
			return nil, err
		}
		orders[i].Table = table
	}

	return &model.ListResponse{
		Items:        orders,
		ItemsPerPage: params.Limit,
		PageIndex:    params.PageIndex,
		TotalItems:   int(totalItems),
	}, nil
}

func (r *OrderRepository) GetRestaurantOrders(ctx context.Context, restaurantID uint, params *model.Params) (*model.ListResponse, error) {
	var orders []model.OrderResponse
	var totalItems int64

	query := r.DB.WithContext(ctx).Table("orders").Where("restaurant_id = ?", restaurantID)

	if err := query.Model(&model.OrderResponse{}).Count(&totalItems).Error; err != nil {
		return nil, err
	}

	query = query.Limit(params.Limit).Offset(params.Offset)

	if params.Order != nil && params.SortVector != nil {
		query.Order(params.Order.(string) + " " + params.SortVector.(string))
	}

	if err := query.Find(&orders).Error; err != nil {
		return nil, err
	}

	for i := 0; i < len(orders); i++ {
		var rest model.Restaurant
		if err := r.DB.WithContext(ctx).Table("restaurants").Where("id = ?", orders[i].RestaurantID).First(&rest).Error; err != nil {
			return nil, err
		}
		var services []model.Service
		if err := r.DB.Raw("SELECT services.* FROM services JOIN restaurant_service ON services.id = restaurant_service.service_id WHERE restaurant_service.restaurant_id = ?", rest.ID).Scan(&services).Error; err != nil {
			return nil, err
		}
		rest.Services = services

		orders[i].Restaurant = rest

		var table model.Table
		if err := r.DB.WithContext(ctx).Table("tables").Where("id = ?", orders[i].TableID).First(&table).Error; err != nil {
			return nil, err
		}
		orders[i].Table = table
	}

	return &model.ListResponse{
		Items:        orders,
		ItemsPerPage: params.Limit,
		PageIndex:    params.PageIndex,
		TotalItems:   int(totalItems),
	}, nil
}
