package postgres

import (
	"context"
	"errors"
	"fmt"
	"ride-hail/internal/core/domain/models"
	"ride-hail/internal/core/domain/types"
	"ride-hail/pkg/executor"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
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
	ex := executor.GetExecutor(ctx, repo.pool)

	query := `INSERT INTO rides (
		ride_number, passenger_id, vehicle_type, status,
		estimated_fare, pickup_coordinate_id, destination_coordinate_id
	) VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING id`

	var id string
	err := ex.QueryRow(
		ctx, query,
		ride.RideNumber,
		ride.PassengerID,
		ride.VehicleType,
		ride.Status,
		ride.EstimatedFare,
		ride.PickupCoordinateId,
		ride.DestinationCoordinateId,
	).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("failed to create ride: %w", err)
	}

	return id, nil
}

func (repo *RideRepository) GetRide(ctx context.Context, id string) (models.Ride, error) {
	ex := executor.GetExecutor(ctx, repo.pool)

	query := `
	SELECT id, created_at, updated_at, ride_number, passenger_id, driver_id, vehicle_type,
	       status, priority, requested_at, matched_at, arrived_at, started_at,
	       completed_at, cancelled_at, cancellation_reason, estimated_fare, final_fare,
	       pickup_coordinate_id, destination_coordinate_id
	FROM rides
	WHERE id = $1
	`

	var ride models.Ride

	err := ex.QueryRow(ctx, query, id).Scan(
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

func (repo *RideRepository) GenerateRideNumber(ctx context.Context) (int, error) {
	ex := executor.GetExecutor(ctx, repo.pool)
	var counter int
	err := ex.QueryRow(ctx,
		`
        INSERT INTO ride_counters (ride_date, counter) 
        VALUES (CURRENT_DATE, 1)
        ON CONFLICT (ride_date) DO UPDATE
        SET counter = ride_counters.counter + 1
        RETURNING counter
    `).Scan(&counter)

	if err != nil {
		return 0, err
	}
	return counter, nil
}

func (repo *RideRepository) UpdateRide(ctx context.Context, rideID string, newStatus string, reason string, t *time.Time) error {
	ex := executor.GetExecutor(ctx, repo.pool)

	var timeField string
	switch newStatus {
	case types.RideStatusMATCHED:
		timeField = "matched_at"
	case types.RideStatusEN_ROUTE:
		timeField = "arrived_at"
	case types.RideStatusARRIVED:
		timeField = "arrived_at"
	case types.RideStatusIN_PROGRESS:
		timeField = "started_at"
	case types.RideStatusCOMPLETED:
		timeField = "completed_at"
	case types.RideStatusCANCELLED:
		timeField = "cancelled_at"
	}

	var query string
	var args []any

	if timeField != "" {
		if newStatus == types.RideStatusCANCELLED {
			query = fmt.Sprintf(`
				UPDATE rides
				SET 
					status = $1,
					%s = COALESCE(%s, $2),
					cancellation_reason = CASE WHEN $3 != '' THEN $3 ELSE cancellation_reason END,
					updated_at = now()
				WHERE id = $4
			`, timeField, timeField)
			args = []any{newStatus, t, reason, rideID}
		} else {
			query = fmt.Sprintf(`
				UPDATE rides
				SET 
					status = $1,
					%s = COALESCE(%s, $2),
					updated_at = now()
				WHERE id = $3
			`, timeField, timeField)
			args = []any{newStatus, t, rideID}
		}
	} else {
		query = `
			UPDATE rides
			SET 
				status = $1,
				updated_at = now()
			WHERE id = $2
		`
		args = []any{newStatus, rideID}
	}

	_, err := ex.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update ride status: %w", err)
	}

	return nil
}

func (r *RideRepository) UpdateMatchedRide(ctx context.Context, rideID, driverID string, matchedAt time.Time) error {
	ex := executor.GetExecutor(ctx, r.pool)

	query := `
UPDATE rides
SET driver_id = $1,
    matched_at = $2,
    updated_at = now()
WHERE id = $3
`

	_, err := ex.Exec(ctx, query, driverID, matchedAt, rideID)
	if err != nil {
		return err
	}

	return nil
}
