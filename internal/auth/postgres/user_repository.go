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

func (p *UserRepository) CreateNewUser(ctx context.Context, user models.User) error {
	log := p.log.Func("CreateNewUser")

	log.Info(ctx,
		action.Registration,
		"start creating user",
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
				ctx,
				action.Registration,
				"user already exists",
				"email", user.Email,
			)
			return types.ErrUserAlreadyExists
		}
		log.Error(
			ctx,
			action.Registration,
			"error when inserting user into postgres",
			"error", err.Error(),
		)

		return err
	}

	log.Info(
		ctx,
		action.Registration,
		"user successfully created",
	)
	return nil
}

func (p *UserRepository) GetGyUserEmail(ctx context.Context, email string) (models.User, error) {
	log := p.log.Func("GetGyUserEmail")

	log.Debug(
		ctx,
		action.Login,
		"start getting user by email",
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
				ctx,
				action.Login,
				"user not found",
			)
			return models.User{}, types.ErrUserNotFound
		}

		log.Error(
			ctx,
			action.Login,
			"error when getting user by email",
			"error", err.Error(),
		)

		return models.User{}, err
	}

	log.Info(
		ctx,
		action.Login,
		"user successfully retrieved",
	)

	return user, nil
}
