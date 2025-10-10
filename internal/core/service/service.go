package service

import (
	"ride-hail/internal/core/domain/models"
	"ride-hail/internal/core/ports"
	"ride-hail/pkg/logger"
)

type Service struct {
	repo ports.Repository
	log  logger.Logger
}

func NewService(repo ports.Repository, log logger.Logger) *Service {
	return &Service{
		repo: repo,
		log:  log,
	}
}

func (s *Service) CreateNewUser(req models.User) (string, error) {
	var id string

	return id, nil
}
