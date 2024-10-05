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
	"time"
)

func NewTableService(repository *repository.Manager, config *config.Config, logger *zap.SugaredLogger) *TableService {
	return &TableService{repository: repository, config: config, logger: logger, FormatParams: infrastructure.NewFormatParams()}
}

type TableService struct {
	repository *repository.Manager
	config     *config.Config
	logger     *zap.SugaredLogger
	FormatParams
}

func (s *TableService) GetTableCategories(ctx context.Context, restaurantID uint) ([]string, error) {
	return s.repository.Table.GetTableCategories(ctx, restaurantID)
}

func (s *TableService) GetRestaurantTables(ctx context.Context, restaurantID uint, params *model.Params) (*model.ListResponse, error) {
	return s.repository.Table.GetRestaurantTables(ctx, restaurantID, params)
}

func (s *TableService) GetRestaurantTable(ctx context.Context, restaurantID uint, tableID uint) (*model.Table, error) {
	return s.repository.Table.GetRestaurantTable(ctx, restaurantID, tableID)
}

func (s *TableService) CreateRestaurantTable(ctx context.Context, restaurantID uint, table *model.Table) (*model.Table, error) {
	err := s.checkOwner(ctx, restaurantID)
	if err != nil {
		return nil, err
	}

	table.RestaurantID = restaurantID

	return s.repository.Table.CreateTable(ctx, table)
}

func (s *TableService) UpdateRestaurantTable(ctx context.Context, restaurantID uint, table *model.Table) (*model.Table, error) {
	err := s.checkOwner(ctx, restaurantID)
	if err != nil {
		return nil, err
	}

	table.RestaurantID = restaurantID

	return s.repository.Table.UpdateTable(ctx, table)
}

func (s *TableService) DeleteRestaurantTable(ctx context.Context, restaurantID uint, tableID uint) error {
	err := s.checkOwner(ctx, restaurantID)
	if err != nil {
		return err
	}

	return s.repository.Table.DeleteTable(ctx, tableID)
}

func (s *TableService) GetAvailableTime(ctx context.Context, restaurantID uint, tableID uint, date time.Time) ([]time.Time, error) {
	//TODO implement me
	panic("implement me")
}

func (s *TableService) checkOwner(ctx context.Context, restaurantID uint) error {
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
