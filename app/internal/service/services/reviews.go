package services

import (
	"context"
	"errors"
	"github.com/alibekabdrakhman1/orynal/config"
	"github.com/alibekabdrakhman1/orynal/internal/model"
	"github.com/alibekabdrakhman1/orynal/internal/repository"
	"github.com/alibekabdrakhman1/orynal/internal/service/infrastructure"
	"github.com/alibekabdrakhman1/orynal/pkg/utils"
	"go.uber.org/zap"
	"time"
)

func NewReviewsService(repository *repository.Manager, config *config.Config, logger *zap.SugaredLogger) *ReviewsService {
	return &ReviewsService{repository: repository, config: config, logger: logger, FormatParams: infrastructure.NewFormatParams()}
}

type ReviewsService struct {
	repository *repository.Manager
	config     *config.Config
	logger     *zap.SugaredLogger
	FormatParams
}

func (s *ReviewsService) GetReviews(ctx context.Context, restaurantID uint, params *model.Params) (*model.ListResponse, error) {
	resp, err := s.repository.Reviews.GetReviews(ctx, restaurantID, params)
	if err != nil {
		return nil, err
	}

	resp.TotalPages = int(resp.TotalItems) / resp.ItemsPerPage
	if int(resp.TotalPages)%resp.ItemsPerPage != 0 {
		resp.TotalPages++
	}

	return resp, nil
}

func (s *ReviewsService) CreateReview(ctx context.Context, review *model.RestaurantReview) (*model.RestaurantReview, error) {
	id, err := utils.GetIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	review.UserID = id
	review.Date = time.Now()

	return s.repository.Reviews.CreateReview(ctx, review)
}

func (s *ReviewsService) DeleteReview(ctx context.Context, id uint) error {
	id, err := utils.GetIDFromContext(ctx)
	if err != nil {
		return err
	}

	review, err := s.repository.Reviews.GetReview(ctx, id)
	if err != nil {
		return err
	}

	if review.UserID != id {
		return errors.New("not author")
	}

	return s.repository.Reviews.DeleteReview(ctx, id)
}
