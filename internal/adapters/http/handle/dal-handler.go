package handle

import (
	"encoding/json"
	"net/http"
	"ride-hail/internal/adapters/http/handle/dto"
	"ride-hail/internal/core/domain/action"
	"ride-hail/internal/core/domain/models"
	"ride-hail/internal/core/domain/types"
	"ride-hail/internal/core/ports"
	"ride-hail/pkg/logger"
	"strings"
)

type DalHandler struct {
	svc ports.DalService
	log *logger.Logger
}

func NewDalHandler(log *logger.Logger) *DalHandler {
	return &DalHandler{log: log}
}

func (h *DalHandler) Registration(w http.ResponseWriter, r *http.Request) {
	log := h.log.Func("DalHandler.Registration")
	ctx := r.Context()

	log.Debug(ctx, action.Registration, "registration request started")

	if logger.GetRole(ctx) != types.RoleCustomer {
		log.Error(ctx, action.Registration, "invalid role")
		writeJSON(w, http.StatusForbidden, "invalid role")
		return
	}

	var data dto.DriverRegistration
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Error(ctx, "decode error", "msg", "err", err.Error())
		writeJSON(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if errMsg := data.Validate(); errMsg != "" {
		log.Error(ctx, "validate error", errMsg)
		writeJSON(w, http.StatusBadRequest, errMsg)
		return
	}

	if logger.GetUserID(ctx) == "" {
		log.Error(ctx, action.Registration, "invalid user_id")
		writeJSON(w, http.StatusBadRequest, "invalid user_id")
		return
	}

	if err := h.svc.CreateNewDriver(ctx, models.Driver{
		ID:            logger.GetUserID(ctx),
		LicenseNumber: data.LicenseNumber,
		VehicleType:   data.VehicleType,
		VehicleAttrs:  data.VehicleAttrs,
		Status:        types.DriverStatusOffline,
	}); err != nil {
		writeJSON(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{
		"status":    "created",
		"driver_id": logger.GetUserID(ctx),
	})

	log.Debug(ctx, action.Registration, "registration request finished")
}

func (h *DalHandler) DriverGoesOnline(w http.ResponseWriter, r *http.Request) {
	log := h.log.Func("DalHandler.DriverGoesOnline")
	ctx := r.Context()
	if logger.GetRole(ctx) != types.RoleDriver {
		log.Error(ctx, action.UpdateStatus, "invalid role")
		writeJSON(w, http.StatusForbidden, "invalid role")
		return
	}

	if logger.GetUserID(ctx) != extractDriverID(r) && logger.GetUserID(ctx) != "" {
		log.Error(ctx, action.UpdateStatus, "invalid driver_id")
		writeJSON(w, http.StatusBadRequest, "invalid driver_id")
		return
	}

	var location dto.Location
	if err := json.NewDecoder(r.Body).Decode(&location); err != nil {
		log.Error(ctx, "decode error", "msg", "err", err.Error())
		writeJSON(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if errMsg := location.Validate(); errMsg != "" {
		log.Error(ctx, "validate error", errMsg)
		writeJSON(w, http.StatusBadRequest, errMsg)
		return
	}

	sessionID, err := h.svc.StatusOnline(ctx, logger.GetUserID(ctx), models.Position{
		Latitude:  location.Latitude,
		Longitude: location.Longitude,
	})

	if err != nil {
		writeJSON(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"status":  "AVAILABLE",
		"session": sessionID,
		"message": "You are now online and ready to accept rides"},
	)
}

func (h *DalHandler) DriverGoesOffline(w http.ResponseWriter, r *http.Request) {
	log := h.log.Func("DalHandler.DriverGoesOffline")
	ctx := r.Context()
	if logger.GetRole(ctx) != types.RoleDriver {
		log.Error(ctx, action.UpdateStatus, "invalid role")
		writeJSON(w, http.StatusForbidden, "invalid role")
	}

	if logger.GetUserID(ctx) != extractDriverID(r) && logger.GetUserID(ctx) != "" {
		log.Error(ctx, action.UpdateStatus, "invalid driver_id")
		writeJSON(w, http.StatusBadRequest, "invalid driver_id")
		return
	}

	if driverInfo, err := h.svc.StatusClose(ctx, logger.GetUserID(ctx)); err != nil {
		writeJSON(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	} else {
		writeJSON(w, http.StatusOK, driverInfo)
		return
	}
}

func extractDriverID(r *http.Request) string {
	path := r.URL.Path
	parts := strings.Split(strings.TrimPrefix(path, "/"), "/")

	for i, part := range parts {
		if part == "drivers" && i+1 < len(parts) {
			return parts[i+1]
		}
	}
	return ""
}
