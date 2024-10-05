package postgre

import (
	"context"
	"errors"
	"github.com/alibekabdrakhman1/orynal/internal/model"
	"github.com/alibekabdrakhman1/orynal/pkg/enums"
	"github.com/alibekabdrakhman1/orynal/pkg/utils"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) (*model.UserResponse, error) {
	result := r.DB.WithContext(ctx).Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, errors.New("unable to create user")
	}

	return &model.UserResponse{
		ID:      user.ID,
		Name:    user.Name,
		Surname: user.Surname,
		Email:   user.Email,
		Phone:   user.Phone,
		Role:    user.Role,
	}, nil
}

func (r *UserRepository) Update(ctx context.Context, user *model.User) (*model.UserResponse, error) {
	var oldUser model.User
	if err := r.DB.WithContext(ctx).First(&oldUser, user.ID).Error; err != nil {
		return nil, err
	}

	if oldUser.Role != user.Role {
		return nil, errors.New("you cannot change user role")
	}

	if err := r.DB.WithContext(ctx).Model(&oldUser).Updates(user).Error; err != nil {
		return nil, err
	}

	return &model.UserResponse{
		ID:      oldUser.ID,
		Name:    oldUser.Name,
		Surname: oldUser.Surname,
		Email:   oldUser.Email,
		Phone:   oldUser.Phone,
		Role:    oldUser.Role,
	}, nil
}

func (r *UserRepository) ChangePassword(ctx context.Context, id uint, pass *model.ChangePasswordRequest) error {
	var oldUser model.User
	if err := r.DB.WithContext(ctx).First(&oldUser, id).Error; err != nil {
		return err
	}

	if utils.CheckPassword(pass.OldPassword, oldUser.Password) != nil {
		return errors.New("wrong old password")
	}

	oldUser.Password = pass.NewPassword

	if err := r.DB.WithContext(ctx).Save(&oldUser).Error; err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id uint) error {
	var user model.User
	if err := r.DB.WithContext(ctx).First(&user, id).Error; err != nil {
		return err
	}

	if err := r.DB.WithContext(ctx).Delete(&user).Error; err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uint) (*model.UserResponse, error) {
	var user model.User
	if err := r.DB.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, err
	}

	return &model.UserResponse{
		ID:      user.ID,
		Name:    user.Name,
		Surname: user.Surname,
		Email:   user.Email,
		Phone:   user.Phone,
		Role:    user.Role,
	}, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user *model.User
	if err := r.DB.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetAllClients(ctx context.Context, params *model.Params) (*model.ListResponse, error) {
	var clients []model.User
	var totalItems int64

	countQuery := r.DB.WithContext(ctx).Where("role = ?", enums.User)
	if params.Query != "" {
		countQuery = countQuery.Where("LOWER(name) LIKE LOWER(?)", "%"+params.Query+"%")
	}
	if err := countQuery.Model(&model.User{}).Count(&totalItems).Error; err != nil {
		return nil, err
	}

	if int64(params.Offset) >= totalItems {
		return nil, errors.New("offset exceeds total items")
	}

	query := r.DB.WithContext(ctx).Where("role = ?", enums.User).
		Limit(params.Limit).
		Offset(params.Offset)

	if params.Query != "" {
		query = query.Where("LOWER(name) LIKE LOWER(?)", "%"+params.Query+"%")
	}

	if params.Order != nil && params.SortVector != nil {
		query.Order(params.Order.(string) + params.SortVector.(string))
	}

	if err := query.Find(&clients).Error; err != nil {
		return nil, err
	}

	var userResponses []model.UserResponse
	for _, user := range clients {
		userResponses = append(userResponses, model.UserResponse{
			ID:      user.ID,
			Name:    user.Name,
			Surname: user.Surname,
			Email:   user.Email,
			Phone:   user.Phone,
			Role:    user.Role,
		})
	}

	return &model.ListResponse{
		Items:        userResponses,
		ItemsPerPage: params.Limit,
		PageIndex:    params.PageIndex,
		TotalItems:   int(totalItems),
	}, nil
}

func (r *UserRepository) GetAllOwners(ctx context.Context, params *model.Params) (*model.ListResponse, error) {
	var owners []model.User
	var totalItems int64

	countQuery := r.DB.WithContext(ctx).Where("role = ?", enums.Owner)
	if params.Query != "" {
		countQuery = countQuery.Where("LOWER(name) LIKE LOWER(?)", "%"+params.Query+"%")
	}
	if err := countQuery.Model(&model.User{}).Count(&totalItems).Error; err != nil {
		return nil, err
	}

	if int64(params.Offset) >= totalItems {
		return nil, errors.New("offset exceeds total items")
	}

	query := r.DB.WithContext(ctx).Where("role = ?", enums.Owner).
		Limit(params.Limit).
		Offset(params.Offset)

	if params.Query != "" {
		query = query.Where("LOWER(name) LIKE LOWER(?)", "%"+params.Query+"%")
	}

	if params.Order != nil && params.SortVector != nil {
		query.Order(params.Order.(string) + params.SortVector.(string))
	}

	if err := query.Find(&owners).Error; err != nil {
		return nil, err
	}

	var userResponses []model.UserResponse
	for _, user := range owners {
		userResponses = append(userResponses, model.UserResponse{
			ID:      user.ID,
			Name:    user.Name,
			Surname: user.Surname,
			Email:   user.Email,
			Phone:   user.Phone,
			Role:    user.Role,
		})
	}

	return &model.ListResponse{
		Items:        userResponses,
		ItemsPerPage: params.Limit,
		PageIndex:    params.PageIndex,
		TotalItems:   int(totalItems),
	}, nil
}
