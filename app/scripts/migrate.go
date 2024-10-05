package main

import (
	"context"
	"github.com/alibekabdrakhman1/orynal/config"
	"github.com/alibekabdrakhman1/orynal/internal/model"
	"github.com/alibekabdrakhman1/orynal/pkg/gorm"
	"github.com/fatih/color"
	"log"
)

func init() {

}

var TablesList = []interface{}{
	model.User{},
	model.UserToken{},
	model.Restaurant{},
	model.RestaurantPhoto{},
	model.Order{},
	model.OrderFood{},
	model.Food{},
	model.Table{},
}

func main() {
	cfg, err := config.LoadConfig("./")
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	log.Println(color.YellowString("Начала миграции..."))

	gormObject, err := gorm.Dial(ctx, cfg.DSN())
	if err != nil {
		return
	}

	err = gormObject.AutoMigrate(TablesList...)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(color.GreenString("Миграция завершена"))
}
