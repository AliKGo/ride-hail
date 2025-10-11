package ports

import (
	"context"
	"ride-hail/internal/core/domain/models"
)

type AuthService interface {
	CreateNewUser(ctx context.Context, reqId string, user models.User) error
	Login(ctx context.Context, reqId string, user models.User) (string, error)
}

type UserRepository interface {
	CreateNewUser(ctx context.Context, reqId string, user models.User) error
	GetGyUserEmail(ctx context.Context, reqId string, email string) (models.User, error)
}
