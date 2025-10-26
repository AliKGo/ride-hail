package postgres

import (
	"context"
	"errors"
	"fmt"

	"ride-hail/internal/core/domain/models"
	"ride-hail/pkg/executor"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DriverRepository struct {
	pool *pgxpool.Pool
}

func NewDriverRepository(pool *pgxpool.Pool) *DriverRepository {
	return &DriverRepository{
		pool: pool,
	}
}

func (repo *DriverRepository) CreateDriver(ctx context.Context, driver models.Driver) (string, error) {
	ex := executor.GetExecutor(ctx, repo.pool)

	query := `INSERT INTO drivers
			(id, created_at, updated_at, license_number, vehicle_type, vehicle_attrs, 
							rating, total_rides, total_earnings, status, is_verified)
			VALUES ($1, now(), now(), $2, $3, $4::jsonb, $5, $6, $7, $8, $9)
			RETURNING id;`

	var id string
	err := ex.QueryRow(
		ctx, query,
		driver.ID,
		driver.LicenseNumber,
		driver.VehicleType,
		driver.VehicleAttrs,
		driver.Rating,
		driver.TotalRides,
		driver.TotalEarnings,
		driver.Status,
		driver.IsVerified,
	).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("failed to create driver: %w", err)
	}

	return id, nil
}

func (repo *DriverRepository) GetDriverByID(ctx context.Context, id string) (*models.Driver, error) {
	ex := executor.GetExecutor(ctx, repo.pool)

	query := `SELECT id, created_at, updated_at, license_number, vehicle_type, 
	vehicle_attrs, rating, total_rides, total_earnings, status, is_verified
			FROM drivers WHERE id = $1;`

	var driver models.Driver

	err := ex.QueryRow(ctx, query, id).Scan(
		&driver.ID,
		&driver.CreatedAt,
		&driver.UpdatedAt,
		&driver.LicenseNumber,
		&driver.VehicleType,
		&driver.VehicleAttrs,
		&driver.Rating,
		&driver.TotalRides,
		&driver.TotalEarnings,
		&driver.Status,
		&driver.IsVerified,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("driver not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get driver by id %s: %w", id, err)
	}
	return &driver, nil
}

func (repo *DriverRepository) UpdateDriver(ctx context.Context, driver models.Driver) error {
	ex := executor.GetExecutor(ctx, repo.pool)

	query := `UPDATE drivers
		SET license_number = $2, vehicle_type = $3, vehicle_attrs = $4::jsonb, rating = $5, 
		total_rides = $6, total_earnings = $7, status = $8, is_verified = $9, updated_at = now()
		WHERE id = $1;`

	result, err := ex.Exec(
		ctx, query,
		driver.ID,
		driver.LicenseNumber,
		driver.VehicleType,
		driver.VehicleAttrs,
		driver.Rating,
		driver.TotalRides,
		driver.TotalEarnings,
		driver.Status,
		driver.IsVerified,
	)
	if err != nil {
		return fmt.Errorf("failed to update driver: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("no driver found with id %s", driver.ID)
	}

	return nil
}

func (repo *DriverRepository) DeleteDriver(ctx context.Context, id string) error {
	ex := executor.GetExecutor(ctx, repo.pool)

	query := `DELETE FROM drivers WHERE id = $1;`

	result, err := ex.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete driver: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("no driver found with id %s", id)
	}
	return nil
}

func (repo *DriverRepository) ListDriversByStatus(ctx context.Context, status string, limit, offset int) ([]models.Driver, error) {
	ex := executor.GetExecutor(ctx, repo.pool)

	query := `SELECT id, created_at, updated_at, license_number, vehicle_type, 
	vehicle_attrs, rating, total_rides, total_earnings, status, is_verified
	FROM drivers 
	WHERE status = $1
	LIMIT $2 OFFSET $3;`

	rows, err := ex.Query(ctx, query, status, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list drivers by status: %w", err)
	}
	defer rows.Close()

	var drivers []models.Driver
	for rows.Next() {
		var driver models.Driver
		err := rows.Scan(
			&driver.ID,
			&driver.CreatedAt,
			&driver.UpdatedAt,
			&driver.LicenseNumber,
			&driver.VehicleType,
			&driver.VehicleAttrs,
			&driver.Rating,
			&driver.TotalRides,
			&driver.TotalEarnings,
			&driver.Status,
			&driver.IsVerified,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan driver: %w", err)
		}
		drivers = append(drivers, driver)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating drivers: %w", err)
	}

	return drivers, nil
}

func (repo *DriverRepository) UpdateDriverStatus(ctx context.Context, driverID string, newStatus string) error {
	ex := executor.GetExecutor(ctx, repo.pool)

	query := `UPDATE drivers
		SET status = $1, updated_at = now()
		WHERE id = $2;`

	result, err := ex.Exec(ctx, query, newStatus, driverID)
	if err != nil {
		return fmt.Errorf("failed to update driver status: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("no driver found with id %s", driverID)
	}

	return nil
}
