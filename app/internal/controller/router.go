package controller

import (
	"github.com/labstack/echo/v4"
)

func (s *Server) SetupRoutes() {
	v1 := s.App.Group("/api")
	s.setupAuthRoutes(v1)
	s.setupAdminRoutes(v1)
	s.setupOrderRoutes(v1)
	s.setupRestaurantRoutes(v1)
	s.setupProfileRoutes(v1)
}

func (s *Server) setupAuthRoutes(g *echo.Group) {
	auth := g.Group("/auth")
	auth.POST("/login", s.handler.User.SignIn)
	auth.POST("/register", s.handler.User.SignUp)
	auth.POST("/refresh-token", s.handler.User.RefreshToken)
}

func (s *Server) setupProfileRoutes(g *echo.Group) {
	profile := g.Group("/profile", s.jwt.ValidateAuth)
	profile.GET("", s.handler.User.Profile)
	profile.PUT("", s.handler.User.UpdateProfile)
	profile.DELETE("", s.handler.User.DeleteProfile)
	profile.PUT("/change-password", s.handler.User.ChangePassword, s.jwt.ValidateAuth)
}

func (s *Server) setupAdminRoutes(g *echo.Group) {
	admin := g.Group("/admin")
	admin.Use(s.jwt.ValidateAuth)
	admin.Use(s.jwt.ValidateAdmin)
	admin.GET("/owners", s.handler.Admin.GetOwners)
	admin.DELETE("/owners/:id", s.handler.Admin.DeleteOwner)
	admin.POST("/owners", s.handler.Admin.CreateOwner)
	s.setupAdminRestaurantRoutes(admin)
	admin.GET("/clients", s.handler.Admin.GetClients)
	admin.DELETE("/clients/:id", s.handler.Admin.DeleteClient)
	admin.POST("/services", s.handler.Admin.CreateService)
	admin.PUT("/services/:id", s.handler.Admin.UpdateService)
	admin.DELETE("/services/:id", s.handler.Admin.DeleteService)
	admin.GET("/services", s.handler.Restaurant.GetServices)
}

func (s *Server) setupAdminRestaurantRoutes(g *echo.Group) {
	restaurants := g.Group("/restaurants")
	restaurants.GET("", s.handler.Admin.GetRestaurants)
	restaurants.POST("", s.handler.Admin.CreateRestaurant)
	restaurants.POST("/services", s.handler.Admin.CreateService)
	restaurants.DELETE("/services/:id", s.handler.Admin.DeleteService)
	restaurants.PUT("/services/:id", s.handler.Admin.UpdateService)
	restaurants.GET("/services", s.handler.Restaurant.GetServices)
	restaurants.PUT("/:id", s.handler.Admin.UpdateRestaurant)
	restaurants.DELETE("/:id", s.handler.Admin.DeleteRestaurant)
	restaurants.GET("/:id", s.handler.Admin.GetRestaurant)
}

func (s *Server) setupOrderRoutes(g *echo.Group) {
	order := g.Group("/orders")
	order.Use(s.jwt.ValidateAuth)
	order.POST("/create", s.handler.Order.CreateOrder)
	order.DELETE("/:id", s.handler.Order.DeleteOrder)
	order.PUT("/:id", s.handler.Order.UpdateOrder)
	order.GET("/:id", s.handler.Order.GetOrder)
	order.GET("", s.handler.Order.GetAllOrders)
}

func (s *Server) setupRestaurantRoutes(g *echo.Group) {
	restaurant := g.Group("/restaurants")
	restaurant.GET("/statistics", s.handler.Restaurant.GetStatistics)
	restaurant.GET("", s.handler.Restaurant.GetRestaurants, s.jwt.RoleToCtx)
	restaurant.GET("/popular", s.handler.Restaurant.PopularRestaurants)
	restaurant.GET("/services", s.handler.Restaurant.GetServices)
	restaurant.GET("/:id", s.handler.Restaurant.GetRestaurantByID)
	restaurant.GET("/:id/reviews", s.handler.Reviews.GetReviews)
	s.setupTableRoutes(restaurant)
	s.setupMenuRoutes(restaurant)
	restaurant.GET("/:id/orders", s.handler.Restaurant.GetRestaurantOrders, s.jwt.ValidateAuth, s.jwt.ValidateOwner)
	restaurant.POST("/:id/reviews", s.handler.Reviews.CreateReview, s.jwt.ValidateAuth, s.jwt.ValidateUser)
	restaurant.DELETE("/:id/reviews/:review_id", s.handler.Reviews.DeleteReview, s.jwt.ValidateAuth, s.jwt.ValidateUser)

}

func (s *Server) setupMenuRoutes(g *echo.Group) {
	menu := g.Group("/:id/menu")
	menu.GET("/categories", s.handler.Menu.GetMenuCategories)
	menu.GET("", s.handler.Menu.GetRestaurantMenu)
	menu.GET("/:food_id", s.handler.Menu.GetRestaurantFood)
	menu.POST("", s.handler.Menu.CreateRestaurantFood, s.jwt.ValidateAuth, s.jwt.ValidateOwner)
	menu.PUT("/:food_id", s.handler.Menu.UpdateRestaurantFood, s.jwt.ValidateAuth, s.jwt.ValidateOwner)
	menu.DELETE("/:food_id", s.handler.Menu.DeleteRestaurantFood, s.jwt.ValidateAuth, s.jwt.ValidateOwner)
}

func (s *Server) setupTableRoutes(g *echo.Group) {
	tables := g.Group("/:id/tables")
	tables.GET("/categories", s.handler.Table.GetTableCategories)
	tables.GET("", s.handler.Table.GetRestaurantTables)
	tables.GET("/:table_id", s.handler.Table.GetRestaurantTable)
	tables.POST("", s.handler.Table.CreateRestaurantTable, s.jwt.ValidateAuth, s.jwt.ValidateOwner)
	tables.PUT("/:table_id", s.handler.Table.UpdateRestaurantTable, s.jwt.ValidateAuth, s.jwt.ValidateOwner)
	tables.DELETE("/:table_id", s.handler.Table.DeleteRestaurantTable, s.jwt.ValidateAuth, s.jwt.ValidateOwner)
}
