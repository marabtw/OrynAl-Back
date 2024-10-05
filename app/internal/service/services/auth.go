package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/alibekabdrakhman1/orynal/config"
	"github.com/alibekabdrakhman1/orynal/internal/model"
	"github.com/alibekabdrakhman1/orynal/internal/repository"
	"github.com/alibekabdrakhman1/orynal/pkg/utils"
	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
	"strconv"
	"time"
)

var ErrExpiredToken = errors.New("expiration date validation error")

func NewAuthService(repository *repository.Manager, config *config.Config, logger *zap.SugaredLogger) *AuthService {
	return &AuthService{repository: repository, config: config, logger: logger}
}

type AuthService struct {
	repository *repository.Manager
	config     *config.Config
	logger     *zap.SugaredLogger
}

func (s *AuthService) Login(ctx context.Context, login model.Login) (*model.JwtTokens, error) {
	user, err := s.repository.User.GetByEmail(ctx, login.Email)
	if err != nil {
		s.logger.Errorf("GetUser request err: %v", err)
		return nil, fmt.Errorf("GetUser request err: %w", err)
	}
	err = utils.CheckPassword(login.Password, user.Password)
	if err != nil {
		s.logger.Error("incorrect password")
		return nil, errors.New("incorrect password")
	}

	userClaim := model.UserClaim{
		Email:  user.Email,
		UserID: user.ID,
		Role:   user.Role,
	}

	tokens, err := s.generateToken(ctx, userClaim)
	if err != nil {
		s.logger.Errorf("generating token err: %v", err)
		return nil, fmt.Errorf("generating token err: %w", err)
	}

	res := &model.JwtTokens{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}
	return res, nil
}

func (s *AuthService) Register(ctx context.Context, user model.Register) (uint, error) {
	s.logger.Info(user)

	pass, err := utils.HashPassword(user.Password)
	if err != nil {
		s.logger.Error(err)
		return 0, err
	}

	req := &model.User{
		Name:     user.Name,
		Surname:  user.Surname,
		Email:    user.Email,
		Phone:    user.Phone,
		Role:     "user",
		Password: pass,
	}

	res, err := s.repository.User.Create(ctx, req)
	if err != nil {
		s.logger.Error(err)
		return 0, err
	}

	return res.ID, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*model.JwtTokens, error) {
	token, err := jwt.Parse(
		refreshToken,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(s.config.Auth.JwtSecretKey), nil
		},
	)

	if err != nil {
		var validationErr *jwt.ValidationError
		if errors.As(err, &validationErr) {
			if validationErr.Errors&jwt.ValidationErrorExpired > 0 {
				return nil, errors.New("expiration date validation error")
			}
		}
		s.logger.Error(err)
		return nil, fmt.Errorf("failed to parse jwt err: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		s.logger.Error(err)
		return nil, fmt.Errorf("unexpected type %T", claims)
	}
	fmt.Println(claims)
	user, err := s.repository.User.GetByEmail(ctx, claims["email"].(string))
	if err != nil {
		s.logger.Error(err)
		return nil, fmt.Errorf("GetUser request err: %w", err)
	}

	userClaim := model.UserClaim{

		UserID: user.ID,
		Role:   user.Role,
	}

	tokens, err := s.generateToken(ctx, userClaim)
	if err != nil {
		s.logger.Error(err)
		return nil, fmt.Errorf("generating token err: %w", err)
	}

	return tokens, nil
}

func (s *AuthService) generateToken(ctx context.Context, user model.UserClaim) (*model.JwtTokens, error) {
	accessTokenExpirationTime := time.Now().Add(time.Hour)
	refreshTokenExpirationTime := time.Now().Add(24 * time.Hour)

	accessTokenClaims := &model.JWTClaim{
		Email:  user.Email,
		UserID: user.UserID,
		Role:   user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: accessTokenExpirationTime.Unix(),
		},
	}

	secretKey := []byte(s.config.Auth.JwtSecretKey)
	accessClaimToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)

	accessTokenString, err := accessClaimToken.SignedString(secretKey)
	if err != nil {
		s.logger.Errorf("AccessToken: SignedStrign err: %v", err)
		return nil, fmt.Errorf("AccessToken: SignedString err: %w", err)
	}

	refreshTokenClaims := &model.RefreshJWTClaim{
		Email:  user.Email,
		UserID: user.UserID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: refreshTokenExpirationTime.Unix(),
		},
	}

	refreshClaimToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)

	refreshTokenString, err := refreshClaimToken.SignedString(secretKey)
	if err != nil {
		s.logger.Errorf("RefreshToken: SignedString err: %v", err)
		return nil, fmt.Errorf("RefreshToken: SignedString err: %w", err)
	}

	userToken := model.UserToken{
		UserID:       user.UserID,
		Role:         user.Role,
		Email:        user.Email,
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	}

	s.logger.Info(userToken)

	err = s.repository.UserToken.CreateUserToken(ctx, userToken)
	if err != nil {
		s.logger.Errorf("CreateUserToken err: %v", err)
		return nil, errors.New(fmt.Sprintf("CreateUserToken err: %v", err))
	}

	jwtToken := &model.JwtTokens{
		AccessToken:  userToken.AccessToken,
		RefreshToken: userToken.RefreshToken,
	}

	return jwtToken, nil
}

func (s *AuthService) GetJwtUserID(jwtToken string) (*model.ContextUserID, error) {
	token, err := jwt.Parse(
		jwtToken,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(s.config.JwtSecretKey), nil
		},
	)

	if err != nil {
		var validationErr *jwt.ValidationError
		if errors.As(err, &validationErr) {
			if validationErr.Errors&jwt.ValidationErrorExpired > 0 {
				return nil, ErrExpiredToken
			}
		}

		return nil, fmt.Errorf("failed to parse jwt err: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("unexpected type %T", claims)
	}

	user, err := s.getUserIDFromJwt(claims)
	if err != nil {
		return nil, fmt.Errorf("failed to get user from jwt err: %w", err)
	}
	return user, nil
}

func (s *AuthService) GetJwtUserRole(jwtToken string) (*model.ContextUserRole, error) {
	token, err := jwt.Parse(
		jwtToken,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(s.config.JwtSecretKey), nil
		},
	)

	if err != nil {
		var validationErr *jwt.ValidationError
		if errors.As(err, &validationErr) {
			if validationErr.Errors&jwt.ValidationErrorExpired > 0 {
				return nil, ErrExpiredToken
			}
		}

		return nil, fmt.Errorf("failed to parse jwt err: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("unexpected type %T", claims)
	}

	user, err := s.getUserRoleFromJwt(claims)
	if err != nil {
		return nil, fmt.Errorf("failed to get user from jwt err: %w", err)
	}
	return user, nil
}

func (s *AuthService) getUserIDFromJwt(claims jwt.MapClaims) (*model.ContextUserID, error) {
	user := &model.ContextUserID{}
	userId, ok := claims["user_id"]
	if !ok {
		return nil, fmt.Errorf("user is not exists in jwt")
	}

	userId = fmt.Sprintf("%.0f", userId)

	parsedUserId, err := strconv.ParseUint(userId.(string), 10, 32)
	if err != nil {
		return nil, fmt.Errorf("unexpected in userID value: %T", userId)
	}

	user.ID = uint(parsedUserId)

	return user, nil
}

func (s *AuthService) getUserRoleFromJwt(claims jwt.MapClaims) (*model.ContextUserRole, error) {
	user := &model.ContextUserRole{}
	role, ok := claims["role"]
	if !ok {
		return nil, fmt.Errorf("user is not exists in jwt")
	}

	user.Role = role.(string)

	return user, nil
}
