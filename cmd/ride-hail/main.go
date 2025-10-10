package ride_hail

import (
	"flag"
	"log/slog"
	"ride-hail/config"
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

}
