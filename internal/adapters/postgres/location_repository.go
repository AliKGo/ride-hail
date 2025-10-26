package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"ride-hail/internal/core/domain/models"
	"ride-hail/pkg/executor"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LocationRepository struct {
	pool *pgxpool.Pool
}

func NewLocationRepository(pool *pgxpool.Pool) *LocationRepository {
	return &LocationRepository{
		pool: pool,
	}
}

func (repo *LocationRepository) SaveLocation(ctx context.Context, location models.LocationHistory) (string, error) {
	ex := executor.GetExecutor(ctx, repo.pool)

	query := `INSERT INTO location_history (id, coordinate_id, driver_id, latitude, 
	longitude, accuracy_meters, speed_kmh, heading_degrees, recorded_at, ride_id)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	RETURNING id;`

	var id string
	err := ex.QueryRow(
		ctx, query,
		location.ID,
		location.CoordinateID,
		location.DriverID,
		location.Latitude,
		location.Longitude,
		location.AccuracyMeters,
		location.SpeedKmh,
		location.HeadingDegrees,
		location.RecordedAt,
		location.RideID,
	).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("failed to save location: %w", err)
	}

	return id, nil
}

func (repo *LocationRepository) GetLastLocationByDriver(ctx context.Context, driverID string) (*models.LocationHistory, error) {
	ex := executor.GetExecutor(ctx, repo.pool)

	query := `SELECT id, coordinate_id, driver_id, latitude, 
	longitude, accuracy_meters, speed_kmh, heading_degrees, recorded_at, ride_id
	FROM location_history
	WHERE driver_id = $1
	ORDER BY recorded_at DESC
	LIMIT 1;`

	var location models.LocationHistory
	err := ex.QueryRow(ctx, query, driverID).Scan(
		&location.ID,
		&location.CoordinateID,
		&location.DriverID,
		&location.Latitude,
		&location.Longitude,
		&location.AccuracyMeters,
		&location.SpeedKmh,
		&location.HeadingDegrees,
		&location.RecordedAt,
		&location.RideID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get last location for driver %s: %w", driverID, err)
	}

	return &location, nil
}

func (repo *LocationRepository) GetLocationHistoryByDriver(ctx context.Context, driverID string, limit int) ([]models.LocationHistory, error) {
	ex := executor.GetExecutor(ctx, repo.pool)

	query := `SELECT id, coordinate_id, driver_id, latitude, 
	longitude, accuracy_meters, speed_kmh, heading_degrees, recorded_at, ride_id
	FROM location_history
	WHERE driver_id = $1
	ORDER BY recorded_at DESC
	LIMIT $2;`

	rows, err := ex.Query(ctx, query, driverID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get location history for driver %s: %w", driverID, err)
	}
	defer rows.Close()

	var locations []models.LocationHistory
	for rows.Next() {
		var location models.LocationHistory
		err := rows.Scan(
			&location.ID,
			&location.CoordinateID,
			&location.DriverID,
			&location.Latitude,
			&location.Longitude,
			&location.AccuracyMeters,
			&location.SpeedKmh,
			&location.HeadingDegrees,
			&location.RecordedAt,
			&location.RideID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan location history: %w", err)
		}
		locations = append(locations, location)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating location history: %w", err)
	}

	return locations, nil
}

func (repo *LocationRepository) DeleteLocationHistory(ctx context.Context, driverID string, before time.Time) error {
	ex := executor.GetExecutor(ctx, repo.pool)

	query := `DELETE FROM location_history WHERE driver_id = $1 AND recorded_at < $2;`

	result, err := ex.Exec(ctx, query, driverID, before)
	if err != nil {
		return fmt.Errorf("failed to delete location history: %w", err)
	}

	_ = result.RowsAffected()

	return nil
}
