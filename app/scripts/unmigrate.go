package main

import (
	"context"
	"github.com/alibekabdrakhman1/orynal/config"
	"github.com/alibekabdrakhman1/orynal/internal/model"
	"github.com/alibekabdrakhman1/orynal/pkg/gorm"
	"github.com/fatih/color"
	"log"
)

func main() {
	cfg, err := config.LoadConfig("./")
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	gormObject, err := gorm.Dial(ctx, cfg.DSN())
	if err != nil {
		return
	}

	log.Println(color.YellowString("Откат миграций..."))

	err = gormObject.Migrator().DropTable(model.User{})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(color.GreenString("Откат миграций завершен"))
}
