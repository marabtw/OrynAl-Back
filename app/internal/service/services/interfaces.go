package services

import (
	"context"
	"github.com/alibekabdrakhman1/orynal/internal/model"
	"github.com/labstack/echo/v4"
	"time"
)

type IAuthService interface {
	Login(ctx context.Context, login model.Login) (*model.JwtTokens, error)
	Register(ctx context.Context, user model.Register) (uint, error)
	RefreshToken(ctx context.Context, refreshToken string) (*model.JwtTokens, error)
	GetJwtUserID(jwtToken string) (*model.ContextUserID, error)
	GetJwtUserRole(jwtToken string) (*model.ContextUserRole, error)
}

type IUserService interface {
	Create(ctx context.Context, user *model.User) (*model.UserResponse, error)
	CreateOwner(ctx context.Context, user *model.User) (*model.UserResponse, error)
	Update(ctx context.Context, user *model.User) (*model.UserResponse, error)
	ChangePassword(ctx context.Context, user *model.ChangePasswordRequest) error
	Delete(ctx context.Context, id uint) error
	Profile(ctx context.Context) (*model.UserResponse, error)
	GetByID(ctx context.Context, id uint) (*model.UserResponse, error)
	GetAllClients(ctx context.Context, params *model.Params) (*model.ListResponse, error)
	GetAllOwners(ctx context.Context, params *model.Params) (*model.ListResponse, error)
	FormatParams
}

type IRestaurantService interface {
	GetRestaurants(ctx context.Context, params *model.Params) (*model.ListResponse, error)
	GetStatistics(ctx context.Context) (*model.Statistics, error)
	GetRestaurantByID(ctx context.Context, id uint) (*model.Restaurant, error)
	CreateRestaurant(ctx context.Context, restaurant *model.Restaurant) (*model.Restaurant, error)
	UpdateRestaurant(ctx context.Context, restaurant *model.Restaurant, id uint) (*model.Restaurant, error)
	DeleteRestaurant(ctx context.Context, id uint) error
	FavoriteRestaurants(ctx context.Context, id uint, params *model.Params) (*model.ListResponse, error)
	PopularRestaurants(ctx context.Context) (*model.ListResponse, error)
	GetRestaurantOrders(ctx context.Context, id uint, params *model.Params) (*model.ListResponse, error)
	CreateService(ctx context.Context, service *model.Service) ([]model.Service, error)
	DeleteService(ctx context.Context, id uint) error
	GetServices(ctx context.Context) ([]model.Service, error)
	UpdateService(ctx context.Context, service *model.Service) ([]model.Service, error)
	FormatParams
}

type FormatParams interface {
	RestaurantsSearchFormatting(model *model.Params, ctx echo.Context) (*model.Params, error)
	TablesSearchFormatting(params *model.Params, ctx echo.Context) (*model.Params, error)
	UserSearchFormatting(params *model.Params, ctx echo.Context) (*model.Params, error)
	OrderSearchFormatting(params *model.Params, ctx echo.Context) (*model.Params, error)
	MenuSearchFormatting(params *model.Params, ctx echo.Context) (*model.Params, error)
	ReviewsSearchFormatting(params *model.Params, ctx echo.Context) (*model.Params, error)
}

type IOrderService interface {
	Create(ctx context.Context, order *model.OrderRequest) (*model.OrderResponse, error)
	Update(ctx context.Context, id uint, order *model.Order) (*model.OrderResponse, error)
	Delete(ctx context.Context, id uint) error
	GetByID(ctx context.Context, id uint) (*model.OrderResponse, error)
	GetAllOrders(ctx context.Context, params *model.Params) (*model.ListResponse, error)
	FormatParams
}

type IMenuService interface {
	GetRestaurantMenu(ctx context.Context, restaurantID uint, params *model.Params) (*model.ListResponse, error)
	GetRestaurantFood(ctx context.Context, restaurantID, foodID uint) (*model.Food, error)
	CreateRestaurantFood(ctx context.Context, restaurantID uint, food *model.Food) (*model.Food, error)
	UpdateRestaurantFood(ctx context.Context, restaurantID uint, food *model.Food) (*model.Food, error)
	DeleteRestaurantFood(ctx context.Context, restaurantID, foodID uint) error
	GetMenuCategories(ctx context.Context, restaurantID uint) ([]string, error)
	FormatParams
}

type ITableService interface {
	GetRestaurantTables(ctx context.Context, restaurantID uint, params *model.Params) (*model.ListResponse, error)
	GetRestaurantTable(ctx context.Context, restaurantID uint, tableID uint) (*model.Table, error)
	CreateRestaurantTable(ctx context.Context, restaurantID uint, table *model.Table) (*model.Table, error)
	UpdateRestaurantTable(ctx context.Context, restaurantID uint, table *model.Table) (*model.Table, error)
	DeleteRestaurantTable(ctx context.Context, restaurantID uint, tableID uint) error
	GetAvailableTime(ctx context.Context, restaurantID uint, tableID uint, date time.Time) ([]time.Time, error)
	GetTableCategories(ctx context.Context, restaurantID uint) ([]string, error)
	FormatParams
}

type IReviewsService interface {
	GetReviews(ctx context.Context, restaurantID uint, params *model.Params) (*model.ListResponse, error)
	CreateReview(ctx context.Context, review *model.RestaurantReview) (*model.RestaurantReview, error)
	DeleteReview(ctx context.Context, id uint) error
	FormatParams
}
