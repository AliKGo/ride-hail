package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
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
func (s *AuthService) Login(ctx context.Context, user models.User) (string, string, error) {
	log := s.log.Func("Login")

	u, err := s.repo.GetGyUserEmail(ctx, user.Email)
	if err != nil {
		log.Error(ctx, action.Login, "error in getting user email", "email", user.Email, "error", err)
		return "", "", err
	}

	ok, err := hash.VerifyPassword(u.Password, user.Password)
	if err != nil {
		log.Error(ctx, action.Login, "error verifying password", "userID", u.ID, "error", err)
		return "", "", err
	}
	if !ok {
		log.Warn(ctx, action.Login, "incorrect password")
		return "", "", types.ErrIncorrectPassword
	}

	claims := models.Claims{
		ClaimsID: newClaimsID(),
		UserID:   u.ID,
		Role:     u.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		log.Error(ctx, action.Login, "error generating JWT token", "error", err)
		return "", "", err
	}

	return tokenString, u.ID, nil
}

func (s *AuthService) CreateNewUser(ctx context.Context, user models.User) error {
	log := s.log.Func("CreateNewUser")

	hashPass, err := hash.HashPassword(user.Password)
	if err != nil {
		log.Error(ctx, action.Registration, "error hashing password", "error", err)
		return err
	}

	user.Password = hashPass
	err = s.repo.CreateNewUser(ctx, user)
	if err != nil {
		log.Error(ctx, action.Registration, "error creating new user", "error", err)
		return err
	}
	return nil
}

func newClaimsID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return hex.EncodeToString(fmt.Appendf(nil, "%d", time.Now().UnixNano()))
	}
	return hex.EncodeToString(b)
}
