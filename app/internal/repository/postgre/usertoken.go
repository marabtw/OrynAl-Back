package postgre

import (
	"context"
	"errors"
	"github.com/alibekabdrakhman1/orynal/internal/model"
	"gorm.io/gorm"
)

type UserTokenRepository struct {
	DB *gorm.DB
}

func NewUserTokenRepository(db *gorm.DB) *UserTokenRepository {
	return &UserTokenRepository{
		DB: db,
	}
}

func (r *UserTokenRepository) CreateUserToken(ctx context.Context, userToken model.UserToken) error {
	var existingToken model.UserToken
	result := r.DB.WithContext(ctx).Where("user_id = ?", userToken.UserID).First(&existingToken)

	if result.Error == nil {
		existingToken.AccessToken = userToken.AccessToken
		existingToken.RefreshToken = userToken.RefreshToken
		result = r.DB.WithContext(ctx).Save(&existingToken)
	} else if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		if err := r.DB.WithContext(ctx).Create(&userToken).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *UserTokenRepository) UpdateUserToken(ctx context.Context, userToken model.UserToken) error {
	if err := r.DB.WithContext(ctx).Save(&userToken).Error; err != nil {
		return err
	}
	return nil
}
