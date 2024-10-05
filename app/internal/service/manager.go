package service

import (
	"github.com/alibekabdrakhman1/orynal/config"
	"github.com/alibekabdrakhman1/orynal/internal/repository"
	"github.com/alibekabdrakhman1/orynal/internal/service/services"
	"go.uber.org/zap"
)

type Manager struct {
	Auth       services.IAuthService
	User       services.IUserService
	Restaurant services.IRestaurantService
	Table      services.ITableService
	Menu       services.IMenuService
	Order      services.IOrderService
	Reviews    services.IReviewsService
}

func NewManager(repository *repository.Manager, config *config.Config, logger *zap.SugaredLogger) *Manager {
	return &Manager{
		Auth:       services.NewAuthService(repository, config, logger),
		User:       services.NewUserService(repository, config, logger),
		Restaurant: services.NewRestaurantService(repository, config, logger),
		Table:      services.NewTableService(repository, config, logger),
		Menu:       services.NewMenuService(repository, config, logger),
		Order:      services.NewOrderService(repository, config, logger),
		Reviews:    services.NewReviewsService(repository, config, logger),
	}
}
