package ride

import (
	"context"
	"log/slog"
	"ride-hail/internal/adapters/http/handle"
	"ride-hail/internal/adapters/http/server"
	"ride-hail/internal/adapters/http/websocket"
	"ride-hail/internal/adapters/postgres"
	rabbit2 "ride-hail/internal/adapters/rabbit"
	"ride-hail/internal/core/ports"
	"ride-hail/internal/core/service"
	"ride-hail/pkg/logger"
	"ride-hail/pkg/rabbit"
	"ride-hail/pkg/txm"

	"ride-hail/config"
	pg "ride-hail/pkg/potgres"
)

type RideService struct {
	server server.Server
	svc    ports.RideService
}

func New(ctx context.Context, cfg config.Config) (*RideService, error) {
	log := logger.NewLogger(
		cfg.Mode, logger.Options{
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

	rb, err := rabbit.New(cfg.RabbitMQ)
	if err != nil {
		return nil, err
	}

	if err = rabbit2.InitRabbitTopology(rb); err != nil {
		return nil, err
	}

	rPub := rabbit.NewPublisher(rb.Conn)
	lCons := rabbit2.NewLocationConsumer(rb.Conn)
	dmCons := rabbit2.NewDriverResponseConsumer(rb.Conn)
	rSCons := rabbit2.NewRideStatusConsumer(rb.Conn)

	tmx := txm.NewTXManager(pg.Pool)

	wsm := websocket.NewPassengerWebSocketManager(log)
	wsh := websocket.NewPassengerWebSocketHandler(wsm, log)

	authServ := service.NewAuthService(cfg, uRepo, log)
	rideServ := service.NewRideService(log, tmx, rRepo, cRepo, wsm, rPub, lCons, dmCons, rSCons)

	authHandle := handle.New(cfg, authServ, log)
	rideHandle := handle.NewRideHandle(rideServ, wsh, log)

	serv, err := server.New(cfg, log, authHandle, rideHandle)
	if err != nil {
		return nil, err
	}

	return &RideService{
		server: serv,
		svc:    rideServ,
	}, nil
}

func (r *RideService) Run(ctx context.Context) {
	go r.svc.StartService(ctx)
	go r.server.Run()

}

func (r *RideService) Stop(ctx context.Context) error {
	return r.server.Stop(ctx)
}
