package service

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"ride-hail/config"
	"ride-hail/internal/core/domain/action"
	"ride-hail/internal/core/domain/models"
	"ride-hail/internal/core/domain/types"
	"ride-hail/internal/core/ports"
	"ride-hail/internal/core/service/hash"
	"ride-hail/pkg/logger"
	"time"
)

type AuthService struct {
	secretKey string
	cfg       config.Config
	repo      ports.UserRepository
	log       *logger.Logger
}

func NewAuthService(cfg config.Config, repo ports.UserRepository, log *logger.Logger) *AuthService {
	return &AuthService{
		secretKey: cfg.JWT.Secret,
		repo:      repo,
		log:       log,
		cfg:       cfg,
	}
}

// Login returning token and error
func (s *AuthService) Login(ctx context.Context, reqId string, user models.User) (string, error) {
	log := s.log.Func("Login")

	u, err := s.repo.GetGyUserEmail(ctx, reqId, user.Email)
	if err != nil {
		return "", err
	}

	ok, err := hash.VerifyPassword(u.Password, user.Password)
	if err != nil {
		log.Error(
			action.Login,
			"error verifying password",
			"requestID", reqId,
			"userID", u.ID,
			"error", err,
		)
		return "", err
	}
	if !ok {
		log.Warn(
			action.Login,
			"incorrect password",
			"requestID", reqId,
		)
		return "", types.ErrIncorrectPassword
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
		log.Error(
			action.Login,
			"error generating JWT token",
			"requestID", reqId,
			"error", err,
		)
		return "", err
	}

	return tokenString, nil
}

func (s *AuthService) CreateNewUser(ctx context.Context, reqId string, user models.User) error {
	log := s.log.Func("CreateNewUser")

	hashPass, err := hash.HashPassword(user.Password)
	if err != nil {
		log.Error(
			action.Registration,
			"error hashing password",
			"requestID", reqId,
			"error", err,
		)
		return err
	}

	user.Password = hashPass
	err = s.repo.CreateNewUser(ctx, reqId, user)
	if err != nil {
		return err
	}
	return nil
}
