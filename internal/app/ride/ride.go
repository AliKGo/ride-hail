package ride

import (
	"context"

	"ride-hail/config"
	// "ride-hail/internal/adapters/http/handle"
	// "ride-hail/internal/adapters/postgres"
	// "ride-hail/internal/core/service"
	// "ride-hail/pkg/logger"
	// pg "ride-hail/pkg/potgres"
	// "ride-hail/pkg/txm"
)

type RideService struct{}

func New(ctx context.Context, cfg config.Config) (*RideService, error) {
	// log := logger.New(cfg.Mode, false)
	// pg, err := pg.New(ctx, cfg.Database)
	// if err != nil {
	// 	return nil, err
	// }

	// uRepo := postgres.NewRepo(pg.Pool)
	// cRepo := postgres.NewCordRepository(pg.Pool)
	// rRepo := postgres.NewRideRepository(pg.Pool)

	// tmx := txm.NewTXManager(pg.Pool)

	// authServ := service.NewAuthService(cfg, uRepo, log)
	// rideServ := service.NewRideService(log, tmx, rRepo, cRepo)

	// authHandle := handle.New(cfg, authServ, log)
	// rideHandle := handle.NewRideHandle(rideServ, log)

	return &RideService{}, nil
}

func (r *RideService) Run() error {
	return nil
}

func (r *RideService) Stop() {}
