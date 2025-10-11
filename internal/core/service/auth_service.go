package service

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"ride-hail/config"
	"ride-hail/internal/core/domain/action"
	"ride-hail/internal/core/domain/models"
	"ride-hail/internal/core/ports"
	"ride-hail/internal/core/service/hash"
	"ride-hail/pkg/logger"
	"time"
)

type AuthService struct {
	secretKey string
	cfg       config.Config
	repo      ports.UserRepository
	log       logger.Logger
}

func NewAuthService(key string, cfg config.Config, repo ports.UserRepository, log logger.Logger) *AuthService {
	return &AuthService{
		secretKey: key,
		repo:      repo,
		log:       log,
		cfg:       cfg,
	}
}

// returning token and error
func (s *AuthService) Login(ctx context.Context, reqId string, user models.User) (string, error) {
	u, err := s.repo.GetGyUserEmail(ctx, reqId, user.Email)
	if err != nil {
		return "", err
	}

	ok, err := hash.VerifyPassword(u.Password, user.Password)
	if err != nil {
		s.log.Error(action.Login, "error in varify password", reqId, "", err)
		return "", err
	}
	if !ok {
		return "", models.ErrIncorrectPassword
	}

	claims := models.Claims{
		UserID: u.ID,
		Role:   u.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		s.log.Error(action.Login, "error in generate token", reqId, "", err)
		return "", err
	}

	return tokenString, nil
}

func (s *AuthService) CreateNewUser(ctx context.Context, reqId string, user models.User) error {
	hashPass, err := hash.HashPassword(user.Password)
	if err != nil {
		s.log.Error(action.Login, "error in hash password", reqId, "", err)
		return err
	}

	user.Password = hashPass
	return s.repo.CreateNewUser(ctx, reqId, user)
}
