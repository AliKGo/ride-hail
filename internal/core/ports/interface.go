package ports

import (
	"context"
	"time"

	"ride-hail/internal/core/domain/models"
)

type AuthServices interface {
	Run() error
}

// auth ports

type AuthService interface {
	CreateNewUser(ctx context.Context, user models.User) error
	Login(ctx context.Context, user models.User) (string, string, error)
}

type UserRepository interface {
	CreateNewUser(ctx context.Context, user models.User) error
	GetGyUserEmail(ctx context.Context, email string) (models.User, error)
}

// ride ports

type PassengerWSManager interface {
	SendRideStatusUpdate(ctx context.Context, passengerID string, update models.RideStatusUpdate) error
	SendDriverLocationUpdate(ctx context.Context, passengerID string, update models.DriverLocationUpdate) error
}

type RideService interface {
	StartService(ctx context.Context)
	CreateNewRide(ctx context.Context, r models.CreateRideRequest) (models.CreateRideResponse, error)
	CloseRide(ctx context.Context, req models.CloseRideRequest) (models.CloseRideResponse, error)
}

type RideProducer interface {
	Producer(exName, queue string, message []byte) error
}

type LocationSubscriber interface {
	Subscribe(ctx context.Context) (<-chan models.DriverLocationUpdate, error)
	Start(ctx context.Context) error
}

type DriverMatchSubscriber interface {
	Subscribe(ctx context.Context) (<-chan models.DriverResponseEvent, error)
	Start(ctx context.Context) error
}
type RideStatusSubscriber interface {
	Subscribe(ctx context.Context) (<-chan models.RideStatusEvent, error)
	Start(ctx context.Context) error
}

type RideRepository interface {
	CreateNewRide(ctx context.Context, ride models.Ride) (string, error)
	GetRide(ctx context.Context, id string) (models.Ride, error)
	UpdateRide(ctx context.Context, rideID string, newStatus string, reason string, t *time.Time) error
	UpdateMatchedRide(ctx context.Context, rideID, driverID string, matchedAt time.Time) error
	GenerateRideNumber(ctx context.Context) (int, error)
}

type CoordinatesRepository interface {
	CreateNewCoordinate(ctx context.Context, c models.Coordinate) (string, error)
	GetCoordinate(ctx context.Context, id string) (models.Coordinate, error)
}
