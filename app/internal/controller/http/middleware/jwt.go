package middleware

import (
	"context"
	"errors"
	"fmt"
	"github.com/alibekabdrakhman1/orynal/internal/model"
	"github.com/alibekabdrakhman1/orynal/internal/service/services"
	"github.com/alibekabdrakhman1/orynal/pkg/enums"
	"github.com/alibekabdrakhman1/orynal/pkg/utils"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

const (
	AuthorizationHeaderKey = "Authorization"
)

type JWTAuth struct {
	jwtKey      []byte
	AuthService services.IAuthService
	logger      *zap.SugaredLogger
}

func NewJWTAuth(jwtKey []byte, service services.IAuthService, logger *zap.SugaredLogger) *JWTAuth {
	return &JWTAuth{jwtKey: jwtKey, AuthService: service, logger: logger}
}

func (m *JWTAuth) RoleToCtx(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		jwtToken, err := m.getTokenFromHeader(c.Request())
		if err != nil {
			return err
		}

		contextUserId, _ := m.AuthService.GetJwtUserID(jwtToken)
		contextUserRole, _ := m.AuthService.GetJwtUserRole(jwtToken)
		ctx := context.WithValue(c.Request().Context(), model.ContextUserIDKey, contextUserId)
		ctx = context.WithValue(ctx, model.ContextUserRoleKey, contextUserRole)

		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}

func (m *JWTAuth) ValidateAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		jwtToken, err := m.getTokenFromHeader(c.Request())
		if err != nil {
			return err
		}

		contextUserId, err := m.AuthService.GetJwtUserID(jwtToken)
		if err != nil {
			if !errors.Is(err, services.ErrExpiredToken) {
				m.logger.Errorf("failed to GetJwtUser err: %v", err)
			}

			return echo.NewHTTPError(http.StatusUnauthorized, fmt.Sprintf("Unauthorized"))
		} else {
			contextUserRole, err := m.AuthService.GetJwtUserRole(jwtToken)
			if err != nil {
				m.logger.Errorf("failed to GetJwtUser err: %v", err)
			}

			ctx := context.WithValue(c.Request().Context(), model.ContextUserIDKey, contextUserId)
			ctx = context.WithValue(ctx, model.ContextUserRoleKey, contextUserRole)

			c.SetRequest(c.Request().WithContext(ctx))
		}
		return next(c)
	}
}

func (m *JWTAuth) ValidateAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := m.validateRole(c, enums.Admin); err != nil {
			return err
		}
		return next(c)
	}
}

func (m *JWTAuth) ValidateOwner(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := m.validateRole(c, enums.Owner); err != nil {
			return err
		}
		return next(c)
	}
}

func (m *JWTAuth) ValidateUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := m.validateRole(c, enums.User); err != nil {
			return err
		}
		return next(c)
	}
}

func (m *JWTAuth) validateRole(c echo.Context, expectedRole string) error {
	role, err := utils.GetRoleFromContext(c.Request().Context())
	if err != nil {
		m.logger.Error(err)
		return err
	}
	if role != expectedRole {
		m.logger.Warn(fmt.Sprintf("you are not %v", expectedRole))
		return echo.NewHTTPError(http.StatusForbidden, fmt.Sprintf("not permitted"))
	}
	return nil
}

func (m *JWTAuth) getTokenFromHeader(r *http.Request) (string, error) {
	if _, ok := r.Header[AuthorizationHeaderKey]; !ok {
		m.logger.Warn("'Authorization' key missing from headers")
		return "", echo.NewHTTPError(http.StatusForbidden, errors.New("authorization' key missing from headers"))
	}

	jwtToken := r.Header.Get(AuthorizationHeaderKey)

	if !(len(jwtToken) > 7 && strings.ToUpper(jwtToken[0:6]) == "BEARER") {
		m.logger.Warn(fmt.Sprintf(
			"failed to getTokenFromHeader invalidToken: %s",
			r.Header.Get(AuthorizationHeaderKey),
		))

		return "", echo.NewHTTPError(http.StatusForbidden, fmt.Errorf("failed to getTokenFromHeader invalidToken: %s", r.Header.Get(AuthorizationHeaderKey)))
	}

	return jwtToken[7:], nil
}
