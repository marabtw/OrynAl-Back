package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/alibekabdrakhman1/orynal/config"
	"github.com/alibekabdrakhman1/orynal/internal/model"
	"github.com/alibekabdrakhman1/orynal/internal/repository"
	"github.com/alibekabdrakhman1/orynal/internal/service/infrastructure"
	"github.com/alibekabdrakhman1/orynal/pkg/enums"
	"github.com/alibekabdrakhman1/orynal/pkg/utils"
	"go.uber.org/zap"
)

func NewUserService(repository *repository.Manager, config *config.Config, logger *zap.SugaredLogger) *UserService {
	return &UserService{repository: repository, config: config, logger: logger, FormatParams: infrastructure.NewFormatParams()}
}

type UserService struct {
	repository *repository.Manager
	config     *config.Config
	logger     *zap.SugaredLogger
	FormatParams
}

func (s *UserService) Create(ctx context.Context, user *model.User) (*model.UserResponse, error) {
	if user.Role == enums.Admin {
		return nil, errors.New("permission denied")
	}

	role, err := utils.GetRoleFromContext(ctx)
	if err == nil {
		if role != enums.Admin {
			return nil, errors.New("permission denied")
		} else {
			if user.Role != enums.Owner {
				return nil, errors.New("admin can create only restaurant owners")
			}
		}
	}

	return s.repository.User.Create(ctx, user)
}

func (s *UserService) CreateOwner(ctx context.Context, user *model.User) (*model.UserResponse, error) {
	pass, err := utils.HashPassword(user.Password)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	user.Password = pass

	role, err := utils.GetRoleFromContext(ctx)
	if err == nil {
		if role != enums.Admin {
			return nil, errors.New("permission denied")
		} else {
			if user.Role != enums.Owner {
				return nil, errors.New("admin can create only restaurant owners")
			}
		}
	}

	return s.repository.User.Create(ctx, user)
}

func (s *UserService) Update(ctx context.Context, user *model.User) (*model.UserResponse, error) {
	if user.Role == enums.Admin {
		return nil, errors.New("permission denied")
	}

	role, err := utils.GetRoleFromContext(ctx)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	id, err := utils.GetIDFromContext(ctx)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	if role != enums.Admin || user.ID != id {
		return nil, errors.New("permission denied")
	}

	return s.repository.User.Update(ctx, user)
}

func (s *UserService) ChangePassword(ctx context.Context, pass *model.ChangePasswordRequest) error {
	id, err := utils.GetIDFromContext(ctx)
	if err != nil {
		s.logger.Error(err)
		return err
	}

	if pass.OldPassword == pass.NewPassword {
		return errors.New("old password is equal new password")
	}

	pass.NewPassword, err = utils.HashPassword(pass.NewPassword)
	if err != nil {
		return fmt.Errorf("error hashing password: %v", err)
	}

	return s.repository.User.ChangePassword(ctx, id, pass)

}

func (s *UserService) Delete(ctx context.Context, id uint) error {
	user, err := s.repository.User.GetByID(ctx, id)
	if user.Role == enums.Admin {
		return errors.New("permission denied")
	}

	role, err := utils.GetRoleFromContext(ctx)
	if err != nil {
		s.logger.Error(err)
		return err
	}

	ctxID, err := utils.GetIDFromContext(ctx)
	if err != nil {
		s.logger.Error(err)
		return err
	}

	if role != enums.Admin && user.ID != ctxID {
		return errors.New("permission denied")
	}

	return s.repository.User.Delete(ctx, id)
}

func (s *UserService) Profile(ctx context.Context) (*model.UserResponse, error) {
	id, err := utils.GetIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	return s.repository.User.GetByID(ctx, id)
}

func (s *UserService) GetByID(ctx context.Context, id uint) (*model.UserResponse, error) {
	return s.repository.User.GetByID(ctx, id)
}

func (s *UserService) GetAllClients(ctx context.Context, params *model.Params) (*model.ListResponse, error) {
	role, err := utils.GetRoleFromContext(ctx)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	if role != enums.Admin {
		return nil, errors.New("permission denied")
	}

	list, err := s.repository.User.GetAllClients(ctx, params)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	list.TotalPages = list.TotalItems / list.ItemsPerPage

	return list, err
}

func (s *UserService) GetAllOwners(ctx context.Context, params *model.Params) (*model.ListResponse, error) {
	role, err := utils.GetRoleFromContext(ctx)
	if err != nil {
		s.logger.Error(err)
		return nil, err
	}

	if role != enums.Admin {
		return nil, errors.New("permission denied")
	}
	return s.repository.User.GetAllOwners(ctx, params)
}
