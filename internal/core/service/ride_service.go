package service

import (
	"ride-hail/internal/core/ports"
	"ride-hail/pkg/logger"
)

type RideService struct {
	rideRepo ports.RideRepository
	CordRepo ports.CoordinatesRepository
	log      *logger.Logger
}

func NewRideService(log *logger.Logger) *RideService {
	return &RideService{
		log: log,
	}
}

//func (svc *RideService) CreateNewRide(ctx context.Context, rideRequest models.CreateRideRequest) (models.CreateRideResponse, error) {
//	//log := svc
//	//return result, nil
//}
//
//func (svc *RideService) CloseRide(ctx context.Context, req models.CloseRideRequest) (models.CloseRideResponse, error) {
//
//}
