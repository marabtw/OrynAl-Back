package gorm

import (
	"context"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Dial(ctx context.Context, url string) (*gorm.DB, error) {
	_ = ctx
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
