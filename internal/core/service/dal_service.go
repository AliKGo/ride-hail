package service

import (
	"context"
	"errors"
	"fmt"
	"ride-hail/internal/core/domain/action"
	"ride-hail/internal/core/domain/models"
	"ride-hail/internal/core/domain/types"
	"ride-hail/internal/core/ports"
	"ride-hail/pkg/logger"
	"ride-hail/pkg/txm"
	"time"
)

type DalService struct {
	log  *logger.Logger
	txm  txm.Manager
	repo dalRepository
}

type dalRepository struct {
	driver ports.DriversRepository
}

func NewDalService(log *logger.Logger, txm txm.Manager, driver ports.DriversRepository) *DalService {
	return &DalService{
		log: log,
		txm: txm,
		repo: dalRepository{
			driver: driver,
		},
	}
}

func (svc *DalService) CreateNewDriver(ctx context.Context, newDriver models.Driver) error {
	log := svc.log.Func("DalService.CreateNewDriver")

	if _, err := svc.repo.driver.Get(ctx, newDriver.ID); err == nil {
		log.Warn(ctx, action.Registration, "driver already exists")
		return types.ErrDriverExists
	}

	if err := svc.repo.driver.Insert(ctx, newDriver); err != nil {
		log.Error(ctx, action.Registration, "error when saving data in the database", err)
		return types.ErrInternalServiceError
	}

	return nil
}

func (svc *DalService) StatusOnline(ctx context.Context, id string, loc models.Position) (string, error) {
	log := svc.log.Func("DalService.StatusOnline")

	driver, err := svc.repo.driver.Get(ctx, id)
	if err != nil {
		log.Error(ctx, action.UpdateStatus, "error when getting data from the database", err)
		return "", types.ErrInternalServiceError
	} else if driver.Status != types.DriverStatusOffline {
		log.Warn(ctx, action.UpdateStatus, "driver is not offline")
		return "", types.ErrDriverOnline
	}

	if session, err := svc.repo.driver.GetLastActiveSession(ctx, driver.ID); err == nil {
		log.Error(ctx, action.UpdateStatus, "failed get last active session", "error", err)
		return "", types.ErrInternalServiceError
	} else if !session.EndedAt.IsZero() {
		log.Warn(ctx, action.UpdateStatus, "the last session didn't end", "error", err)
		if err = svc.repo.driver.UpdateStatus(ctx, driver.ID, types.DriverStatusOffline); err != nil {
			log.Error(ctx, action.UpdateStatus, "error when saving data in the database", err)
		}
		return "", types.ErrInternalServiceError
	}

	// проверка времени техосмотра
	if inspectionDate, err := time.Parse(time.DateOnly, driver.VehicleAttrs.InspectionDate); err != nil {
		log.Error(ctx, action.UpdateStatus, "error when parsing inspection date", err)
		return "", types.ErrInternalServiceError
	} else if isInspectionExpired(inspectionDate) {
		log.Warn(ctx, action.UpdateStatus, "inspection date is expired")
		return "", fmt.Errorf("inspection date is expired")
	}

	// срок действий страховки
	if insuranceExpiry, err := time.Parse(time.DateOnly, driver.VehicleAttrs.InsuranceExpiry); err != nil {
		log.Error(ctx, action.UpdateStatus, "error when parsing insurance expiry", err)
		return "", types.ErrInternalServiceError
	} else if err = validateInsurance(insuranceExpiry); err != nil {
		log.Error(ctx, action.UpdateStatus, "the insurance period has expired")
		return "", err
	}

	// срок действий лицензии на такси
	if taxiLicenseExpiry, err := time.Parse(time.DateOnly, driver.VehicleAttrs.TaxiLicenseExpiry); err != nil {
		log.Error(ctx, action.UpdateStatus, "error when parsing taxi license expiry", err)
		return "", fmt.Errorf("error when parsing taxi license expiry: %w", err)
	} else if err = validateInsurance(taxiLicenseExpiry); err != nil {
		log.Error(ctx, action.UpdateStatus, "the insurance period has expired")
		return "", err
	}

	var sessionId string
	fn := func(ctx context.Context) error {
		if sessionId, err = svc.repo.driver.InsertSession(ctx, id); err != nil {
			log.Error(ctx, action.UpdateStatus, "error when saving data in the database", err)
			return types.ErrInternalServiceError
		}
		return nil
	}

	svc.txm.Do(ctx, fn)
	return sessionId, nil
}

func isInspectionExpired(inspectionDate time.Time) bool {
	return time.Now().After(inspectionDate.AddDate(0, 6, 0))
}

func validateInsurance(expiry time.Time) error {
	if time.Now().After(expiry) {
		return errors.New("the insurance period has expired")
	}
	return nil
}
