package controller

import (
	"context"
	"errors"
	"fmt"
	"github.com/alibekabdrakhman1/orynal/config"
	http "github.com/alibekabdrakhman1/orynal/internal/controller/http/handler"
	"github.com/alibekabdrakhman1/orynal/internal/controller/http/middleware"
	"github.com/labstack/echo/v4"
	middleware2 "github.com/labstack/echo/v4/middleware"
	"log"
	http2 "net/http"

	"time"
)

type Server struct {
	cfg     *config.Config
	handler *http.Manager
	App     *echo.Echo
	jwt     *middleware.JWTAuth
}

func NewServer(cfg *config.Config, handler *http.Manager, jwt *middleware.JWTAuth) *Server {
	return &Server{
		cfg:     cfg,
		handler: handler,
		jwt:     jwt,
	}
}

func (s *Server) StartHTTPServer(ctx context.Context) error {
	s.App = s.BuildEngine()
	s.SetupRoutes()
	go func() {
		if err := s.App.Start(fmt.Sprintf(":%v", s.cfg.HttpServer.Port)); err != nil && !errors.Is(err, http2.ErrServerClosed) {
			log.Fatalf("listen:%v\n", err)
		}
	}()
	<-ctx.Done()

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()
	if err := s.App.Shutdown(ctxShutDown); err != nil {
		log.Fatalf("controller Shutdown Failed:%v", err)
	}
	log.Print("controller exited properly")
	return nil
}

func (s *Server) BuildEngine() *echo.Echo {
	e := echo.New()
	e.Use(middleware2.CORSWithConfig(middleware2.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"*"},
	}))

	return e
}
