package app

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"ride-hail/config"
	"ride-hail/internal/adapters/http/handle"
	"ride-hail/internal/adapters/http/server"
	postgres2 "ride-hail/internal/adapters/postgres"
	rabbit2 "ride-hail/internal/adapters/rabbit"
	"ride-hail/internal/app/ride"
	"ride-hail/internal/core/domain/types"
	"ride-hail/internal/core/ports"
	"ride-hail/internal/core/service"
	"ride-hail/pkg/logger"
	postgres "ride-hail/pkg/potgres"
	"ride-hail/pkg/rabbit"
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

func New(cfg config.Config) (*App, error) {

	Svc := initService(cfg)

	return &App{
		authServ: authSvc,
	}, nil
}

func (app *App) Start() {
	err := app.authServ.Run()
	if err != nil {
		return
	}
}

func initService(cfg config.Config) (Service, error) {
	switch cfg.Mode {
	case types.ModeAdmin:
	case types.ModeDAL:
	case types.ModeRide:
		ride.New(cfg)
	default:
		return nil, fmt.Errorf("unknown mode: %s", cfg.Mode)
	}
	return nil, nil
}
