package handlers

import (
	"fmt"
	"github.com/alibekabdrakhman1/orynal/internal/model"
	"github.com/alibekabdrakhman1/orynal/internal/service"
	"github.com/alibekabdrakhman1/orynal/pkg/response"
	"github.com/alibekabdrakhman1/orynal/pkg/utils"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
)

func NewAdminHandler(service *service.Manager, logger *zap.SugaredLogger) *AdminHandler {
	return &AdminHandler{
		service: service,
		logger:  logger,
	}
}

type AdminHandler struct {
	service *service.Manager
	logger  *zap.SugaredLogger
}

func (h *AdminHandler) CreateService(c echo.Context) error {
	var createService model.Service
	if err := c.Bind(&createService); err != nil {
		return c.JSON(http.StatusBadRequest, response.CustomResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid request body",
			Data:    err.Error(),
		})
	}

	services, err := h.service.Restaurant.CreateService(c.Request().Context(), &createService)
	if err != nil {
		h.logger.Error("Failed to create owner:", err)
		return c.JSON(http.StatusInternalServerError, response.CustomResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to create owner",
			Data:    err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, response.CustomResponse{
		Status:  http.StatusCreated,
		Message: "Success",
		Data:    services,
	})
}

func (h *AdminHandler) DeleteService(c echo.Context) error {
	serviceID, err := utils.ConvertIdToUint(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.CustomResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid service ID",
			Data:    err.Error(),
		})
	}

	err = h.service.Restaurant.DeleteService(c.Request().Context(), serviceID)
	if err != nil {
		h.logger.Error("Failed to delete service:", err)
		return c.JSON(http.StatusInternalServerError, response.CustomResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to delete service",
			Data:    err.Error(),
		})
	}

	return c.JSON(http.StatusOK, response.CustomResponse{
		Status:  http.StatusOK,
		Message: "Success",
		Data:    nil,
	})
}

func (h *AdminHandler) UpdateService(c echo.Context) error {
	var updateService model.Service
	if err := c.Bind(&updateService); err != nil {
		return c.JSON(http.StatusBadRequest, response.CustomResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid request body",
			Data:    err.Error(),
		})
	}

	services, err := h.service.Restaurant.UpdateService(c.Request().Context(), &updateService)
	if err != nil {
		h.logger.Error("Failed to create owner:", err)
		return c.JSON(http.StatusInternalServerError, response.CustomResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to create owner",
			Data:    err.Error(),
		})
	}

	return c.JSON(http.StatusOK, response.CustomResponse{
		Status:  http.StatusOK,
		Message: "Success",
		Data:    services,
	})
}

func (h *AdminHandler) GetClients(c echo.Context) error {
	searchParams, err := h.service.User.UserSearchFormatting(model.NewParams(), c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.CustomResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed reading params",
			Data:    err.Error(),
		})
	}

	clients, err := h.service.User.GetAllClients(c.Request().Context(), searchParams)
	if err != nil {
		h.logger.Error("Failed to get clients:", err)
		return c.JSON(http.StatusInternalServerError, response.CustomResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to get clients",
			Data:    err.Error(),
		})
	}
	return c.JSON(http.StatusOK, response.CustomResponse{
		Status:  http.StatusOK,
		Message: "Success",
		Data:    clients,
	})
}

func (h *AdminHandler) GetClient(c echo.Context) error {
	clientID, err := utils.ConvertIdToUint(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.CustomResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid client id",
			Data:    err.Error(),
		})
	}

	client, err := h.service.User.GetByID(c.Request().Context(), clientID)
	if err != nil {
		h.logger.Error("Failed to get client:", err)
		return c.JSON(http.StatusInternalServerError, response.CustomResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to get client",
			Data:    err.Error(),
		})
	}

	return c.JSON(http.StatusOK, response.CustomResponse{
		Status:  http.StatusOK,
		Message: "Success",
		Data:    client,
	})
}

func (h *AdminHandler) DeleteClient(c echo.Context) error {
	clientID, err := utils.ConvertIdToUint(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.CustomResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid client id",
			Data:    err.Error(),
		})
	}

	err = h.service.User.Delete(c.Request().Context(), clientID)
	if err != nil {
		h.logger.Error("Failed to delete client:", err)
		return c.JSON(http.StatusInternalServerError, response.CustomResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to delete client",
			Data:    err.Error(),
		})
	}

	return c.JSON(http.StatusOK, response.CustomResponse{
		Status:  http.StatusOK,
		Message: "Deleted",
	})
}

func (h *AdminHandler) GetOwners(c echo.Context) error {
	searchParams, err := h.service.User.UserSearchFormatting(model.NewParams(), c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.CustomResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed reading params",
			Data:    err.Error(),
		})
	}

	owners, err := h.service.User.GetAllOwners(c.Request().Context(), searchParams)
	if err != nil {
		h.logger.Error("Failed to get owners:", err)
		return c.JSON(http.StatusInternalServerError, response.CustomResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to get owners",
			Data:    err.Error(),
		})
	}

	return c.JSON(http.StatusOK, response.CustomResponse{
		Status:  http.StatusOK,
		Message: "Success",
		Data:    owners,
	})
}

func (h *AdminHandler) CreateOwner(c echo.Context) error {
	var owner model.User
	if err := c.Bind(&owner); err != nil {
		return c.JSON(http.StatusBadRequest, response.CustomResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid request body",
			Data:    err.Error(),
		})
	}

	createdOwner, err := h.service.User.CreateOwner(c.Request().Context(), &owner)
	if err != nil {
		h.logger.Error("Failed to create owner:", err)
		return c.JSON(http.StatusInternalServerError, response.CustomResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to create owner",
			Data:    err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, response.CustomResponse{
		Status:  http.StatusCreated,
		Message: "Success",
		Data:    createdOwner,
	})
}

func (h *AdminHandler) DeleteOwner(c echo.Context) error {
	ownerID, err := utils.ConvertIdToUint(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.CustomResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid owner ID",
			Data:    err.Error(),
		})
	}

	err = h.service.User.Delete(c.Request().Context(), ownerID)
	if err != nil {
		h.logger.Error("Failed to delete owner:", err)
		return c.JSON(http.StatusInternalServerError, response.CustomResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to delete owner",
			Data:    err.Error(),
		})
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *AdminHandler) GetRestaurants(c echo.Context) error {
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

func (h *AdminHandler) GetRestaurant(c echo.Context) error {
	restaurantID, err := utils.ConvertIdToUint(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.CustomResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid restaurant ID",
			Data:    err.Error(),
		})
	}

	restaurant, err := h.service.Restaurant.GetRestaurantByID(c.Request().Context(), restaurantID)
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
		Message: "Success",
		Data:    restaurant,
	})
}

func (h *AdminHandler) CreateRestaurant(c echo.Context) error {
	var restaurant model.Restaurant
	if err := c.Bind(&restaurant); err != nil {
		return c.JSON(http.StatusBadRequest, response.CustomResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid request body",
			Data:    err.Error(),
		})
	}

	fmt.Println(restaurant)

	createdRestaurant, err := h.service.Restaurant.CreateRestaurant(c.Request().Context(), &restaurant)
	if err != nil {
		h.logger.Error("Failed to create restaurant:", err)
		return c.JSON(http.StatusInternalServerError, response.CustomResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to create restaurant",
			Data:    err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, response.CustomResponse{
		Status:  http.StatusCreated,
		Message: "Success",
		Data:    createdRestaurant,
	})
}

func (h *AdminHandler) DeleteRestaurant(c echo.Context) error {
	restaurantID, err := utils.ConvertIdToUint(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.CustomResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid restaurant ID",
			Data:    err.Error(),
		})
	}

	err = h.service.Restaurant.DeleteRestaurant(c.Request().Context(), restaurantID)
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
		Message: "Success",
		Data:    nil,
	})
}

func (h *AdminHandler) UpdateRestaurant(c echo.Context) error {
	restaurantID, err := utils.ConvertIdToUint(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.CustomResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid restaurant ID",
			Data:    err.Error(),
		})
	}

	var updatedRestaurant model.Restaurant
	if err := c.Bind(&updatedRestaurant); err != nil {
		return c.JSON(http.StatusBadRequest, response.CustomResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid request body",
			Data:    err.Error(),
		})
	}

	updatedRestaurant.ID = restaurantID
	restaurant, err := h.service.Restaurant.UpdateRestaurant(c.Request().Context(), &updatedRestaurant, restaurantID)
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
		Message: "Success",
		Data:    restaurant,
	})
}
