package postgre

import (
	"context"
	"errors"
	"github.com/alibekabdrakhman1/orynal/internal/model"
	"gorm.io/gorm"
)

func NewRestaurantRepository(db *gorm.DB) *RestaurantRepository {
	return &RestaurantRepository{
		DB: db,
	}
}

type RestaurantRepository struct {
	DB *gorm.DB
}

func (r *RestaurantRepository) GetPopularRestaurants(ctx context.Context) (*model.ListResponse, error) {
	var restaurants []model.Restaurant

	err := r.DB.Table("orders").
		WithContext(ctx).
		Select("restaurants.id, restaurants.name, restaurants.address, restaurants.description, restaurants.city, restaurants.status, restaurants.phone, restaurants.owner_id, restaurants.mode_from, restaurants.mode_to, restaurants.icon_id, count(orders.id) as order_count").
		Joins("JOIN restaurants ON restaurants.id = orders.restaurant_id").
		Group("restaurants.id").
		Order("order_count DESC").
		Limit(10).
		Find(&restaurants).Error

	for i := 0; i < len(restaurants); i++ {
		var owner model.UserResponse
		if err := r.DB.Table("users").Where("id = ?", restaurants[i].OwnerID).First(&owner).Error; err != nil {
			return nil, err
		}

		var icon model.Photo
		if err := r.DB.Table("photos").Where("id = ?", restaurants[i].IconID).First(&icon).Error; err != nil {
			return nil, err
		}

		var services []model.Service
		if err := r.DB.Raw("SELECT services.* FROM services JOIN restaurant_service ON services.id = restaurant_service.service_id WHERE restaurant_service.restaurant_id = ?", restaurants[i].ID).Scan(&services).Error; err != nil {
			return nil, err
		}

		var photos []model.Photo
		if err := r.DB.Raw("SELECT photos.* FROM photos JOIN restaurant_photos ON photos.id = restaurant_photos.photo_id WHERE restaurant_photos.restaurant_id = ?", restaurants[i].ID).Scan(&photos).Error; err != nil {
			return nil, err
		}

		restaurants[i].Owner = owner
		restaurants[i].Services = services
		restaurants[i].Icon = icon
		restaurants[i].Photos = photos
	}

	if err != nil {
		return nil, err
	}

	return &model.ListResponse{
		Items:        restaurants,
		ItemsPerPage: 10,
		TotalPages:   1,
		PageIndex:    1,
		TotalItems:   10,
	}, nil
}

func (r *RestaurantRepository) GetRestaurants(ctx context.Context, params *model.Params) (*model.ListResponse, error) {
	var restaurants []model.Restaurant
	var totalItems int64

	countQuery := r.DB.WithContext(ctx)
	if params.Query != "" {
		countQuery = countQuery.Where("LOWER(name) LIKE LOWER(?)", "%"+params.Query+"%")
	}
	if err := countQuery.Table("restaurants").Count(&totalItems).Error; err != nil {
		return nil, err
	}

	if int(totalItems) <= params.Offset {
		return nil, errors.New("offset cannot be less than total items")
	}

	query := r.DB.WithContext(ctx).Table("restaurants").
		Limit(params.Limit).
		Offset(params.Offset)

	if params.Query != "" {
		query = query.Where("LOWER(name) LIKE LOWER(?)", "%"+params.Query+"%")
	}

	query = query.Find(&restaurants)

	if err := query.Find(&restaurants).Error; err != nil {
		return nil, err
	}

	for i := 0; i < len(restaurants); i++ {
		var owner model.UserResponse
		if err := r.DB.Table("users").Where("id = ?", restaurants[i].OwnerID).First(&owner).Error; err != nil {
			return nil, err
		}

		var icon model.Photo
		if err := r.DB.Table("photos").Where("id = ?", restaurants[i].IconID).First(&icon).Error; err != nil {
			return nil, err
		}

		var services []model.Service
		if err := r.DB.Raw("SELECT services.* FROM services JOIN restaurant_service ON services.id = restaurant_service.service_id WHERE restaurant_service.restaurant_id = ?", restaurants[i].ID).Scan(&services).Error; err != nil {
			return nil, err
		}

		restaurants[i].Owner = owner
		restaurants[i].Services = services
		restaurants[i].Icon = icon
	}

	return &model.ListResponse{
		Items:        restaurants,
		ItemsPerPage: params.Limit,
		PageIndex:    params.PageIndex,
		TotalItems:   int(totalItems),
	}, nil
}

func (r *RestaurantRepository) GetStatistics(ctx context.Context) (*model.Statistics, error) {
	var countRestaurants int64
	if err := r.DB.WithContext(ctx).Table("restaurants").Model(&model.Restaurant{}).Count(&countRestaurants).Error; err != nil {
		return nil, err
	}
	var countOrders int64
	if err := r.DB.WithContext(ctx).Table("orders").Model(&model.Order{}).Count(&countOrders).Error; err != nil {
		return nil, err
	}

	return &model.Statistics{
		OrderCount:       countOrders,
		PeopleCount:      countOrders * 5,
		RestaurantsCount: countRestaurants,
	}, nil
}

func (r *RestaurantRepository) GetRestaurantByID(ctx context.Context, id uint) (*model.Restaurant, error) {
	var restaurantResponse model.Restaurant

	if err := r.DB.WithContext(ctx).
		Table("restaurants").
		Preload("Photo").
		First(&restaurantResponse, id).Error; err != nil {
		return &model.Restaurant{}, err
	}

	var owner model.UserResponse
	if err := r.DB.Table("users").Where("id = ?", restaurantResponse.OwnerID).First(&owner).Error; err != nil {
		return nil, err
	}

	var icon model.Photo
	if err := r.DB.Table("photos").Where("id = ?", restaurantResponse.IconID).First(&icon).Error; err != nil {
		return nil, err
	}

	var services []model.Service
	if err := r.DB.Raw("SELECT services.* FROM services JOIN restaurant_service ON services.id = restaurant_service.service_id WHERE restaurant_service.restaurant_id = ?", id).Scan(&services).Error; err != nil {
		return nil, err
	}

	var photos []model.Photo
	if err := r.DB.Raw("SELECT photos.* FROM photos JOIN restaurant_photos ON photos.id = restaurant_photos.photo_id WHERE restaurant_photos.restaurant_id = ?", id).Scan(&photos).Error; err != nil {
		return nil, err
	}

	restaurantResponse.Owner = owner
	restaurantResponse.Services = services
	restaurantResponse.Photos = photos
	restaurantResponse.Icon = icon

	return &restaurantResponse, nil
}

func (r *RestaurantRepository) GetRestaurantsByOwner(ctx context.Context, ownerID uint, params *model.Params) (*model.ListResponse, error) {
	var restaurants []model.Restaurant
	var totalItems int64

	countQuery := r.DB.WithContext(ctx)
	if params.Query != "" {
		countQuery = countQuery.Where("LOWER(name) LIKE LOWER(?)", "%"+params.Query+"%")
	}
	if err := countQuery.Table("restaurants").Where("owner_id = ?", ownerID).Count(&totalItems).Error; err != nil {
		return nil, err
	}

	if int(totalItems) <= params.Offset {
		return nil, errors.New("offset cannot be less than total items")
	}

	query := r.DB.WithContext(ctx).Table("restaurants").
		Where("owner_id = ?", ownerID).
		Limit(params.Limit).
		Offset(params.Offset)

	if params.Query != "" {
		query = query.Where("LOWER(name) LIKE LOWER(?)", "%"+params.Query+"%")
	}

	query = query.Find(&restaurants)

	if err := query.Find(&restaurants).Error; err != nil {
		return nil, err
	}

	for i := 0; i < len(restaurants); i++ {
		var owner model.UserResponse
		if err := r.DB.Table("users").Where("id = ?", restaurants[i].OwnerID).First(&owner).Error; err != nil {
			return nil, err
		}

		var icon model.Photo
		if err := r.DB.Table("photos").Where("id = ?", restaurants[i].IconID).First(&icon).Error; err != nil {
			return nil, err
		}

		var services []model.Service
		if err := r.DB.Raw("SELECT services.* FROM services JOIN restaurant_service ON services.id = restaurant_service.service_id WHERE restaurant_service.restaurant_id = ?", restaurants[i].ID).Scan(&services).Error; err != nil {
			return nil, err
		}

		restaurants[i].Owner = owner
		restaurants[i].Services = services
		restaurants[i].Icon = icon
	}

	return &model.ListResponse{
		Items:        restaurants,
		ItemsPerPage: params.Limit,
		PageIndex:    params.PageIndex,
		TotalItems:   int(totalItems),
	}, nil
}

func (r *RestaurantRepository) GetFavoriteRestaurants(ctx context.Context, userID uint, params *model.Params) (*model.ListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RestaurantRepository) CreateRestaurant(ctx context.Context, restaurant *model.Restaurant) (*model.Restaurant, error) {
	tx := r.DB.WithContext(ctx).Begin()

	if restaurant.Icon.Route != "" {
		iconPhoto := model.Photo{Route: restaurant.Icon.Route}
		if err := tx.Table("photos").Create(&iconPhoto).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		restaurant.IconID = iconPhoto.ID
	}

	createRestaurant := model.Restaurant{
		Name:        restaurant.Name,
		Address:     restaurant.Address,
		Description: restaurant.Description,
		City:        restaurant.City,
		Status:      restaurant.Status,
		OwnerID:     restaurant.OwnerID,
		Phone:       restaurant.Phone,
		ModeFrom:    restaurant.ModeFrom,
		ModeTo:      restaurant.ModeTo,
		IconID:      restaurant.IconID,
	}

	if err := tx.Table("restaurants").Create(&createRestaurant).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	for _, service := range restaurant.Services {
		restaurantService := model.RestaurantService{
			ServiceID:    service.ID,
			RestaurantID: createRestaurant.ID,
		}
		if err := tx.Table("restaurant_service").Create(&restaurantService).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if len(restaurant.Photos) > 0 {
		var photos []model.Photo
		for _, photo := range restaurant.Photos {
			photos = append(photos, model.Photo{Route: photo.Route})
		}
		if err := tx.Table("photos").Create(&photos).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		var restaurantPhotos []model.RestaurantPhoto
		for _, photo := range photos {
			restaurantPhoto := model.RestaurantPhoto{
				PhotoID:      photo.ID,
				RestaurantID: createRestaurant.ID,
			}
			restaurantPhotos = append(restaurantPhotos, restaurantPhoto)
		}
		if err := tx.Table("restaurant_photos").Create(&restaurantPhotos).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return r.GetRestaurantByID(ctx, createRestaurant.ID)
}

func (r *RestaurantRepository) DeleteRestaurant(ctx context.Context, restaurantID uint) error {
	if err := r.DB.WithContext(ctx).Table("restaurants").Delete(&model.Restaurant{}, restaurantID).Error; err != nil {
		return err
	}
	return nil
}

func (r *RestaurantRepository) UpdateRestaurant(ctx context.Context, restaurantID uint, restaurant *model.Restaurant) (*model.Restaurant, error) {
	var existingRestaurant model.Restaurant
	if err := r.DB.WithContext(ctx).Table("restaurants").First(&existingRestaurant, restaurantID).Error; err != nil {
		return nil, err
	}

	if err := r.DB.WithContext(ctx).Table("restaurants").Model(&existingRestaurant).Updates(restaurant).Error; err != nil {
		return nil, err
	}

	if err := r.UpdateRestaurantPhotos(ctx, restaurantID, restaurant.Photos); err != nil {
		return nil, err
	}

	if err := r.UpdateRestaurantServices(ctx, restaurantID, restaurant.Services); err != nil {
		return nil, err
	}

	return r.GetRestaurantByID(ctx, restaurantID)
}

func (r *RestaurantRepository) UpdateRestaurantPhotos(ctx context.Context, restaurantID uint, photos []model.Photo) error {
	var existingPhotos []model.RestaurantPhoto
	if err := r.DB.WithContext(ctx).Table("restaurant_photos").Where("restaurant_id = ?", restaurantID).Find(&existingPhotos).Error; err != nil {
		return err
	}

	if err := r.DB.WithContext(ctx).Table("restaurant_photos").Where("restaurant_id = ?", restaurantID).Delete(&model.RestaurantPhoto{}).Error; err != nil {
		return err
	}

	for _, photo := range photos {
		newPhoto := model.Photo{Route: photo.Route}
		if err := r.DB.WithContext(ctx).Table("photos").Create(&newPhoto).Error; err != nil {
			return err
		}

		newRestaurantPhoto := model.RestaurantPhoto{
			PhotoID:      newPhoto.ID,
			RestaurantID: restaurantID,
		}
		if err := r.DB.WithContext(ctx).Table("restaurant_photos").Create(&newRestaurantPhoto).Error; err != nil {
			return err
		}
	}

	for _, oldPhoto := range existingPhotos {
		var count int64
		if err := r.DB.WithContext(ctx).Table("restaurant_photos").Where("photo_id = ?", oldPhoto.PhotoID).Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			if err := r.DB.WithContext(ctx).Table("photos").Delete(&model.Photo{ID: oldPhoto.PhotoID}).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *RestaurantRepository) UpdateRestaurantServices(ctx context.Context, restaurantID uint, services []model.Service) error {
	if err := r.DB.WithContext(ctx).Table("restaurant_service").Where("restaurant_id = ?", restaurantID).Delete(&model.RestaurantService{}).Error; err != nil {
		return err
	}

	for _, service := range services {
		if err := r.DB.WithContext(ctx).Table("restaurant_service").Create(&model.RestaurantService{
			RestaurantID: restaurantID,
			ServiceID:    service.ID,
		}).Error; err != nil {
			return err
		}
	}

	return nil
}
