package dal

import (
	"context"

	"ride-hail/config"
)

type DriverService struct{}

func New(ctx context.Context, cfg config.Config) (*DriverService, error) {
	return &DriverService{}, nil
}

func (r *DriverService) Run() error {
	return nil
}

func (r *DriverService) Stop() {}
