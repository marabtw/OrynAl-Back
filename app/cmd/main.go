package main

import (
	"github.com/alibekabdrakhman1/orynal/config"
	"github.com/alibekabdrakhman1/orynal/internal/app"
	"go.uber.org/zap"
)

//	@title			Reviews
//	@version		1.0.0
//	@description	Reviews Service

//	@host		localhost:5000
//	@BasePath	/api

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	l := logger.Sugar()
	l = l.With(zap.String("app", "orynal-app"))

	cfg, err := config.LoadConfig("./")
	if err != nil {
		l.Error(err)
		l.Fatalf("failed to load configs err: %v", err)
	}

	app := app.New(l, &cfg)
	app.Run()
}
