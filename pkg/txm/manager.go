package txm

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"ride-hail/internal/core/domain/models"
)

type TXManager struct {
	pool *pgxpool.Pool
}

func NewTXManager(pool *pgxpool.Pool) *TXManager {
	return &TXManager{
		pool: pool,
	}
}

type Manager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

func (T *TXManager) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := T.pool.Begin(ctx)
	if err != nil {
		return err
	}
	ctx = context.WithValue(ctx, models.GetTxKey(), tx)

	if err = fn(ctx); err != nil {
		if errR := tx.Rollback(ctx); errR != nil {
			return fmt.Errorf("tx rollback: %w: %w", errR, err)
		}
		return err
	}
	return tx.Commit(ctx)
}
