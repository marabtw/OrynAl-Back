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

type ReviewsHandler struct {
	service *service.Manager
	logger  *zap.SugaredLogger
}

func NewReviewsHandler(service *service.Manager, logger *zap.SugaredLogger) *ReviewsHandler {
	return &ReviewsHandler{
		service: service,
		logger:  logger,
	}
}

func (h *ReviewsHandler) DeleteReview(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, response.CustomResponse{
			Status:  http.StatusBadRequest,
			Message: "Review ID is required",
		})
	}

	reviewID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.CustomResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid review ID",
			Data:    err.Error(),
		})
	}

	err = h.service.Reviews.DeleteReview(c.Request().Context(), uint(reviewID))
	if err != nil {
		h.logger.Error("Failed to delete review:", err)
		return c.JSON(http.StatusInternalServerError, response.CustomResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to delete review",
			Data:    err.Error(),
		})
	}

	return c.JSON(http.StatusOK, response.CustomResponse{
		Status:  http.StatusOK,
		Message: "Review deleted successfully",
	})
}

func (h *ReviewsHandler) CreateReview(c echo.Context) error {
	var review model.RestaurantReview
	if err := c.Bind(&review); err != nil {
		return c.JSON(http.StatusBadRequest, response.CustomResponse{
			Status:  http.StatusBadRequest,
			Message: "Failed to parse review data",
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

	review.RestaurantID = uint(id)

	createdReview, err := h.service.Reviews.CreateReview(c.Request().Context(), &review)
	if err != nil {
		h.logger.Error("Failed to create review:", err)
		return c.JSON(http.StatusInternalServerError, response.CustomResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to create review",
			Data:    err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, response.CustomResponse{
		Status:  http.StatusCreated,
		Message: "Review created successfully",
		Data:    createdReview,
	})
}

func (h *ReviewsHandler) GetReviews(c echo.Context) error {
	searchParams, err := h.service.Reviews.ReviewsSearchFormatting(model.NewParams(), c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.CustomResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed reading params",
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

	reviews, err := h.service.Reviews.GetReviews(c.Request().Context(), uint(id), searchParams)
	if err != nil {
		h.logger.Error("Failed to get reviews:", err)
		return c.JSON(http.StatusInternalServerError, response.CustomResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to get reviews",
			Data:    err.Error(),
		})
	}

	return c.JSON(http.StatusOK, response.CustomResponse{
		Status:  http.StatusOK,
		Message: "Success",
		Data:    reviews,
	})
}
