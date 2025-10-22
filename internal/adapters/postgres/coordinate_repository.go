package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"ride-hail/internal/core/domain/models"
)

type CordRepository struct {
	pool *pgxpool.Pool
}

func NewCordRepository(pool *pgxpool.Pool) *CordRepository {
	return &CordRepository{
		pool: pool,
	}
}

func (repo *CordRepository) CreateNewCoordinate(ctx context.Context, c models.Coordinate) (string, error) {
	query := `INSERT INTO coordinates 
		(entity_id, entity_type, address, latitude, longitude,
		fare_amount, distance_km, duration_minutes, is_current)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	RETURNING id`

	var id string
	err := repo.pool.QueryRow(
		ctx, query,
		c.EntityID, c.EntityType, c.Address,
		c.Latitude, c.Longitude, c.FareAmount, c.DistanceKM, c.DurationMinutes, c.IsCurrent,
	).Scan(&id)

	if err != nil {
		return "", fmt.Errorf("failed to insert coordinate: %w", err)
	}

	return id, nil
}

func (repo *CordRepository) GetCoordinate(ctx context.Context, id string) (models.Coordinate, error) {
	query := `SELECT id, created_at, updated_at, entity_id, entity_type, address,
       latitude, longitude, fare_amount, distance_km, duration_minutes, is_current
FROM coordinates
WHERE id = $1
`
	row := repo.pool.QueryRow(ctx, query, id)
	var coordinate models.Coordinate
	err := row.Scan(
		&coordinate.ID,
		&coordinate.CreatedAt,
		&coordinate.UpdatedAt,
		&coordinate.EntityID,
		&coordinate.EntityType,
		&coordinate.Address,
		&coordinate.Latitude,
		&coordinate.Longitude,
		&coordinate.FareAmount,
		&coordinate.DistanceKM,
		&coordinate.DurationMinutes,
		&coordinate.IsCurrent,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Coordinate{}, fmt.Errorf("failed to get coordinate by id %s: %w", id, err)
		}

	}

	return coordinate, err
}
