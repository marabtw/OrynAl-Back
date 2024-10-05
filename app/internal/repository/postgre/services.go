package postgre

import (
	"context"
	"github.com/alibekabdrakhman1/orynal/internal/model"
	"gorm.io/gorm"
)

func NewServicesRepository(db *gorm.DB) *ServicesRepository {
	return &ServicesRepository{
		DB: db,
	}
}

type ServicesRepository struct {
	DB *gorm.DB
}

func (r *ServicesRepository) CreateService(ctx context.Context, service *model.Service) ([]model.Service, error) {
	if err := r.DB.WithContext(ctx).Table("services").Create(service).Error; err != nil {
		return nil, err
	}

	services, err := r.GetServices(ctx)
	if err != nil {
		return nil, err
	}

	return services, nil
}

func (r *ServicesRepository) DeleteService(ctx context.Context, id uint) error {
	if err := r.DB.WithContext(ctx).Table("services").Where("id = ?", id).Delete(&model.Service{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *ServicesRepository) GetServices(ctx context.Context) ([]model.Service, error) {
	var services []model.Service
	if err := r.DB.WithContext(ctx).Table("services").Find(&services).Error; err != nil {
		return nil, err
	}
	return services, nil
}

func (r *ServicesRepository) UpdateService(ctx context.Context, service *model.Service) ([]model.Service, error) {
	if err := r.DB.WithContext(ctx).Table("services").Model(&model.Service{}).Where("id = ?", service.ID).Updates(service).Error; err != nil {
		return nil, err
	}

	services, err := r.GetServices(ctx)
	if err != nil {
		return nil, err
	}

	return services, nil
}
