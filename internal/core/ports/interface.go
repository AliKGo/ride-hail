package ports

import (
	"context"
	"net/http"
	"ride-hail/internal/core/domain/models"
)

type AuthServices interface {
	Run() error
}

// auth ports

type AuthHandler interface {
	Login(w http.ResponseWriter, r *http.Request)
	Registration(w http.ResponseWriter, r *http.Request)
}

type AuthService interface {
	CreateNewUser(ctx context.Context, user models.User) error
	Login(ctx context.Context, user models.User) (string, error)
}

type UserRepository interface {
	CreateNewUser(ctx context.Context, user models.User) error
	GetGyUserEmail(ctx context.Context, email string) (models.User, error)
}

//ride ports

type RideHandler interface {
	CreateNewRide(ctx context.Context, ride models.Ride) error
}

type RideService interface {
	CreateNewRide(ctx context.Context, rideRequest models.CreateRideRequest) (models.CreateRideResponse, error)
	CloseRide(ctx context.Context, req models.CloseRideRequest) (models.CloseRideResponse, error)
}

type RideRepository interface {
	CreateNewRide(ctx context.Context, ride models.Ride) (string, error)
	GetRide(ctx context.Context, id string) (models.Ride, error)
}

type CoordinatesRepository interface {
	CreateNewCoordinate(ctx context.Context, c models.Coordinate) (string, error)
	GetCoordinate(ctx context.Context, id string) (models.Coordinate, error)
}
