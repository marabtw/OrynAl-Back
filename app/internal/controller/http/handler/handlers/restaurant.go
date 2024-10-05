package handlers

import (
	"github.com/alibekabdrakhman1/orynal/internal/model"
	"github.com/alibekabdrakhman1/orynal/internal/service"
	"github.com/alibekabdrakhman1/orynal/internal/service/infrastructure"
	"github.com/alibekabdrakhman1/orynal/pkg/response"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

func NewRestaurantHandler(service *service.Manager, logger *zap.SugaredLogger) *RestaurantHandler {
	return &RestaurantHandler{
		service:      service,
		logger:       logger,
		FormatParams: infrastructure.NewFormatParams(),
	}
}

type RestaurantHandler struct {
	service *service.Manager
	logger  *zap.SugaredLogger
	*infrastructure.FormatParams
}

func (h *RestaurantHandler) SavedRestaurants(c echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (h *RestaurantHandler) SaveRestaurant(c echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (h *RestaurantHandler) UnsaveRestaurant(c echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (h *RestaurantHandler) PopularRestaurants(c echo.Context) error {
	restaurants, err := h.service.Restaurant.PopularRestaurants(c.Request().Context())
	if err != nil {
		h.logger.Error("Failed to get restaurants:", err)
		return c.JSON(http.StatusInternalServerError, response.CustomResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to get restaurants",
			Data:    err.Error(),
		})
	}

	return c.JSON(http.StatusOK, response.CustomResponse{
		Status:  http.StatusOK,
		Message: "Success",
		Data:    restaurants,
	})
}

func (h *RestaurantHandler) GetServices(c echo.Context) error {
	services, err := h.service.Restaurant.GetServices(c.Request().Context())
	if err != nil {
		h.logger.Error("Failed to get services:", err)
		return c.JSON(http.StatusInternalServerError, response.CustomResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to get services",
			Data:    err.Error(),
		})
	}

	return c.JSON(http.StatusOK, response.CustomResponse{
		Status:  http.StatusOK,
		Message: "Success",
		Data:    services,
	})
}

func (h *RestaurantHandler) DeleteRestaurant(c echo.Context) error {
	restaurantID := c.Param("id")
	if restaurantID == "" {
		return c.JSON(http.StatusBadRequest, response.CustomResponse{
			Status:  http.StatusBadRequest,
			Message: "Restaurant ID is required",
		})
	}

	id, err := strconv.ParseUint(restaurantID, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.CustomResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid restaurant ID",
			Data:    err.Error(),
		})
	}

	err = h.service.Restaurant.DeleteRestaurant(c.Request().Context(), uint(id))
	if err != nil {
		h.logger.Error("Failed to delete restaurant:", err)
		return c.JSON(http.StatusInternalServerError, response.CustomResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to delete restaurant",
			Data:    err.Error(),
		})
	}

	return c.JSON(http.StatusOK, response.CustomResponse{
		Status:  http.StatusOK,
		Message: "Restaurant deleted successfully",
	})
}

func (h *RestaurantHandler) UpdateRestaurant(c echo.Context) error {
	var restaurant model.Restaurant
	if err := c.Bind(&restaurant); err != nil {
		return c.JSON(http.StatusBadRequest, response.CustomResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid request body",
			Data:    err.Error(),
		})
	}

	restaurantID := c.Param("id")
	if restaurantID == "" {
		return c.JSON(http.StatusBadRequest, response.CustomResponse{
			Status:  http.StatusBadRequest,
			Message: "Restaurant ID is required",
		})
	}

	id, err := strconv.ParseUint(restaurantID, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.CustomResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid restaurant ID",
			Data:    err.Error(),
		})
	}

	restaurant.ID = uint(id)

	updatedRestaurant, err := h.service.Restaurant.UpdateRestaurant(c.Request().Context(), &restaurant, uint(id))
	if err != nil {
		h.logger.Error("Failed to update restaurant:", err)
		return c.JSON(http.StatusInternalServerError, response.CustomResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to update restaurant",
			Data:    err.Error(),
		})
	}

	return c.JSON(http.StatusOK, response.CustomResponse{
		Status:  http.StatusOK,
		Message: "Restaurant updated successfully",
		Data:    updatedRestaurant,
	})
}

func (h *RestaurantHandler) GetRestaurants(c echo.Context) error {
	searchParams, err := h.service.Restaurant.RestaurantsSearchFormatting(model.NewParams(), c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.CustomResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed reading params",
			Data:    err.Error(),
		})
	}

	restaurants, err := h.service.Restaurant.GetRestaurants(c.Request().Context(), searchParams)
	if err != nil {
		h.logger.Error("Failed to get restaurants:", err)
		return c.JSON(http.StatusInternalServerError, response.CustomResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to get restaurants",
			Data:    err.Error(),
		})
	}

	return c.JSON(http.StatusOK, response.CustomResponse{
		Status:  http.StatusOK,
		Message: "Success",
		Data:    restaurants,
	})
}

func (h *RestaurantHandler) GetStatistics(c echo.Context) error {
	statistics, err := h.service.Restaurant.GetStatistics(c.Request().Context())
	if err != nil {
		h.logger.Error("Failed to get restaurants:", err)
		return c.JSON(http.StatusInternalServerError, response.CustomResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to get restaurants",
			Data:    err.Error(),
		})
	}

	return c.JSON(http.StatusOK, response.CustomResponse{
		Status:  http.StatusOK,
		Message: "Success",
		Data:    statistics,
	})
}

func (h *RestaurantHandler) GetRestaurantByID(c echo.Context) error {
	restaurantID := c.Param("id")
	if restaurantID == "" {
		return c.JSON(http.StatusBadRequest, response.CustomResponse{
			Status:  http.StatusBadRequest,
			Message: "Restaurant ID is required",
		})
	}

	id, err := strconv.ParseUint(restaurantID, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.CustomResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid restaurant ID",
			Data:    err.Error(),
		})
	}

	restaurant, err := h.service.Restaurant.GetRestaurantByID(c.Request().Context(), uint(id))
	if err != nil {
		h.logger.Error("Failed to get restaurant:", err)
		return c.JSON(http.StatusInternalServerError, response.CustomResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to get restaurant",
			Data:    err.Error(),
		})
	}

	return c.JSON(http.StatusOK, response.CustomResponse{
		Status:  http.StatusOK,
		Message: "Restaurant retrieved successfully",
		Data:    restaurant,
	})
}

func (h *RestaurantHandler) GetRestaurantOrders(c echo.Context) error {
	restaurantID := c.Param("id")
	if restaurantID == "" {
		return c.JSON(http.StatusBadRequest, response.CustomResponse{
			Status:  http.StatusBadRequest,
			Message: "Restaurant ID is required",
		})
	}

	id, err := strconv.ParseUint(restaurantID, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.CustomResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid restaurant ID",
			Data:    err.Error(),
		})
	}

	searchParams, err := h.service.Order.OrderSearchFormatting(model.NewParams(), c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.CustomResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed reading params",
			Data:    err.Error(),
		})
	}

	orders, err := h.service.Restaurant.GetRestaurantOrders(c.Request().Context(), uint(id), searchParams)
	if err != nil {
		h.logger.Error("Failed to get orders for restaurant:", err)
		return c.JSON(http.StatusInternalServerError, response.CustomResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to get orders for restaurant",
			Data:    err.Error(),
		})
	}

	return c.JSON(http.StatusOK, response.CustomResponse{
		Status:  http.StatusOK,
		Message: "Orders retrieved successfully",
		Data:    orders,
	})
}

//func (h *RestaurantHandler) FavoriteRestaurants(c echo.Context) error {
//	//TODO implement me
//	panic("implement me")
//}
