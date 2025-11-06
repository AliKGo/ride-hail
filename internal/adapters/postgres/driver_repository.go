package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"ride-hail/internal/core/domain/models"
	"ride-hail/pkg/executor"
)

type DriverRepository struct {
	pool *pgxpool.Pool
}

func NewDriverRepository(pool *pgxpool.Pool) *DriverRepository {
	return &DriverRepository{
		pool: pool,
	}
}

func (r *DriverRepository) Insert(ctx context.Context, newDriver models.Driver) error {
	ex := executor.GetExecutor(ctx, r.pool)

	query := `INSERT INTO drivers (id, license_number, vehicle_type, vehicle_attrs, status) VALUES ($1,$2, $3, $4, $5)`

	_, err := ex.Exec(ctx, query,
		newDriver.ID,
		newDriver.LicenseNumber,
		newDriver.VehicleType,
		newDriver.VehicleAttrs,
		newDriver.Status,
	)
	if err != nil {
		return fmt.Errorf("failed to insert driver: %w", err)
	}

	return nil
}

func (r *DriverRepository) Get(ctx context.Context, id string) (models.Driver, error) {
	ex := executor.GetExecutor(ctx, r.pool)

	query := `
		SELECT id, license_number, vehicle_type, vehicle_attrs, rating, total_rides, total_earnings, status, is_verified
		FROM drivers
		WHERE id = $1
	`
	var attrs []byte
	var driver models.Driver
	err := ex.QueryRow(ctx, query, id).Scan(
		&driver.ID,
		&driver.LicenseNumber,
		&driver.VehicleType,
		&attrs,
		&driver.Rating,
		&driver.TotalRides,
		&driver.TotalEarnings,
		&driver.Status,
		&driver.IsVarified,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Driver{}rrorf("driver not found: %w", err)
		}
		return models.Driver{}, fmt.Errorf("failed to get driver: %w", err)
	}

	if err = json.Unmarshal(attrs, &driver.VehicleAttrs); err != nil {
		return models.Driver{}, fmt.Errorf("failed to unmarshal driver attrs: %w", err)
	}

	return driver, nil
}

func (r *DriverRepository) UpdateStatus(ctx context.Context, id, status string) error {
	ex := executor.GetExecutor(ctx, r.pool)

	query := `UPDATE drivers SET status = $1 WHERE id = $2`
	cmdTag, err := ex.Exec(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("failed to update driver status: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("no driver found with id %s", id)
	}

	return nil
}

func (r *DriverRepository) InsertSession(ctx context.Context, driverID string) (string, error) {
	ex := executor.GetExecutor(ctx, r.pool)

	if driverID == "" {
		return "", fmt.Errorf("driverID cannot be empty")
	}

	query := `
		INSERT INTO driver_sessions (driver_id)
		VALUES ($1)
		RETURNING id
	`

	var sessionID string
	if err := ex.QueryRow(ctx, query, driverID).Scan(&sessionID); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23503":
				return "", fmt.Errorf("driver %s does not exist", driverID)
			default:
				return "", fmt.Errorf("insert driver session failed: %v (pgcode=%s)", pgErr.Message, pgErr.Code)
			}
		}
		return "", fmt.Errorf("failed to insert driver session: %w", err)
	}

	return sessionID, nil
}

func (r *DriverRepository) GetLastActiveSession(ctx context.Context, driverID string) (models.DriverSession, error) {
	ex := executor.GetExecutor(ctx, r.pool)

	query := `
		SELECT id, driver_id, started_at, ended_at, total_rides, total_earnings
		FROM driver_sessions
		WHERE driver_id = $1
		ORDER BY started_at DESC
		LIMIT 1
	`

	var session models.DriverSession
	err := ex.QueryRow(ctx, query, driverID).Scan(
		&session.ID,
		&session.DriverID,
		&session.StartedAt,
		&session.EndedAt,
		&session.TotalRides,
		&session.TotalEarnings,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.DriverSession{}, fmt.Errorf("no sessions found for driver %s", driverID)
		}
		return models.DriverSession{}, err
	}

	return session, nil
}
