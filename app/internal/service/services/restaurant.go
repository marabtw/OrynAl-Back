package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/alibekabdrakhman1/orynal/config"
	"github.com/alibekabdrakhman1/orynal/internal/model"
	"github.com/alibekabdrakhman1/orynal/internal/repository"
	"github.com/alibekabdrakhman1/orynal/internal/service/infrastructure"
	"github.com/alibekabdrakhman1/orynal/pkg/enums"
	"github.com/alibekabdrakhman1/orynal/pkg/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func NewRestaurantService(repository *repository.Manager, config *config.Config, logger *zap.SugaredLogger) *RestaurantService {
	return &RestaurantService{repository: repository, config: config, logger: logger, FormatParams: infrastructure.NewFormatParams()}
}

type RestaurantService struct {
	repository *repository.Manager
	config     *config.Config
	logger     *zap.SugaredLogger
	FormatParams
}

func (s *RestaurantService) GetStatistics(ctx context.Context) (*model.Statistics, error) {
	return s.repository.Restaurant.GetStatistics(ctx)
}

func (s *RestaurantService) PopularRestaurants(ctx context.Context) (*model.ListResponse, error) {
	return s.repository.Restaurant.GetPopularRestaurants(ctx)
}

func (s *RestaurantService) CreateService(ctx context.Context, service *model.Service) ([]model.Service, error) {
	return s.repository.Services.CreateService(ctx, service)
}

func (s *RestaurantService) DeleteService(ctx context.Context, id uint) error {
	return s.repository.Services.DeleteService(ctx, id)
}

func (s *RestaurantService) GetServices(ctx context.Context) ([]model.Service, error) {
	return s.repository.Services.GetServices(ctx)
}

func (s *RestaurantService) UpdateService(ctx context.Context, service *model.Service) ([]model.Service, error) {
	return s.repository.Services.UpdateService(ctx, service)
}

func (s *RestaurantService) GetRestaurants(ctx context.Context, params *model.Params) (*model.ListResponse, error) {
	role, err := utils.GetRoleFromContext(ctx)
	if err != nil {
		return s.repository.Restaurant.GetRestaurants(ctx, params)
	}

	if role == enums.Owner {
		id, err := utils.GetIDFromContext(ctx)
		if err != nil {
			s.logger.Error(err)
			return nil, err
		}
		return s.repository.Restaurant.GetRestaurantsByOwner(ctx, id, params)
	}

	return s.repository.Restaurant.GetRestaurants(ctx, params)
}

func (s *RestaurantService) GetRestaurantByID(ctx context.Context, id uint) (*model.Restaurant, error) {
	return s.repository.Restaurant.GetRestaurantByID(ctx, id)
}

func (s *RestaurantService) CreateRestaurant(ctx context.Context, restaurant *model.Restaurant) (*model.Restaurant, error) {
	role, err := utils.GetRoleFromContext(ctx)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	if role != enums.Admin {
		return nil, errors.New("permission denied")
	}

	owner, err := s.repository.User.GetByID(ctx, restaurant.OwnerID)
	if err != nil {
		s.logger.Error(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("owner not found")
		}
		return nil, errors.New(fmt.Sprintf("error get owner from db: %s", err.Error()))
	}

	if owner.Role != enums.Owner {
		return nil, errors.New("by this id is not owner")
	}

	return s.repository.Restaurant.CreateRestaurant(ctx, restaurant)
}

func (s *RestaurantService) UpdateRestaurant(ctx context.Context, restaurant *model.Restaurant, id uint) (*model.Restaurant, error) {
	role, err := utils.GetRoleFromContext(ctx)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	switch role {
	case enums.Admin:
		return s.repository.Restaurant.UpdateRestaurant(ctx, id, restaurant)
	case enums.Owner:
		if err := s.checkOwner(ctx, id); err != nil {
			return nil, err
		}
		return s.repository.Restaurant.UpdateRestaurant(ctx, id, restaurant)
	default:
		return nil, errors.New("permission denied")
	}
}

func (s *RestaurantService) DeleteRestaurant(ctx context.Context, id uint) error {
	role, err := utils.GetRoleFromContext(ctx)
	if err != nil {
		s.logger.Error(err)
		return err
	}

	switch role {
	case enums.Admin:
		return s.repository.Restaurant.DeleteRestaurant(ctx, id)
	case enums.Owner:
		if err := s.checkOwner(ctx, id); err != nil {
			return err
		}
		return s.repository.Restaurant.DeleteRestaurant(ctx, id)
	default:
		return errors.New("permission denied")
	}
}

func (s *RestaurantService) FavoriteRestaurants(ctx context.Context, id uint, params *model.Params) (*model.ListResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *RestaurantService) GetRestaurantOrders(ctx context.Context, id uint, params *model.Params) (*model.ListResponse, error) {
	role, err := utils.GetRoleFromContext(ctx)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	switch role {
	case enums.Admin:
		return s.repository.Order.GetRestaurantOrders(ctx, id, params)
	case enums.Owner:
		if err := s.checkOwner(ctx, id); err != nil {
			return nil, err
		}
		return s.repository.Order.GetRestaurantOrders(ctx, id, params)
	default:
		return nil, errors.New("permission denied")
	}
}

func (s *RestaurantService) checkOwner(ctx context.Context, restaurantID uint) error {
	restaurant, err := s.repository.Restaurant.GetRestaurantByID(ctx, restaurantID)
	if err != nil {
		s.logger.Error(fmt.Errorf("there is not restaurant by id: %v\n%w", restaurantID, err))
		return fmt.Errorf("there is not restaurant by id: %v", restaurantID)
	}

	userID, err := utils.GetIDFromContext(ctx)
	if err != nil {
		s.logger.Error(err)
		return fmt.Errorf("there is not user id from ctx")
	}

	if restaurant.Owner.ID != userID {
		s.logger.Error(fmt.Errorf("the user id is not owner of restaurant"))
		return fmt.Errorf("permission denied")
	}

	return nil
}
