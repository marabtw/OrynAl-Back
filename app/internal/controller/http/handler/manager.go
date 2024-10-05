package http

import (
	"github.com/alibekabdrakhman1/orynal/internal/controller/http/handler/handlers"
	"github.com/alibekabdrakhman1/orynal/internal/service"
	"go.uber.org/zap"
)

type Manager struct {
	User       IUserHandler
	Admin      IAdminHandler
	Order      IOrderHandler
	Restaurant IRestaurantHandler
	Table      ITableHandler
	Menu       IMenuHandler
	Reviews    IReviewsHandler
}

func NewManager(srv *service.Manager, logger *zap.SugaredLogger) *Manager {
	return &Manager{
		User:       handlers.NewUserHandler(srv, logger),
		Admin:      handlers.NewAdminHandler(srv, logger),
		Order:      handlers.NewOrderHandler(srv, logger),
		Restaurant: handlers.NewRestaurantHandler(srv, logger),
		Table:      handlers.NewTableHandler(srv, logger),
		Menu:       handlers.NewMenuHandler(srv, logger),
		Reviews:    handlers.NewReviewsHandler(srv, logger),
	}
}
