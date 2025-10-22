package config

import (
	"bufio"
	"errors"
	"os"
	"ride-hail/internal/core/domain/types"
	"ride-hail/pkg/potgres"
	"ride-hail/pkg/rabbit"
	"strconv"
	"strings"
)

type Config struct {
	Mode      string
	Database  postgres.Config
	RabbitMQ  rabbit.Config
	WebSocket struct {
		Port int
	}
	Services struct {
		RideService           int
		DriverLocationService int
		AdminService          int
	}
	JWT struct {
		Secret      string
		ExpireHours int
	}
}

func New(configPath, mode string) (*Config, error) {
	cfg, err := parseConfig(configPath)
	if err != nil {
		return nil, err
	}

	if !cfg.parseMode(mode) {
		return nil, errors.New("invalid Mode")
	}

	cfg.printConfig()
	return cfg, nil
}

func parseConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return &Config{}, err
	}
	defer file.Close()

	var cfg Config
	section := ""

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if !strings.Contains(line, ":") {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// --- поддержка ${ENV:-default} ---
		if strings.HasPrefix(value, "${") {
			inner := strings.TrimSuffix(strings.TrimPrefix(value, "${"), "}")
			innerParts := strings.SplitN(inner, ":-", 2)

			envName := innerParts[0]
			defVal := ""
			if len(innerParts) == 2 {
				defVal = innerParts[1]
			}

			if envVal, ok := os.LookupEnv(envName); ok {
				value = envVal
			} else {
				value = defVal
			}
		}

		switch key {
		case "postgres", "rabbitmq", "websocket", "services", "jwt":
			section = key

		default:
			switch section {
			case "postgres":
				switch key {
				case "host":
					cfg.Database.Host = value
				case "port":
					cfg.Database.Port = value
				case "user":
					cfg.Database.User = value
				case "password":
					cfg.Database.Password = value
				case "database":
					cfg.Database.Database = value
				}
			case "rabbitmq":
				switch key {
				case "host":
					cfg.RabbitMQ.Host = value
				case "port":
					cfg.RabbitMQ.Port, _ = strconv.Atoi(value)
				case "user":
					cfg.RabbitMQ.User = value
				case "password":
					cfg.RabbitMQ.Password = value
				}
			case "websocket":
				if key == "port" {
					cfg.WebSocket.Port, _ = strconv.Atoi(value)
				}
			case "services":
				switch key {
				case "ride_service":
					cfg.Services.RideService, _ = strconv.Atoi(value)
				case "driver_location_service":
					cfg.Services.DriverLocationService, _ = strconv.Atoi(value)
				case "admin_service":
					cfg.Services.AdminService, _ = strconv.Atoi(value)
				}
			case "jwt":
				switch key {
				case "secret":
					cfg.JWT.Secret = value
				case "expire_hours":
					cfg.JWT.ExpireHours, _ = strconv.Atoi(value)
				}
			}
		}
	}

	if cfg.Database.MaxOpenConns == 0 {
		cfg.Database.MaxOpenConns = 25
	}
	if cfg.Database.MaxIdleTime == "" {
		cfg.Database.MaxIdleTime = "15m"
	}

	return &cfg, scanner.Err()
}

func (cfg *Config) parseMode(mode string) bool {
	switch mode {
	case types.ModeAdmin, types.ModeRide, types.ModeDAL:
		cfg.Mode = mode
	default:
		return false
	}
	return true
}
