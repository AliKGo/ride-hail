package ports

import (
	"context"
	"net/http"
	"ride-hail/internal/core/domain/models"
)

type AuthService interface {
	CreateNewUser(ctx context.Context, user models.User) error
	Login(ctx context.Context, user models.User) (string, error)
}

type UserRepository interface {
	CreateNewUser(ctx context.Context, user models.User) error
	GetGyUserEmail(ctx context.Context, email string) (models.User, error)
}

type AuthHandler interface {
	Login(w http.ResponseWriter, r *http.Request)
	Registration(w http.ResponseWriter, r *http.Request)
}

type AuthServices interface {
	Run() error
}
