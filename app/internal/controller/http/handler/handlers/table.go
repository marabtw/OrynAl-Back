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

func NewTableHandler(service *service.Manager, logger *zap.SugaredLogger) *TableHandler {
	return &TableHandler{
		service: service,
		logger:  logger,
	}
}

type TableHandler struct {
	service *service.Manager
	logger  *zap.SugaredLogger
}

func (h *TableHandler) GetTableCategories(c echo.Context) error {
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

	types, err := h.service.Table.GetTableCategories(c.Request().Context(), uint(id))
	if err != nil {
		h.logger.Error("Failed to get tables for restaurant:", err)
		return c.JSON(http.StatusInternalServerError, response.CustomResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to get tables for restaurant",
			Data:    err.Error(),
		})
	}

	return c.JSON(http.StatusOK, response.CustomResponse{
		Status:  http.StatusOK,
		Message: "Types retrieved successfully",
		Data:    types,
	})
}

func (h *TableHandler) GetRestaurantTables(c echo.Context) error {
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

	searchParams, err := h.service.Table.TablesSearchFormatting(model.NewParams(), c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.CustomResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed reading params",
			Data:    err.Error(),
		})
	}

	tables, err := h.service.Table.GetRestaurantTables(c.Request().Context(), uint(id), searchParams)
	if err != nil {
		h.logger.Error("Failed to get tables for restaurant:", err)
		return c.JSON(http.StatusInternalServerError, response.CustomResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to get tables for restaurant",
			Data:    err.Error(),
		})
	}

	return c.JSON(http.StatusOK, response.CustomResponse{
		Status:  http.StatusOK,
		Message: "Tables retrieved successfully",
		Data:    tables,
	})
}

func (h *TableHandler) GetRestaurantTable(c echo.Context) error {
	restaurantID := c.Param("id")
	if restaurantID == "" {
		return c.JSON(http.StatusBadRequest, response.CustomResponse{
			Status:  http.StatusBadRequest,
			Message: "Restaurant ID is required",
		})
	}

	restaurantId, err := strconv.ParseUint(restaurantID, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.CustomResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid restaurant ID",
			Data:    err.Error(),
		})
	}

	tableID := c.Param("table_id")
	if tableID == "" {
		return c.JSON(http.StatusBadRequest, response.CustomResponse{
			Status:  http.StatusBadRequest,
			Message: "Table ID is required",
		})
	}

	tableId, err := strconv.ParseUint(tableID, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.CustomResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid table ID",
			Data:    err.Error(),
		})
	}

	table, err := h.service.Table.GetRestaurantTable(c.Request().Context(), uint(restaurantId), uint(tableId))
	if err != nil {
		h.logger.Error("Failed to get table:", err)
		return c.JSON(http.StatusInternalServerError, response.CustomResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to get table",
			Data:    err.Error(),
		})
	}

	return c.JSON(http.StatusOK, response.CustomResponse{
		Status:  http.StatusOK,
		Message: "Table retrieved successfully",
		Data:    table,
	})
}

func (h *TableHandler) CreateRestaurantTable(c echo.Context) error {
	restaurantID := c.Param("id")
	if restaurantID == "" {
		return c.JSON(http.StatusBadRequest, response.CustomResponse{
			Status:  http.StatusBadRequest,
			Message: "Restaurant ID is required",
		})
	}

	restaurantId, err := strconv.ParseUint(restaurantID, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.CustomResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid restaurant ID",
			Data:    err.Error(),
		})
	}
	var table model.Table
	if err := c.Bind(&table); err != nil {
		return c.JSON(http.StatusBadRequest, response.CustomResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid request body",
			Data:    err.Error(),
		})
	}

	createdTable, err := h.service.Table.CreateRestaurantTable(c.Request().Context(), uint(restaurantId), &table)
	if err != nil {
		h.logger.Error("Failed to create table:", err)
		return c.JSON(http.StatusInternalServerError, response.CustomResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to create table",
			Data:    err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, response.CustomResponse{
		Status:  http.StatusCreated,
		Message: "Table created successfully",
		Data:    createdTable,
	})
}

func (h *TableHandler) UpdateRestaurantTable(c echo.Context) error {
	var table model.Table
	if err := c.Bind(&table); err != nil {
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

	restaurantId, err := strconv.ParseUint(restaurantID, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.CustomResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid restaurant ID",
			Data:    err.Error(),
		})
	}

	tableID := c.Param("table_id")
	if tableID == "" {
		return c.JSON(http.StatusBadRequest, response.CustomResponse{
			Status:  http.StatusBadRequest,
			Message: "Table ID is required",
		})
	}

	tableId, err := strconv.ParseUint(tableID, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.CustomResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid table ID",
			Data:    err.Error(),
		})
	}
	table.ID = uint(tableId)

	updatedTable, err := h.service.Table.UpdateRestaurantTable(c.Request().Context(), uint(restaurantId), &table)
	if err != nil {
		h.logger.Error("Failed to update table:", err)
		return c.JSON(http.StatusInternalServerError, response.CustomResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to update table",
			Data:    err.Error(),
		})
	}

	return c.JSON(http.StatusOK, response.CustomResponse{
		Status:  http.StatusOK,
		Message: "Table updated successfully",
		Data:    updatedTable,
	})
}

func (h *TableHandler) DeleteRestaurantTable(c echo.Context) error {
	restaurantID := c.Param("id")
	if restaurantID == "" {
		return c.JSON(http.StatusBadRequest, response.CustomResponse{
			Status:  http.StatusBadRequest,
			Message: "Restaurant ID is required",
		})
	}

	restaurantId, err := strconv.ParseUint(restaurantID, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.CustomResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid restaurant ID",
			Data:    err.Error(),
		})
	}

	tableID := c.Param("table_id")
	if tableID == "" {
		return c.JSON(http.StatusBadRequest, response.CustomResponse{
			Status:  http.StatusBadRequest,
			Message: "Table ID is required",
		})
	}

	tableId, err := strconv.ParseUint(tableID, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.CustomResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid table ID",
			Data:    err.Error(),
		})
	}

	err = h.service.Table.DeleteRestaurantTable(c.Request().Context(), uint(restaurantId), uint(tableId))
	if err != nil {
		h.logger.Error("Failed to delete table:", err)
		return c.JSON(http.StatusInternalServerError, response.CustomResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to delete table",
			Data:    err.Error(),
		})
	}

	return c.JSON(http.StatusOK, response.CustomResponse{
		Status:  http.StatusOK,
		Message: "Table deleted successfully",
	})
}

func (h *TableHandler) GetAvailableTime(c echo.Context) error {
	//TODO implement me
	panic("implement me")
}
