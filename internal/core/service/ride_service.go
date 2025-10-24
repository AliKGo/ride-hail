package service

import (
	"context"
	"ride-hail/internal/core/domain/action"
	"ride-hail/internal/core/domain/models"
	"ride-hail/internal/core/domain/types"
	"ride-hail/internal/core/ports"
	"ride-hail/internal/core/service/calculator"
	"ride-hail/pkg/logger"
	"ride-hail/pkg/txm"
)

type RideService struct {
	log  *logger.Logger
	repo Repository
	txm  txm.Manager
}

type Repository struct {
	ride ports.RideRepository
	cord ports.CoordinatesRepository
}

func NewRideService(log *logger.Logger, txm txm.Manager, rideRepo ports.RideRepository, cordRepo ports.CoordinatesRepository) *RideService {
	return &RideService{
		log: log,
		txm: txm,
		repo: Repository{
			ride: rideRepo,
			cord: cordRepo,
		},
	}
}

func (svc *RideService) CreateNewRide(ctx context.Context, r models.CreateRideRequest) (models.CreateRideResponse, error) {
	log := svc.log.Func("RideHandle.CreateNewRide")
	var rideID string

	dist := calculator.Distance(r.PickupLatitude, r.PickupLongitude, r.DestinationLatitude, r.DestinationLongitude)
	minute := calculator.Duration(dist)
	fareAmount, err := calculator.CalculateFare(r.RideType, dist, minute)
	if err != nil {
		log.Error(ctx, action.CreateRide, "error calculating fare amount", "error", err)
		return models.CreateRideResponse{}, err
	}

	fn := func(ctx context.Context) error {
		pickupID, err := svc.repo.cord.CreateNewCoordinate(ctx, models.Coordinate{
			EntityID:        logger.GetUserID(ctx),
			EntityType:      types.EntityRolePassenger,
			Address:         r.PickupAddress,
			Latitude:        r.PickupLatitude,
			Longitude:       r.PickupLongitude,
			FareAmount:      fareAmount,
			DurationMinutes: minute,
			DistanceKM:      dist,
			IsCurrent:       true,
		})
		if err != nil {
			return err
		}

		destinationID, err := svc.repo.cord.CreateNewCoordinate(ctx, models.Coordinate{
			EntityID:        logger.GetUserID(ctx),
			EntityType:      types.EntityRolePassenger,
			Address:         r.DestinationAddress,
			Latitude:        r.DestinationLatitude,
			Longitude:       r.DestinationLongitude,
			FareAmount:      fareAmount,
			DurationMinutes: minute,
			DistanceKM:      dist,
			IsCurrent:       true,
		})
		if err != nil {
			return err
		}

		rideID, err = svc.repo.ride.CreateNewRide(ctx, models.Ride{
			RideNumber:              "",
			PassengerID:             logger.GetUserID(ctx),
			VehicleType:             r.RideType,
			Status:                  types.RideStatusREQUESTED,
			EstimatedFare:           fareAmount,
			PickupCoordinateId:      pickupID,
			DestinationCoordinateId: destinationID,
		})
		if err != nil {
			return err
		}
		return nil
	}

	err = svc.txm.Do(ctx, fn)
	if err != nil {
		log.Error(ctx, action.CreateRide, "error in saved data", "error", err)
		return models.CreateRideResponse{}, err
	}

	return models.CreateRideResponse{
		RideID:                   rideID,
		RideNumber:               "",
		Status:                   types.RideStatusREQUESTED,
		EstimatedFare:            fareAmount,
		EstimatedDurationMinutes: minute,
		EstimatedDistanceKm:      dist,
	}, nil
}
