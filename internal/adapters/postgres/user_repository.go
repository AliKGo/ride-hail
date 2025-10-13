package postgres

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"ride-hail/internal/core/domain/action"
	"ride-hail/internal/core/domain/models"
	"ride-hail/internal/core/domain/types"
	"ride-hail/pkg/logger"
)

type UserRepository struct {
	pool *pgxpool.Pool
	log  *logger.Logger
}

func NewRepo(pool *pgxpool.Pool, log *logger.Logger) *UserRepository {
	return &UserRepository{
		pool: pool,
		log:  log,
	}
}

func (p *UserRepository) CreateNewUser(ctx context.Context, reqId string, user models.User) error {
	log := p.log.Func("CreateNewUser")

	log.Info(
		action.Registration,
		"start creating user",
		"requestID", reqId,
	)

	query := `
	INSERT INTO users (email, role, password_hash)
	VALUES ($1, $2, $3)
	`

	_, err := p.pool.Exec(ctx, query, user.Email, user.Role, user.Password)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			log.Warn(
				action.Registration,
				"user already exists",
				"requestID", reqId,
				"email", user.Email,
			)
			return types.ErrUserAlreadyExists
		}
		log.Error(
			action.Registration,
			"error when inserting user into postgres",
			"requestID", reqId,
			"error", err.Error(),
		)

		return err
	}

	log.Info(
		action.Registration,
		"user successfully created",
		"requestID", reqId,
	)
	return nil
}

// GetGyUserEmail dsvsdf
func (p *UserRepository) GetGyUserEmail(ctx context.Context, reqId string, email string) (models.User, error) {
	log := p.log.Func("GetGyUserEmail")

	log.Debug(
		action.Login,
		"start getting user by email",
		"requestID", reqId,
		"email", email,
	)

	query := `SELECT id, email, role, status, password_hash FROM users WHERE email = $1`

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
			log.Warn(
				action.Login,
				"user not found",
				"request_id", reqId,
			)
			return models.User{}, types.ErrUserNotFound
		}

		log.Error(
			action.Login,
			"error when getting user by email",
			"request_id", reqId,
			"error", err.Error(),
		)

		return models.User{}, err
	}

	log.Info(
		action.Login,
		"user successfully retrieved",
		"requestID", reqId,
	)

	return user, nil
}
