package handle

import (
	"encoding/json"
	"net/http"
	"ride-hail/internal/adapters/http/handle/dto"
	"ride-hail/internal/adapters/http/websocket"
	"ride-hail/internal/core/domain/action"
	"ride-hail/internal/core/domain/models"
	"ride-hail/internal/core/domain/types"
	"ride-hail/internal/core/ports"
	"ride-hail/pkg/logger"
	"strings"
)

type RideHandle struct {
	svc ports.RideService
	log *logger.Logger
}

func NewRideHandle(svc ports.RideService, wsm websocket.PassengerWSHandler, log *logger.Logger) *RideHandle {
	return &RideHandle{
		svc: svc,
		log: log,
	}
}

type RideHandler interface {
	CreateNewRide(w http.ResponseWriter, r *http.Request)
	CancelRide(w http.ResponseWriter, r *http.Request)
}

func (h *RideHandle) CreateNewRide(w http.ResponseWriter, r *http.Request) {
	log := h.log.Func("RideHandle.CreateNewRide")
	ctx := r.Context()

	log.Debug(ctx, action.CreateRide, "request to create a Ride has been launched")

	if logger.GetRole(ctx) != types.RoleCustomer {
		log.Error(ctx, action.CreateRide, "invalid role", "role", logger.GetRole(ctx))
		writeJSON(w, http.StatusForbidden, http.StatusText(http.StatusForbidden))
		return
	}
	var rideDto models.CreateRideRequest

	if err := json.NewDecoder(r.Body).Decode(&rideDto); err != nil {
		log.Error(ctx, "decode error", "msg", "err", err.Error())
		writeJSON(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if ok, err := dto.ValidateRideDTO(rideDto); !ok {
		log.Warn(ctx, action.CreateRide, "invalid request")
		writeJSON(w, http.StatusBadRequest, err)
		return
	}

	if resp, err := h.svc.CreateNewRide(ctx, rideDto); err != nil {
		writeJSON(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	} else {
		log.Debug(ctx, action.CreateRide, "the request to create a trip was successfully completed")
		writeJSON(w, http.StatusCreated, resp)
		return
	}
}

func (h *RideHandle) CancelRide(w http.ResponseWriter, r *http.Request) {
	log := h.log.Func("RideHandle.CancelRide")
	ctx := r.Context()

	log.Debug(ctx, action.CloseRide, "a request to close the ride has been launched")
	if logger.GetRole(ctx) != types.RoleCustomer {
		log.Error(ctx, action.CloseRide, "invalid role")
		writeJSON(w, http.StatusForbidden, http.StatusText(http.StatusForbidden))
		return
	}

	closeReq := models.CloseRideRequest{
		RideID: getRideID(r),
	}

	if err := json.NewDecoder(r.Body).Decode(&closeReq); err != nil {
		log.Error(ctx, action.CloseRide, "error decoding body", "error", err)
		writeJSON(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if resp, err := h.svc.CloseRide(ctx, closeReq); err != nil {
		writeJSON(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	} else {
		log.Debug(ctx, action.CloseRide, "the request to cancel the ride has been completed")
		writeJSON(w, http.StatusOK, resp)
		return
	}
}

func getRideID(r *http.Request) string {
	path := r.URL.Path
	parts := strings.Split(path, "/")
	return parts[2]
}
