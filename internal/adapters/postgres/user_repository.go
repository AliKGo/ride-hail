package postgres

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"ride-hail/internal/core/domain/models"
	"ride-hail/internal/core/domain/types"
	"ride-hail/pkg/executor"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewRepo(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		pool: pool,
	}
}

func (repo *UserRepository) CreateNewUser(ctx context.Context, user models.User) error {
	ex := executor.GetExecutor(ctx, repo.pool)

	query := `
	INSERT INTO users (email, role, password_hash)
	VALUES ($1, $2, $3)
	`

	_, err := ex.Exec(ctx, query, user.Email, user.Role, user.Password)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return types.ErrUserAlreadyExists
		}
		return err
	}
	return nil
}

func (repo *UserRepository) GetGyUserEmail(ctx context.Context, email string) (models.User, error) {
	ex := executor.GetExecutor(ctx, repo.pool)

	query := `SELECT id, email, role, status, password_hash FROM users WHERE email = $1`

	var user models.User
	err := ex.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Role,
		&user.Status,
		&user.Password,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, types.ErrUserNotFound
		}
		return models.User{}, err
	}
	return user, nil
}
