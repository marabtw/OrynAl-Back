package handlers

import (
	"github.com/alibekabdrakhman1/orynal/internal/model"
	"github.com/alibekabdrakhman1/orynal/internal/service"
	"github.com/alibekabdrakhman1/orynal/pkg/response"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

func NewOrderHandler(service *service.Manager, logger *zap.SugaredLogger) *OrderHandler {
	return &OrderHandler{
		service: service,
		logger:  logger,
	}
}

type OrderHandler struct {
	service *service.Manager
	logger  *zap.SugaredLogger
}

func (h *OrderHandler) CreateOrder(c echo.Context) error {
	var order model.OrderRequest
	if err := c.Bind(&order); err != nil {
		return c.JSON(http.StatusBadRequest, response.CustomResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid request body",
			Data:    err.Error(),
		})
	}

	createdOrder, err := h.service.Order.Create(c.Request().Context(), &order)
	if err != nil {
		h.logger.Error("Failed to create order:", err)
		return c.JSON(http.StatusInternalServerError, response.CustomResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to create order",
			Data:    err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, response.CustomResponse{
		Status:  http.StatusCreated,
		Message: "Order created successfully",
		Data:    createdOrder,
	})
}

func (h *OrderHandler) DeleteOrder(c echo.Context) error {
	orderID := c.Param("id")
	if orderID == "" {
		return c.JSON(http.StatusBadRequest, response.CustomResponse{
			Status:  http.StatusBadRequest,
			Message: "Order ID is required",
		})
	}

	id, err := strconv.ParseUint(orderID, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.CustomResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid order ID",
			Data:    err.Error(),
		})
	}

	err = h.service.Order.Delete(c.Request().Context(), uint(id))
	if err != nil {
		h.logger.Error("Failed to delete order:", err)
		return c.JSON(http.StatusInternalServerError, response.CustomResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to delete order",
			Data:    err.Error(),
		})
	}

	return c.JSON(http.StatusOK, response.CustomResponse{
		Status:  http.StatusOK,
		Message: "Order deleted successfully",
	})
}

func (h *OrderHandler) UpdateOrder(c echo.Context) error {
	var order model.Order
	if err := c.Bind(&order); err != nil {
		return c.JSON(http.StatusBadRequest, response.CustomResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid request body",
			Data:    err.Error(),
		})
	}

	orderID := c.Param("id")
	if orderID == "" {
		return c.JSON(http.StatusBadRequest, response.CustomResponse{
			Status:  http.StatusBadRequest,
			Message: "Order ID is required",
		})
	}

	id, err := strconv.ParseUint(orderID, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.CustomResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid order ID",
			Data:    err.Error(),
		})
	}

	order.ID = uint(id)

	updatedOrder, err := h.service.Order.Update(c.Request().Context(), uint(id), &order)
	if err != nil {
		h.logger.Error("Failed to update order:", err)
		return c.JSON(http.StatusInternalServerError, response.CustomResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to update order",
			Data:    err.Error(),
		})
	}

	return c.JSON(http.StatusOK, response.CustomResponse{
		Status:  http.StatusOK,
		Message: "Order updated successfully",
		Data:    updatedOrder,
	})
}

func (h *OrderHandler) GetOrder(c echo.Context) error {
	orderID := c.Param("id")
	if orderID == "" {
		return c.JSON(http.StatusBadRequest, response.CustomResponse{
			Status:  http.StatusBadRequest,
			Message: "Order ID is required",
		})
	}

	id, err := strconv.ParseUint(orderID, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.CustomResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid order ID",
			Data:    err.Error(),
		})
	}

	order, err := h.service.Order.GetByID(c.Request().Context(), uint(id))
	if err != nil {
		h.logger.Error("Failed to get order:", err)
		return c.JSON(http.StatusInternalServerError, response.CustomResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to get order",
			Data:    err.Error(),
		})
	}

	return c.JSON(http.StatusOK, response.CustomResponse{
		Status:  http.StatusOK,
		Message: "Order retrieved successfully",
		Data:    order,
	})
}

func (h *OrderHandler) GetAllOrders(c echo.Context) error {
	searchParams, err := h.service.Order.OrderSearchFormatting(model.NewParams(), c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.CustomResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed reading params",
			Data:    err.Error(),
		})
	}

	orders, err := h.service.Order.GetAllOrders(c.Request().Context(), searchParams)
	if err != nil {
		h.logger.Error("Failed to get all orders:", err)
		return c.JSON(http.StatusInternalServerError, response.CustomResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to get all orders",
			Data:    err.Error(),
		})
	}

	return c.JSON(http.StatusOK, response.CustomResponse{
		Status:  http.StatusOK,
		Message: "All orders retrieved successfully",
		Data:    orders,
	})
}
