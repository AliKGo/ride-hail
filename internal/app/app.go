package app

import (
	"context"
	"fmt"
	"ride-hail/config"
	"ride-hail/internal/app/ride"
	"ride-hail/internal/core/domain/types"
)

type Service interface {
	Run()
	Stop(ctx context.Context) error
}

type App struct {
	svc Service
}

func New(ctx context.Context, cfg config.Config) (*App, error) {
	svc, err := initService(ctx, cfg)
	if err != nil {
		return &App{}, err
	}
	return &App{
		svc: svc,
	}, nil
}

func (app *App) Start() {
	app.svc.Run()
}

func initService(ctx context.Context, cfg config.Config) (Service, error) {
	switch cfg.Mode {
	case types.ModeAdmin:
	case types.ModeDAL:
	case types.ModeRide:
		return ride.New(ctx, cfg)
	default:
		return nil, fmt.Errorf("unknown mode: %s", cfg.Mode)
	}
	return nil, nil
}
