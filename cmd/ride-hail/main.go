package ride_hail

import (
	"context"
	"flag"
	"log/slog"
	"ride-hail/config"
	"ride-hail/internal/app"
)

var (
	modeFlag       = flag.String("mode", "-", "mode service")
	modeConfigPAth = flag.String("config-path", "./config.yaml", "path to config file")
)

func Run() {
	flag.Parse()
	cfg, err := config.New(*modeConfigPAth, *modeFlag)
	if err != nil {
		slog.Error("error in parsing config", err)
		return
	}

	ctx := context.Background()
	application, err := app.New(ctx, *cfg)
	if err != nil {
		slog.Error("error in creating app", err)
		return
	}
	application.Start(ctx)

}
