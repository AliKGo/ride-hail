package handle

import (
	"net/http"

	"ride-hail/internal/core/ports"
	"ride-hail/pkg/logger"
)

type DalHandle struct {
	svc ports.DalService
	log *logger.Logger
}

func NewDalHandle(svc ports.DalService, log *logger.Logger) *DalHandle {
	return &DalHandle{
		svc: svc,
		log: log,
	}
}

type DalHandler interface {
	DriverGoesOnline(w http.ResponseWriter, r *http.Request)
	DriverGoesOffline(w http.ResponseWriter, r *http.Request)
	UpdateDriverLocation(w http.ResponseWriter, r *http.Request)
	StartRide(w http.ResponseWriter, r *http.Request)
	CompleteRide(w http.ResponseWriter, r *http.Request)
}
