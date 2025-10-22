package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"ride-hail/internal/core/domain/models"
	"ride-hail/internal/core/domain/types"
)

type RideRepository struct {
	pool *pgxpool.Pool
}

func NewRideRepository(pool *pgxpool.Pool) *RideRepository {
	return &RideRepository{
		pool: pool,
	}
}

func (repo *RideRepository) CreateNewRide(ctx context.Context, ride models.Ride) (string, error) {
	query := `INSERT INTO rides (
		ride_number, passenger_id, vehicle_type, status,
		estimated_fare, pickup_coordinate_id
	) VALUES ($1, $2, $3, $4, $5, $6,)
	RETURNING id`

	var id string
	err := repo.pool.QueryRow(
		ctx, query,
		ride.RideNumber,
		ride.PassengerID,
		ride.VehicleType,
		ride.Status,
		ride.EstimatedFare,
		ride.PickupCoordinateId,
	).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("failed to create ride: %w", err)
	}

	return id, nil
}

func (repo *RideRepository) GetRide(ctx context.Context, id string) (models.Ride, error) {
	query := `
	SELECT id, created_at, updated_at, ride_number, passenger_id, driver_id, vehicle_type,
	       status, priority, requested_at, matched_at, arrived_at, started_at,
	       completed_at, cancelled_at, cancellation_reason, estimated_fare, final_fare,
	       pickup_coordinate_id, destination_coordinate_id
	FROM rides
	WHERE id = $1
	`

	var ride models.Ride

	err := repo.pool.QueryRow(ctx, query, id).Scan(
		&ride.ID,
		&ride.CreatedAt,
		&ride.UpdatedAt,
		&ride.RideNumber,
		&ride.PassengerID,
		&ride.DriverID,
		&ride.VehicleType,
		&ride.Status,
		&ride.Priority,
		&ride.RequestedAt,
		&ride.MatchedAt,
		&ride.ArrivedAt,
		&ride.StartedAt,
		&ride.CompletedAt,
		&ride.CancelledAt,
		&ride.CancellationReason,
		&ride.EstimatedFare,
		&ride.FinalFare,
		&ride.PickupCoordinateId,
		&ride.DestinationCoordinateId,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Ride{}, types.ErrRideNotFound
		}
		return models.Ride{}, fmt.Errorf("failed to get ride by id %s: %w", id, err)
	}

	return ride, nil
}
