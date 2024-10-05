package app

import (
	"context"
	"github.com/alibekabdrakhman1/orynal/config"
	"github.com/alibekabdrakhman1/orynal/internal/controller"
	http "github.com/alibekabdrakhman1/orynal/internal/controller/http/handler"
	"github.com/alibekabdrakhman1/orynal/internal/controller/http/middleware"
	"github.com/alibekabdrakhman1/orynal/internal/repository"
	"github.com/alibekabdrakhman1/orynal/internal/service"
	"github.com/alibekabdrakhman1/orynal/pkg/gorm"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
)

type App struct {
	logger *zap.SugaredLogger
	config *config.Config
}

func New(logger *zap.SugaredLogger, cfg *config.Config) *App {
	return &App{
		config: cfg,
		logger: logger,
	}
}

func (a *App) Run() error {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()
	gracefullyShutdown(cancel)

	db, err := gorm.Dial(ctx, a.config.DSN())
	if err != nil {
		log.Fatalf("cannot сonnect to DB '%s:%d': %v", a.config.Database.Host, a.config.Database.Port, err)
	}

	repo := repository.NewManager(db)

	srv := service.NewManager(repo, a.config, a.logger)

	endPointHandler := http.NewManager(srv, a.logger)

	jwt := middleware.NewJWTAuth([]byte(a.config.Auth.JwtSecretKey), srv.Auth, a.logger)

	HTTPServer := controller.NewServer(a.config, endPointHandler, jwt)
	return HTTPServer.StartHTTPServer(ctx)
}

func gracefullyShutdown(c context.CancelFunc) {
	osC := make(chan os.Signal, 1)
	signal.Notify(osC, os.Interrupt)
	go func() {
		log.Print(<-osC)
		c()
	}()
}
