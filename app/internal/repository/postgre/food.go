package postgre

import (
	"context"
	"errors"
	"fmt"
	"github.com/alibekabdrakhman1/orynal/internal/model"
	"gorm.io/gorm"
)

func NewFoodRepository(db *gorm.DB) *FoodRepository {
	return &FoodRepository{
		DB: db,
	}
}

type FoodRepository struct {
	DB *gorm.DB
}

func (r *FoodRepository) GetMenuCategories(ctx context.Context, restaurantID uint) ([]string, error) {
	var types []string

	if err := r.DB.WithContext(ctx).Table("foods").
		Select("DISTINCT type").
		Where("restaurant_id = ?", restaurantID).
		Pluck("type", &types).Error; err != nil {
		return nil, err
	}

	return types, nil
}

func (r *FoodRepository) GetRestaurantMenu(ctx context.Context, restaurantID uint, params *model.Params) (*model.ListResponse, error) {
	var foods []model.Food
	var totalItems int64

	countQuery := r.DB.WithContext(ctx).
		Model(&model.Food{}).
		Table("foods").
		Where("restaurant_id = ?", restaurantID)

	if params.Query != "" {
		countQuery = countQuery.Where("LOWER(type) = LOWER(?)", params.Query)
	}

	if err := countQuery.Count(&totalItems).Error; err != nil {
		return nil, err
	}

	fmt.Println(totalItems)

	if int64(params.Offset) >= totalItems {
		return nil, errors.New("offset exceeds total items")
	}

	query := r.DB.WithContext(ctx).
		Table("foods").
		Select("foods.*").
		Where("restaurant_id = ?", restaurantID)

	if params.Query != "" {
		query = query.Where("LOWER(type) = LOWER(?)", params.Query)
	}

	query = query.Limit(params.Limit).
		Offset(params.Offset)

	if params.Order != nil && params.SortVector != nil {
		query = query.Order(params.Order.(string) + " " + params.SortVector.(string))
	}

	if err := query.Find(&foods).Error; err != nil {
		return nil, err
	}

	for i := 0; i < len(foods); i++ {
		if foods[i].PhotoID != 0 {
			var photo model.Photo
			if err := r.DB.Table("photos").Where("id = ?", foods[i].PhotoID).First(&photo).Error; err != nil {
				return nil, err
			}

			foods[i].Photo = photo
		}

	}

	return &model.ListResponse{
		Items:        foods,
		ItemsPerPage: params.Limit,
		PageIndex:    params.PageIndex,
		TotalItems:   int(totalItems),
	}, nil
}

func (r *FoodRepository) GetRestaurantFood(ctx context.Context, restaurantID uint, foodID uint) (*model.Food, error) {
	var food model.Food
	if err := r.DB.WithContext(ctx).Where("restaurant_id = ? AND id = ?", restaurantID, foodID).First(&food).Error; err != nil {
		return nil, err
	}

	var photo model.Photo
	if err := r.DB.Table("photos").Where("id = ?", food.PhotoID).First(&photo).Error; err != nil {
		return nil, err
	}

	food.Photo = photo

	return &food, nil
}

func (r *FoodRepository) CreateRestaurantFood(ctx context.Context, food *model.Food) (*model.Food, error) {
	if err := r.DB.WithContext(ctx).Create(food).Error; err != nil {
		return nil, err
	}
	return food, nil
}

func (r *FoodRepository) UpdateRestaurantFood(ctx context.Context, food *model.Food) (*model.Food, error) {
	var existingFood model.Food
	if err := r.DB.WithContext(ctx).First(&existingFood, food.ID).Error; err != nil {
		return nil, err
	}

	if err := r.DB.WithContext(ctx).Model(&existingFood).Updates(food).Error; err != nil {
		return nil, err
	}

	if existingFood.PhotoID != food.PhotoID {
		if existingFood.PhotoID != 0 {
			if err := r.DeleteFoodPhoto(ctx, existingFood.PhotoID); err != nil {
				return nil, err
			}
		}

		if food.Photo.ID == 0 {
			if err := r.DB.WithContext(ctx).Create(&food.Photo).Error; err != nil {
				return nil, err
			}
		} else {
			if err := r.DB.WithContext(ctx).Model(&food.Photo).Updates(food.Photo).Error; err != nil {
				return nil, err
			}
		}

		existingFood.PhotoID = food.Photo.ID
	}

	return food, nil
}

func (r *FoodRepository) DeleteRestaurantFood(ctx context.Context, foodID uint) error {
	var food model.Food
	if err := r.DB.WithContext(ctx).First(&food, foodID).Error; err != nil {
		return err
	}

	if food.PhotoID != 0 {
		if err := r.DeleteFoodPhoto(ctx, food.PhotoID); err != nil {
			return err
		}
	}

	if err := r.DB.WithContext(ctx).Delete(&model.Food{}, foodID).Error; err != nil {
		return err
	}

	return nil
}

func (r *FoodRepository) DeleteFoodPhoto(ctx context.Context, photoID uint) error {
	if err := r.DB.WithContext(ctx).Delete(&model.Photo{}, photoID).Error; err != nil {
		return err
	}
	return nil
}
