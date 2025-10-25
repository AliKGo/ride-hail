package postgres

import (
	"context"
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

func (repo *DriverRepository) UpdateDriverStatus(ctx context.Context, dal models.Driver) (string, error) {
	ex := executor.GetExecutor(ctx, repo.pool)

	query := `
		UPDATE drivers
		SET status = $1,
			updated_at = now()
		WHERE id = $2
		RETURNING id;
		`

	var id string
	err := ex.QueryRow(
		ctx, query,
		dal.Status,
		dal.ID,
	).Scan(&id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", fmt.Errorf("no driver found with id %s", dal.ID)
		}
		return "", fmt.Errorf("failed to update status: %w", err)
	}

	return id, nil
}
