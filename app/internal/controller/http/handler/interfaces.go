package http

import "github.com/labstack/echo/v4"

type IUserHandler interface {
	SignIn(c echo.Context) error
	SignUp(c echo.Context) error
	RefreshToken(c echo.Context) error
	Profile(c echo.Context) error
	UpdateProfile(c echo.Context) error
	ChangePassword(c echo.Context) error
	DeleteProfile(c echo.Context) error
}

type IAdminHandler interface {
	GetClients(c echo.Context) error
	GetClient(c echo.Context) error
	DeleteClient(c echo.Context) error
	GetOwners(c echo.Context) error
	CreateOwner(c echo.Context) error
	DeleteOwner(c echo.Context) error
	GetRestaurants(c echo.Context) error
	GetRestaurant(c echo.Context) error
	CreateRestaurant(c echo.Context) error
	DeleteRestaurant(c echo.Context) error
	UpdateRestaurant(c echo.Context) error
	CreateService(c echo.Context) error
	DeleteService(c echo.Context) error
	UpdateService(c echo.Context) error
}

type IRestaurantHandler interface {
	GetRestaurants(c echo.Context) error
	GetStatistics(c echo.Context) error
	GetRestaurantByID(c echo.Context) error
	SavedRestaurants(c echo.Context) error
	SaveRestaurant(c echo.Context) error
	UnsaveRestaurant(c echo.Context) error
	PopularRestaurants(c echo.Context) error
	GetRestaurantOrders(c echo.Context) error
	DeleteRestaurant(c echo.Context) error
	UpdateRestaurant(c echo.Context) error
	GetServices(c echo.Context) error
}

type IReviewsHandler interface {
	CreateReview(c echo.Context) error
	GetReviews(c echo.Context) error
	DeleteReview(c echo.Context) error
}

type ITableHandler interface {
	GetTableCategories(c echo.Context) error
	GetRestaurantTables(c echo.Context) error
	GetRestaurantTable(c echo.Context) error
	CreateRestaurantTable(c echo.Context) error
	UpdateRestaurantTable(c echo.Context) error
	DeleteRestaurantTable(c echo.Context) error
	GetAvailableTime(c echo.Context) error
}

type IMenuHandler interface {
	GetMenuCategories(c echo.Context) error
	GetRestaurantMenu(c echo.Context) error
	GetRestaurantFood(c echo.Context) error
	CreateRestaurantFood(c echo.Context) error
	UpdateRestaurantFood(c echo.Context) error
	DeleteRestaurantFood(c echo.Context) error
}

type IOrderHandler interface {
	CreateOrder(c echo.Context) error
	DeleteOrder(c echo.Context) error
	UpdateOrder(c echo.Context) error
	GetOrder(c echo.Context) error
	GetAllOrders(c echo.Context) error
}
