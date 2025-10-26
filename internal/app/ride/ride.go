package ride

import (
	"context"
	"log/slog"
	"ride-hail/internal/adapters/http/handle"
	"ride-hail/internal/adapters/http/server"
	"ride-hail/internal/adapters/postgres"
	"ride-hail/internal/core/service"
	"ride-hail/pkg/logger"
	"ride-hail/pkg/txm"

	"ride-hail/config"
	pg "ride-hail/pkg/potgres"
)

type RideService struct {
	server server.Server
}

func New(ctx context.Context, cfg config.Config) (*RideService, error) {
	log := logger.NewLogger(
		cfg.Mode, logger.LoggerOptions{
			Pretty: true,
			Level:  slog.LevelDebug,
		},
	)
	pg, err := pg.New(ctx, cfg.Database)
	if err != nil {
		return nil, err
	}

	uRepo := postgres.NewRepo(pg.Pool)
	cRepo := postgres.NewCordRepository(pg.Pool)
	rRepo := postgres.NewRideRepository(pg.Pool)

	tmx := txm.NewTXManager(pg.Pool)

	authServ := service.NewAuthService(cfg, uRepo, log)
	rideServ := service.NewRideService(log, tmx, rRepo, cRepo)

	authHandle := handle.New(cfg, authServ, log)
	rideHandle := handle.NewRideHandle(rideServ, log)

	serv, err := server.New(cfg, log, authHandle, rideHandle)
	if err != nil {
		return nil, err
	}

	return &RideService{
		server: serv,
	}, nil
}

func (r *RideService) Run() {
	r.server.Run()
}

func (r *RideService) Stop(ctx context.Context) error {
	return r.server.Stop(ctx)
}
