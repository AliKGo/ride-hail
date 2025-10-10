package ports

import "ride-hail/internal/core/domain/models"

type Service interface {
	CreateNewUser(newUser models.User) (string, error)
}

type Repository interface {
	CreateNewUser(models.User) (string, error)
	GetUser(email string) (models.User, error)
}
