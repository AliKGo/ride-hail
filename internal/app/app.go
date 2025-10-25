package app

import (
	"context"
	"fmt"

	"ride-hail/config"
	"ride-hail/internal/app/ride"
	"ride-hail/internal/core/domain/types"
	"ride-hail/internal/core/ports"
	"ride-hail/pkg/logger"
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
	Svc, err := initService(ctx, cfg)
	if err != nil {
		return &App{}, nil
	}
	return &App{
		authServ: Svc,
	}, nil
}

func (app *App) Start() {
	err := app.authServ.Run()
	if err != nil {
		return
	}
}

func initService(ctx context.Context, cfg config.Config) (Service, error) {
	switch cfg.Mode {
	case types.ModeAdmin:
	case types.ModeDAL:
	case types.ModeRide:
		ride.New(ctx, cfg)
	default:
		return nil, fmt.Errorf("unknown mode: %s", cfg.Mode)
	}
	return nil, nil
}
