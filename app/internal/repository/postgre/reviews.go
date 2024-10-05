package postgre

import (
	"context"
	"github.com/alibekabdrakhman1/orynal/internal/model"
	"gorm.io/gorm"
)

func NewReviewsRepository(db *gorm.DB) *ReviewsRepository {
	return &ReviewsRepository{
		DB: db,
	}
}

type ReviewsRepository struct {
	DB *gorm.DB
}

func (r *ReviewsRepository) GetReviews(ctx context.Context, restaurantID uint, params *model.Params) (*model.ListResponse, error) {
	var reviews []*model.RestaurantReview
	query := r.DB.Table("restaurant_reviews").
		Where("restaurant_id = ?", restaurantID).
		Preload("User").
		Limit(params.Limit).
		Offset(params.Offset)

	if err := query.Find(&reviews).Error; err != nil {
		return nil, err
	}

	for i := 0; i < len(reviews); i++ {
		var user model.UserResponse
		if err := r.DB.Table("users").Where("id = ?", reviews[i].UserID).First(&user).Error; err != nil {
			return nil, err
		}

		reviews[i].User = user
	}

	var totalItems int64
	if err := r.DB.Model(&model.RestaurantReview{}).Where("restaurant_id = ?", restaurantID).Count(&totalItems).Error; err != nil {
		return nil, err
	}

	return &model.ListResponse{
		Items:        reviews,
		ItemsPerPage: params.Limit,
		PageIndex:    params.PageIndex,
		TotalItems:   int(totalItems),
	}, nil
}

func (r *ReviewsRepository) CreateReview(ctx context.Context, review *model.RestaurantReview) (*model.RestaurantReview, error) {
	if err := r.DB.WithContext(ctx).Table("restaurant_reviews").Create(review).Error; err != nil {
		return nil, err
	}

	var user model.UserResponse
	if err := r.DB.Table("users").Where("id = ?", review.UserID).First(&user).Error; err != nil {
		return nil, err
	}

	review.User = user

	return review, nil
}

func (r *ReviewsRepository) DeleteReview(ctx context.Context, id uint) error {
	if err := r.DB.WithContext(ctx).Table("restaurant_reviews").Where("id = ?", id).Delete(&model.RestaurantReview{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *ReviewsRepository) GetReview(ctx context.Context, id uint) (*model.RestaurantReview, error) {
	var review model.RestaurantReview
	if err := r.DB.WithContext(ctx).Table("restaurant_reviews").First(&review, id).Error; err != nil {
		return nil, err
	}

	var user model.UserResponse
	if err := r.DB.Table("users").Where("id = ?", review.UserID).First(&user).Error; err != nil {
		return nil, err
	}

	review.User = user
	return &review, nil
}
