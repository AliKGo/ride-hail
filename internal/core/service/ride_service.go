package service

import (
	"context"
	"encoding/json"
	"errors"
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
	wsm       ports.PassengerWSManager
	txm       txm.Manager
	msgBroker MsgBroker
}

type MsgBroker struct {
	producer            ports.RideProducer
	consumerLocation    ports.LocationSubscriber
	consumerDriverMatch ports.DriverMatchSubscriber
	consumerRideStatus  ports.RideStatusSubscriber
}

type Repository struct {
	ride ports.RideRepository
	cord ports.CoordinatesRepository
}

func NewRideService(log *logger.Logger, txm txm.Manager, rideRepo ports.RideRepository, cordRepo ports.CoordinatesRepository, wsm ports.PassengerWSManager, rPub ports.RideProducer, consumerLocation ports.LocationSubscriber, consumerDriverMatch ports.DriverMatchSubscriber, consumerRideStatus ports.RideStatusSubscriber) *RideService {
	return &RideService{
		log: log,
		txm: txm,
		wsm: wsm,
		repo: Repository{
			ride: rideRepo,
			cord: cordRepo,
		},
		msgBroker: MsgBroker{
			producer:            rPub,
			consumerLocation:    consumerLocation,
			consumerRideStatus:  consumerRideStatus,
			consumerDriverMatch: consumerDriverMatch,
		},
	}
}

const exchangeName = "ride_topic"

func (svc *RideService) StartService(ctx context.Context) {
	log := svc.log.Func("RideService.StartService")

	go svc.runWithRetry(ctx, svc.driverMatch, "driverMatch")
	go svc.runWithRetry(ctx, svc.driverLocation, "driverLocation")
	go svc.runWithRetry(ctx, svc.rideStatus, "rideStatus")

	log.Debug(ctx, action.ServiceRide, "RideService started")
	<-ctx.Done()
	log.Debug(ctx, action.ServiceRide, "RideService stopping")
}

func (svc *RideService) runWithRetry(ctx context.Context, fn func(ctx context.Context) error, name string) {
	log := svc.log.Func("RideService." + name)
	backoff := time.Second

	for {
		select {
		case <-ctx.Done():
			log.Debug(ctx, action.ServiceRide, "stopping service")
			return
		default:
			err := fn(ctx)
			if err != nil {
				log.Error(ctx, action.ServiceRide, fmt.Sprintf("%s failed, retrying", name), "error", err)
				time.Sleep(backoff)
				if backoff < 30*time.Second {
					backoff *= 2
				}
			} else {
				backoff = time.Second
			}
		}
	}
}

func (svc *RideService) rideStatus(ctx context.Context) error {
	log := svc.log.Func("RideService.rideStatus")

	ch, err := svc.msgBroker.consumerRideStatus.Subscribe(ctx)
	if err != nil {
		log.Error(ctx, action.ServiceRide, "failed to subscribe ride status", "error", err)
		return err
	}

	for {
		select {
		case <-ctx.Done():
			log.Debug(ctx, action.ServiceRide, "ride status stopped")
			return nil
		case msg, ok := <-ch:
			if !ok {
				log.Debug(ctx, action.ServiceRide, "ride status channel closed, exiting")
				return fmt.Errorf("ride status channel closed")
			}
			svc.parsingRideStatus(ctx, msg)
		}
	}
}

func (svc *RideService) parsingRideStatus(ctx context.Context, msg models.RideStatusEvent) {
	log := svc.log.Func("RideService.parsingRideStatus")
	ctxNew := logger.WithRequestID(ctx, msg.CorrelationID)

	if err := svc.repo.ride.UpdateRide(ctx, msg.RideID, msg.Status, "", &msg.Timestamp); err != nil {
		log.Error(ctxNew, action.ServiceRide, "failed to update ride in database", "error", err)
		return
	}
	if err := svc.wsm.SendRideStatusUpdate(ctxNew, msg.RideID, models.RideStatusUpdate{
		RideID:        msg.RideID,
		Status:        msg.Status,
		Timestamp:     msg.Timestamp,
		DriverID:      msg.DriverID,
		CorrelationID: msg.CorrelationID,
	}); err != nil {
		log.Error(ctxNew, action.ServiceRide, "failed to send ride status update", "error", err)
		return
	}
}

func (svc *RideService) driverMatch(ctx context.Context) error {
	log := svc.log.Func("RideService.driverMatchResp")

	ch, err := svc.msgBroker.consumerDriverMatch.Subscribe(ctx)
	if err != nil {
		log.Error(ctx, action.ServiceRide, "Failed to subscribe to msg broker")
		return fmt.Errorf("failed to subscribe to DriverMatch %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			log.Debug(ctx, action.ServiceRide, "driverMatch stopped")
			return nil
		case msg, ok := <-ch:
			if !ok {
				log.Debug(ctx, action.ServiceRide, "driverMatch channel closed")
				return fmt.Errorf("driverMatch channel closed")
			}
			go svc.parsingDriverMatch(ctx, msg)
		}
	}
}

func (svc *RideService) parsingDriverMatch(ctx context.Context, driverResp models.DriverResponseEvent) {
	log := svc.log.Func("RideService.parsingDriverStatus")
	ctxNew := logger.WithRequestID(ctx, driverResp.CorrelationID)
	now := time.Now()

	ride, err := svc.repo.ride.GetRide(ctx, driverResp.RideID)
	if err != nil {
		log.Error(ctxNew, action.ServiceRide, "failed to get ride")
		return
	}

	if err = svc.repo.ride.UpdateMatchedRide(ctx, driverResp.RideID, driverResp.DriverID, now); err != nil {
		log.Error(ctxNew, action.ServiceRide, "failed to update matched ride")
		return
	}

	if err = svc.wsm.SendRideStatusUpdate(ctx, ride.ID, models.RideStatusUpdate{
		RideID:        driverResp.RideID,
		Status:        types.RideStatusMATCHED,
		Timestamp:     now,
		DriverID:      driverResp.DriverID,
		CorrelationID: driverResp.CorrelationID,
	}); err != nil {
		log.Error(ctxNew, action.ServiceRide, "failed to send ride-status update")
		return
	}
}

func (svc *RideService) driverLocation(ctx context.Context) error {
	log := svc.log.Func("RideService.StartService")

	ch, err := svc.msgBroker.consumerLocation.Subscribe(ctx)
	if err != nil {
		log.Error(ctx, action.ServiceRide, "error in getting consumer location subscription", "error", err)
		return err
	}

	for {
		select {
		case <-ctx.Done():
			log.Debug(ctx, action.ServiceRide, "driverLocation stopped")
			return nil
		case msg, ok := <-ch:
			if !ok {
				log.Debug(ctx, action.ServiceRide, "driverLocation channel closed")
				return fmt.Errorf("driverLocation channel closed")
			}
			go svc.processingMsg(ctx, msg)
		}
	}

}

func (svc *RideService) processingMsg(ctx context.Context, msg models.DriverLocationUpdate) {
	log := svc.log.Func("RideService.processingMsg")

	ride, err := svc.repo.ride.GetRide(ctx, msg.RideID)
	if err != nil {
		log.Error(ctx, action.ServiceRide, "failed to get ride", "error", err)
		return
	}

	if err = svc.wsm.SendDriverLocationUpdate(ctx, ride.PassengerID, msg); err != nil {
		log.Error(ctx, action.ServiceRide, "failed to send driver location update", "error", err)
		return
	}
}

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

		if data, err := json.Marshal(models.RideRequestRideType{
			RideID:              newRide.ID,
			RideNumber:          newRide.RideNumber,
			PickupLocation:      models.Location{Lat: r.PickupLatitude, Lng: r.PickupLongitude, Address: r.PickupAddress},
			DestinationLocation: models.Location{Lat: r.DestinationLatitude, Lng: r.DestinationLongitude, Address: r.DestinationAddress},
			RideType:            r.RideType,
			EstimatedFare:       fareAmount,
			MaxDistanceKm:       dist,
			TimeoutSeconds:      30,
			CorrelationID:       logger.GetRequestID(ctx),
		}); err != nil {
			log.Error(ctx, action.CreateRide, "error marshalling new ride", "error", err)
			return err
		} else {
			routingKey := fmt.Sprintf("ride.request.%s", r.RideType)
			if err = svc.msgBroker.producer.Producer(exchangeName, routingKey, data); err != nil {
				log.Error(ctx, action.CreateRide, "error publishing ride", "error", err)
				return err
			}
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
	log := svc.log.Func("RideService.CloseRide")

	ride, err := svc.repo.ride.GetRide(ctx, req.RideID)
	if err != nil {
		log.Error(ctx, action.CloseRide, "error retrieving ride", "error", err)
		return models.CloseRideResponse{}, err
	}

	if ride.Status != types.RideStatusREQUESTED {
		log.Error(ctx, action.CloseRide, "ride.status is not requested", "error", err)
		return models.CloseRideResponse{}, errors.New("ride.status is not requested")
	}

	fn := func(ctx context.Context) error {
		now := time.Now()

		if err = svc.repo.ride.UpdateRide(ctx, ride.ID, types.RideStatusCANCELLED, req.Reason, &now); err != nil {
			log.Error(ctx, action.CloseRide, "error updating ride", "error", err)
			return err
		}

		data, err := json.Marshal(models.RideStatusUpdate{
			RideID:        ride.ID,
			Status:        types.RideStatusCANCELLED,
			Timestamp:     now,
			DriverID:      ride.DriverID,
			CorrelationID: logger.GetRequestID(ctx),
		})
		if err != nil {
			log.Error(ctx, action.CloseRide, "error marshalling ride", "error", err)
			return err
		}

		routingKey := fmt.Sprintf("ride.status.%s", types.RideStatusCANCELLED)
		if err = svc.msgBroker.producer.Producer(exchangeName, routingKey, data); err != nil {
			log.Error(ctx, action.CloseRide, "error publishing ride status", "error", err)
		}

		return nil
	}

	if err = svc.txm.Do(ctx, fn); err != nil {
		return models.CloseRideResponse{}, err
	}

	return models.CloseRideResponse{}, nil
}
