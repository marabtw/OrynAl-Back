package services

import (
	"context"
	"fmt"
	"github.com/alibekabdrakhman1/orynal/config"
	"github.com/alibekabdrakhman1/orynal/internal/model"
	"github.com/alibekabdrakhman1/orynal/internal/repository"
	"github.com/alibekabdrakhman1/orynal/internal/service/infrastructure"
	"github.com/alibekabdrakhman1/orynal/pkg/utils"
	"go.uber.org/zap"
)

func NewMenuService(repository *repository.Manager, config *config.Config, logger *zap.SugaredLogger) *MenuService {
	return &MenuService{repository: repository, config: config, logger: logger, FormatParams: infrastructure.NewFormatParams()}
}

type MenuService struct {
	repository *repository.Manager
	config     *config.Config
	logger     *zap.SugaredLogger
	FormatParams
}

func (s *MenuService) GetMenuCategories(ctx context.Context, restaurantID uint) ([]string, error) {
	return s.repository.Food.GetMenuCategories(ctx, restaurantID)
}

func (s *MenuService) GetRestaurantMenu(ctx context.Context, restaurantID uint, params *model.Params) (*model.ListResponse, error) {
	foods, err := s.repository.Food.GetRestaurantMenu(ctx, restaurantID, params)
	if err != nil {
		return nil, err
	}

	return foods, nil
}

func (s *MenuService) GetRestaurantFood(ctx context.Context, restaurantID, foodID uint) (*model.Food, error) {
	food, err := s.repository.Food.GetRestaurantFood(ctx, restaurantID, foodID)
	if err != nil {
		return nil, err
	}

	return food, nil
}

func (s *MenuService) CreateRestaurantFood(ctx context.Context, restaurantID uint, food *model.Food) (*model.Food, error) {
	if err := s.checkOwner(ctx, restaurantID); err != nil {
		return nil, err
	}
	food.RestaurantID = restaurantID
	createdFood, err := s.repository.Food.CreateRestaurantFood(ctx, food)
	if err != nil {
		return nil, err
	}

	return createdFood, nil
}

func (s *MenuService) UpdateRestaurantFood(ctx context.Context, restaurantID uint, food *model.Food) (*model.Food, error) {
	if err := s.checkOwner(ctx, restaurantID); err != nil {
		return nil, err
	}

	food.RestaurantID = restaurantID

	updatedFood, err := s.repository.Food.UpdateRestaurantFood(ctx, food)
	if err != nil {
		return nil, err
	}

	return updatedFood, nil
}

func (s *MenuService) DeleteRestaurantFood(ctx context.Context, restaurantID, foodID uint) error {
	if err := s.checkOwner(ctx, restaurantID); err != nil {
		return err
	}

	err := s.repository.Food.DeleteRestaurantFood(ctx, foodID)
	if err != nil {
		return err
	}

	return nil
}

func (s *MenuService) checkOwner(ctx context.Context, restaurantID uint) error {
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
