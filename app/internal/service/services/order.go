package services

import (
	"context"
	"errors"
	"github.com/alibekabdrakhman1/orynal/config"
	"github.com/alibekabdrakhman1/orynal/internal/model"
	"github.com/alibekabdrakhman1/orynal/internal/repository"
	"github.com/alibekabdrakhman1/orynal/internal/service/infrastructure"
	"github.com/alibekabdrakhman1/orynal/pkg/enums"
	"github.com/alibekabdrakhman1/orynal/pkg/utils"
	"go.uber.org/zap"
)

func NewOrderService(repository *repository.Manager, config *config.Config, logger *zap.SugaredLogger) *OrderService {
	return &OrderService{repository: repository, config: config, logger: logger, FormatParams: infrastructure.NewFormatParams()}
}

type OrderService struct {
	repository *repository.Manager
	config     *config.Config
	logger     *zap.SugaredLogger
	FormatParams
}

func (s *OrderService) Create(ctx context.Context, order *model.OrderRequest) (*model.OrderResponse, error) {
	userID, err := utils.GetIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	createdOrder, err := s.repository.Order.CreateOrder(ctx, &model.Order{
		RestaurantID: order.RestaurantID,
		TotalSum:     order.TotalSum,
		UserID:       userID,
		TableID:      order.TableID,
		Date:         order.Date,
		Status:       order.Status,
		OrderFoods:   order.OrderFoods,
	})
	if err != nil {
		return nil, err
	}
	return createdOrder, nil
}

func (s *OrderService) Update(ctx context.Context, id uint, order *model.Order) (*model.OrderResponse, error) {
	userID, err := utils.GetIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	role, err := utils.GetRoleFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if role == enums.Owner {
		return s.repository.Order.UpdateOrder(ctx, order)
	}

	oldOrder, err := s.repository.Order.GetOrder(ctx, id)
	if err != nil {
		return nil, err
	}

	if role == enums.User && oldOrder.UserID != userID {
		return nil, errors.New("permission denied")
	}

	switch oldOrder.Status {
	case enums.Canceled:
		return nil, errors.New("order status is canceled")
	case enums.Completed:
		return nil, errors.New("order status is completed")
	}

	updatedOrder, err := s.repository.Order.UpdateOrder(ctx, order)
	if err != nil {
		return nil, err
	}
	return updatedOrder, nil
}

func (s *OrderService) Delete(ctx context.Context, id uint) error {
	userID, err := utils.GetIDFromContext(ctx)
	if err != nil {
		return err
	}

	role, err := utils.GetRoleFromContext(ctx)
	if err != nil {
		return err
	}

	order, err := s.repository.Order.GetOrder(ctx, id)
	if err != nil {
		return err
	}

	if role == enums.User && order.UserID != userID {
		return errors.New("permission denied")
	}

	err = s.repository.Order.DeleteOrder(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *OrderService) GetByID(ctx context.Context, id uint) (*model.OrderResponse, error) {
	userID, err := utils.GetIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	role, err := utils.GetRoleFromContext(ctx)
	if err != nil {
		return nil, err
	}

	order, err := s.repository.Order.GetOrder(ctx, id)
	if err != nil {
		return nil, err
	}

	if role == enums.User && order.UserID != userID {
		return nil, errors.New("permission denied")
	}

	return order, nil
}

func (s *OrderService) GetAllOrders(ctx context.Context, params *model.Params) (*model.ListResponse, error) {
	id, err := utils.GetIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	orders, err := s.repository.Order.GetAllOrders(ctx, id, params)
	if err != nil {
		return nil, err
	}
	return orders, nil
}
