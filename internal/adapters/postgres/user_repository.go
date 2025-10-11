package postgres

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"ride-hail/internal/core/domain/action"
	"ride-hail/internal/core/domain/models"
	"ride-hail/pkg/logger"
	postgres "ride-hail/pkg/potgres"
)

type UserRepository struct {
	pool *pgxpool.Pool
	cfg  postgres.Config
	log  logger.Logger
}

func (p *UserRepository) CreateNewUser(ctx context.Context, reqId string, user models.User) error {
	query := `
	INSERT INTO users (email, role, password_hash)
	VALUES ($1, $2, $3)
	`

	_, err := p.pool.Exec(ctx, query, user.Email, user.Role, user.Password)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return models.ErrUserAlreadyExists
			}
		}
		p.log.Error(action.Registration, "error when inserting to the postgres", reqId, "", err)
		return err
	}
	return nil
}

func (p *UserRepository) GetGyUserEmail(ctx context.Context, reqId string, email string) (models.User, error) {
	query := "SELECT id, email, role, status, password_hash FROM users WHERE email = $1"

	var user models.User
	err := p.pool.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Role,
		&user.Status,
		&user.Password,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, models.ErrUserNotFound
		}
		p.log.Error(action.Login, "error in get user", reqId, "", err)
		return models.User{}, err
	}
	return user, nil
}
