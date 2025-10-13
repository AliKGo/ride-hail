package app

import (
	"context"
	"ride-hail/config"
	"ride-hail/internal/adapters/http/handle"
	"ride-hail/internal/adapters/http/server"
	postgres2 "ride-hail/internal/adapters/postgres"
	"ride-hail/internal/core/service"
	"ride-hail/pkg/logger"
	postgres "ride-hail/pkg/potgres"
)

type Service interface {
	Run() error
	Stop()
}

type App struct {
	svc Service
	api *server.API
	log *logger.Logger
}

func New(ctx context.Context, cfg config.Config) (*App, error) {
	db, err := postgres.New(ctx, cfg.Database)
	if err != nil {
		return nil, err
	}

	log := logger.New(cfg.Mode)

	repo := postgres2.NewRepo(db.Pool, log)
	authSvc := service.NewAuthService(cfg, repo, log)

	handle := handle.New(cfg, authSvc, log)

	api := server.New(handle, cfg)

	return &App{
		api: api,
		log: log,
	}, nil
}

func (app *App) Start() {
	err := app.api.Run()
	if err != nil {
		return
	}
}
