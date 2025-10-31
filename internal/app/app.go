package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"ride-hail/pkg/logger"
	"syscall"
	"time"

	"ride-hail/config"
	"ride-hail/internal/app/ride"
	"ride-hail/internal/core/domain/action"
	"ride-hail/internal/core/domain/types"
)

type Service interface {
	Run(ctx context.Context)
	Stop(ctx context.Context) error
}

type App struct {
	svc Service
	log *logger.Logger
}

func New(ctx context.Context, cfg config.Config) (*App, error) {
	log := logger.NewLogger(
		cfg.Mode, logger.Options{
			Pretty: true,
			Level:  slog.LevelDebug,
		},
	)

	svc, err := initService(ctx, log, cfg)
	if err != nil {
		log.Func("New").Error(ctx, action.StartApplication, "failed to initialize service", "error", err)
		return &App{}, err
	}

	log.Func("New").Debug(ctx, action.StartApplication, "service initialized successfully")

	return &App{
		svc: svc,
		log: log,
	}, nil
}

func (app *App) Start(ctx context.Context) {
	log := app.log.Func("Start")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	log.Info(ctx, action.StartApplication, "starting service")
	go app.svc.Run(ctx)

	sig := <-sigChan
	log.Warn(ctx, action.StopApplication, "received shutdown signal", "signal", sig.String())

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Info(ctx, action.StopApplication, "initiating graceful shutdown with 30 second timeout")

	if err := app.svc.Stop(shutdownCtx); err != nil {
		log.Error(ctx, action.StopApplication, "error during service shutdown", "error", err)
		return
	}

	log.Info(ctx, action.StopApplication, "application stopped successfully")
}

func initService(ctx context.Context, log *logger.Logger, cfg config.Config) (Service, error) {
	funcLog := log.Func("initService")

	funcLog.Debug(ctx, action.StartApplication, "initializing service", "mode", cfg.Mode)

	switch cfg.Mode {
	case types.ModeAdmin:
		funcLog.Debug(ctx, action.StartApplication, "admin service mode detected")
	case types.ModeDAL:
		funcLog.Debug(ctx, action.StartApplication, "driver location service mode detected")
	case types.ModeRide:
		funcLog.Debug(ctx, action.StartApplication, "ride service mode detected")
		return ride.New(ctx, log, cfg)
	default:
		err := fmt.Errorf("unknown mode: %s", cfg.Mode)
		funcLog.Error(ctx, action.StartApplication, "unsupported service mode", "mode", cfg.Mode, "error", err)
		return nil, err
	}
	return nil, nil
}
