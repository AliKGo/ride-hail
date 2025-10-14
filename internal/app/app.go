package app

import (
	"context"
	"ride-hail/config"
	"ride-hail/internal/adapters/http/handle"
	"ride-hail/internal/adapters/http/server"
	postgres2 "ride-hail/internal/adapters/postgres"
	"ride-hail/internal/core/ports"
	"ride-hail/internal/core/service"
	"ride-hail/pkg/logger"
	postgres "ride-hail/pkg/potgres"
)

type Service interface {
	Run() error
	Stop()
}

type App struct {
	svc      Service
	authServ ports.AuthServices
	log      *logger.Logger
}

func New(ctx context.Context, cfg config.Config) (*App, error) {
	db, err := postgres.New(ctx, cfg.Database)
	if err != nil {
		return nil, err
	}
	log := logger.New(cfg.Mode, false)

	authSvc := initAuth(db, cfg, log)

	return &App{
		authServ: authSvc,
		log:      log,
	}, nil
}

func (app *App) Start() {
	err := app.authServ.Run()
	if err != nil {
		return
	}
}

func initAuth(db *postgres.Postgres, cfg config.Config, log *logger.Logger) ports.AuthServices {
	repo := postgres2.NewRepo(db.Pool, log)
	authSvc := service.NewAuthService(cfg, repo, log)

	h := handle.New(cfg, authSvc, log)

	return server.New(h, cfg)
}
