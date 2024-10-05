package handlers

import (
	"encoding/json"
	"github.com/alibekabdrakhman1/orynal/internal/model"
	"github.com/alibekabdrakhman1/orynal/internal/service"
	"github.com/alibekabdrakhman1/orynal/pkg/response"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
)

func NewUserHandler(service *service.Manager, logger *zap.SugaredLogger) *UserHandler {
	return &UserHandler{
		service: service,
		logger:  logger,
	}
}

type UserHandler struct {
	service *service.Manager
	logger  *zap.SugaredLogger
}

func (h *UserHandler) Profile(c echo.Context) error {
	user, err := h.service.User.Profile(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.CustomResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
			Data:    err,
		})
	}
	return c.JSON(http.StatusOK, response.CustomResponse{
		Status:  http.StatusOK,
		Message: "User profile retrieved successfully",
		Data:    user,
	})
}

func (h *UserHandler) UpdateProfile(c echo.Context) error {
	var user model.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, response.CustomResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid request body",
			Data:    err.Error(),
		})
	}

	updatedUser, err := h.service.User.Update(c.Request().Context(), &user)
	if err != nil {
		h.logger.Error("Failed to update user profile:", err)
		return c.JSON(http.StatusInternalServerError, response.CustomResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to update user profile",
			Data:    err.Error(),
		})
	}

	return c.JSON(http.StatusOK, response.CustomResponse{
		Status:  http.StatusOK,
		Message: "User profile updated successfully",
		Data:    updatedUser,
	})
}

func (h *UserHandler) ChangePassword(c echo.Context) error {
	var pass model.ChangePasswordRequest
	if err := c.Bind(&pass); err != nil {
		return c.JSON(http.StatusBadRequest, response.CustomResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid request body",
			Data:    err.Error(),
		})
	}

	err := h.service.User.ChangePassword(c.Request().Context(), &pass)
	if err != nil {
		h.logger.Error("Failed to update user password:", err)
		return c.JSON(http.StatusInternalServerError, response.CustomResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to update user password",
			Data:    err.Error(),
		})
	}

	return c.JSON(http.StatusOK, response.CustomResponse{
		Status:  http.StatusOK,
		Message: "User password updated successfully",
	})
}

func (h *UserHandler) DeleteProfile(c echo.Context) error {
	currentUser := c.Get("user").(*model.UserResponse)
	err := h.service.User.Delete(c.Request().Context(), currentUser.ID)
	if err != nil {
		h.logger.Error("Failed to delete user profile:", err)
		return c.JSON(http.StatusInternalServerError, response.CustomResponse{
			Status:  http.StatusInternalServerError,
			Message: "Failed to delete user profile",
			Data:    err.Error(),
		})
	}

	return c.JSON(http.StatusOK, response.CustomResponse{
		Status:  http.StatusOK,
		Message: "User profile deleted successfully",
	})
}
func (h *UserHandler) SignIn(c echo.Context) error {
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		h.logger.Error(err)
		return c.JSON(http.StatusBadRequest, response.CustomResponse{
			Status:  -1,
			Message: "request body reading error",
			Data:    err,
		})
	}
	var request model.Login

	err = json.Unmarshal(body, &request)
	if err != nil {
		h.logger.Error(err)
		return c.JSON(http.StatusBadRequest, response.CustomResponse{
			Status:  -1,
			Message: "register model unmarshalling error",
			Data:    err,
		})
	}

	//return c.JSON(http.StatusCreated, response.CustomResponse{
	//	Status:  0,
	//	Message: "OK",
	//	Data:    request,
	//})

	userToken, err := h.service.Auth.Login(c.Request().Context(), request)
	if err != nil {
		h.logger.Error(err)
		return c.JSON(http.StatusBadRequest, response.CustomResponse{
			Status:  -1,
			Message: "login error",
			Data:    err.Error(),
		})
	}

	res := model.JwtTokens{
		AccessToken:  userToken.AccessToken,
		RefreshToken: userToken.RefreshToken,
	}

	h.logger.Info(res)
	return c.JSON(http.StatusCreated, response.CustomResponse{
		Status:  0,
		Message: "OK",
		Data:    res,
	})
}

func (h *UserHandler) SignUp(c echo.Context) error {
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.CustomResponse{
			Status:  -1,
			Message: "Request Body reading error",
			Data:    err,
		})
	}

	var request model.Register
	err = json.Unmarshal(body, &request)
	if err != nil {
		h.logger.Error(err)
		return c.JSON(http.StatusBadRequest, response.CustomResponse{
			Status:  -1,
			Message: "Register model unmarshalling error",
			Data:    err,
		})
	}

	userId, err := h.service.Auth.Register(c.Request().Context(), request)
	if err != nil {
		h.logger.Error(err)
		return c.JSON(http.StatusBadRequest, response.CustomResponse{
			Status:  -1,
			Message: "Register error",
			Data:    err,
		})
	}
	h.logger.Info(userId)
	return c.JSON(http.StatusCreated, response.CustomResponse{
		Status:  0,
		Message: "OK",
		Data: response.IDResponse{
			ID: userId,
		},
	})
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (h *UserHandler) RefreshToken(c echo.Context) error {
	var r refreshRequest
	err := c.Bind(&r)
	if err != nil {
		h.logger.Error(err)
		return c.JSON(http.StatusBadRequest, response.CustomResponse{
			Status:  -1,
			Message: "Request Body reading error",
			Data:    err,
		})
	}

	tokens, err := h.service.Auth.RefreshToken(c.Request().Context(), r.RefreshToken)
	if err != nil {
		h.logger.Error(err)
		return c.JSON(http.StatusBadRequest, response.CustomResponse{
			Status:  -1,
			Message: "Refresh Token error",
			Data:    err,
		})
	}

	h.logger.Info(tokens)
	return c.JSON(http.StatusCreated, response.CustomResponse{
		Status:  0,
		Message: "OK",
		Data:    tokens,
	})
}
