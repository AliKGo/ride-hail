package service

import (
	"context"
	"encoding/json"
	"fmt"
	"ride-hail/internal/core/domain/action"
	"ride-hail/internal/core/domain/models"
	"ride-hail/internal/core/domain/types"
	"ride-hail/internal/core/ports"
	"ride-hail/internal/core/service/calculator"
	"ride-hail/pkg/logger"
	"ride-hail/pkg/txm"
	"time"
)

type RideService struct {
	log       *logger.Logger
	repo      Repository
	txm       txm.Manager
	msgBroker MsgBroker
}

type MsgBroker struct {
	publisher ports.RidePublisher
}

type Repository struct {
	ride ports.RideRepository
	cord ports.CoordinatesRepository
}

func NewRideService(log *logger.Logger, txm txm.Manager, rideRepo ports.RideRepository, cordRepo ports.CoordinatesRepository, rPub ports.RidePublisher) *RideService {
	return &RideService{
		log: log,
		txm: txm,
		repo: Repository{
			ride: rideRepo,
			cord: cordRepo,
		},
		msgBroker: MsgBroker{
			publisher: rPub,
		},
	}
}

const exchangeName = "ride_topic"

var (
	queueRideRequests = "ride_requests"
)

func (svc *RideService) CreateNewRide(ctx context.Context, r models.CreateRideRequest) (models.CreateRideResponse, error) {
	log := svc.log.Func("RideService.CreateNewRide")

	dist := calculator.Distance(r.PickupLatitude, r.PickupLongitude, r.DestinationLatitude, r.DestinationLongitude)
	minute := calculator.Duration(dist)
	fareAmount, err := calculator.CalculateFare(r.RideType, dist, minute)
	if err != nil {
		log.Error(ctx, action.CreateRide, "error calculating fare amount", "error", err)
		return models.CreateRideResponse{}, err
	}

	newRide := models.Ride{
		PassengerID:   logger.GetUserID(ctx),
		VehicleType:   r.RideType,
		Status:        types.RideStatusREQUESTED,
		EstimatedFare: fareAmount,
	}

	fn := func(ctx context.Context) error {
		if newRide.PickupCoordinateId, err = svc.repo.cord.CreateNewCoordinate(ctx, models.Coordinate{
			EntityID:        logger.GetUserID(ctx),
			EntityType:      types.EntityRolePassenger,
			Address:         r.PickupAddress,
			Latitude:        r.PickupLatitude,
			Longitude:       r.PickupLongitude,
			FareAmount:      fareAmount,
			DurationMinutes: minute,
			DistanceKM:      dist,
			IsCurrent:       true,
		}); err != nil {
			log.Error(ctx, action.CreateRide, "error creating new coordinate", "error", err)
			return err
		}

		if newRide.DestinationCoordinateId, err = svc.repo.cord.CreateNewCoordinate(ctx, models.Coordinate{
			EntityID:        logger.GetUserID(ctx),
			EntityType:      types.EntityRolePassenger,
			Address:         r.DestinationAddress,
			Latitude:        r.DestinationLatitude,
			Longitude:       r.DestinationLongitude,
			FareAmount:      fareAmount,
			DurationMinutes: minute,
			DistanceKM:      dist,
			IsCurrent:       true,
		}); err != nil {
			log.Error(ctx, action.CreateRide, "error creating new coordinate", "error", err)
			return err
		}

		if rNumber, err := svc.repo.ride.GenerateRideNumber(ctx); err != nil {
			log.Error(ctx, action.CreateRide, "error generating ride number", "error", err)
			return err
		} else {
			newRide.RideNumber = fmt.Sprintf("RIDE_%s_%03d", time.Now().Format("20060102"), rNumber)
		}

		newRide.ID, err = svc.repo.ride.CreateNewRide(ctx, newRide)
		if err != nil {
			log.Error(ctx, action.CreateRide, "error creating new ride", "error", err)
			return err
		}

		if data, err := json.Marshal(newRide); err != nil {
			log.Error(ctx, action.CreateRide, "error marshalling new ride", "error", err)
			return err
		} else if err = svc.msgBroker.publisher.Publish(exchangeName, queueRideRequests, data); err != nil {
			log.Error(ctx, action.CreateRide, "error publishing ride", "error", err)
		}

		return nil
	}

	if err = svc.txm.Do(ctx, fn); err != nil {
		return models.CreateRideResponse{}, err
	}

	return models.CreateRideResponse{
		RideID:                   newRide.ID,
		RideNumber:               newRide.RideNumber,
		Status:                   types.RideStatusREQUESTED,
		EstimatedFare:            fareAmount,
		EstimatedDurationMinutes: minute,
		EstimatedDistanceKm:      dist,
	}, nil
}

func (svc *RideService) CloseRide(ctx context.Context, req models.CloseRideRequest) (models.CloseRideResponse, error) {
	return models.CloseRideResponse{}, nil
}
