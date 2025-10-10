package handle

import (
	"encoding/json"
	"errors"
	"net/http"
	"ride-hail/internal/adapters/http/handle/dto"
	"ride-hail/internal/core/domain/action"
	"ride-hail/internal/core/domain/models"
	"ride-hail/internal/core/ports"
	"ride-hail/pkg/logger"
)

type ctxKey string

const reqIDKey ctxKey = "reqID"

type Handle struct {
	svc ports.Service
	log logger.Logger
}

func New(svc ports.Service) *Handle {
	return &Handle{svc: svc}
}

var (
	ErrorInValidateLogin = errors.New("validation error")
)

func (h *Handle) Registration(w http.ResponseWriter, r *http.Request) {
	var req models.User
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error(action.Registration, "error in parsing request", GetReqID(r), "", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if er, msg := dto.ValidateLogin(req.Email, req.Password); !er {
		h.log.Error(action.Registration, msg, "", "", ErrorInValidateLogin)
		http.Error(w, msg, http.StatusUnauthorized)
		return
	}

	id, err := h.svc.CreateNewUser(req)
	if err != nil {

	}
}

func GetReqID(r *http.Request) string {
	if v := r.Context().Value(reqIDKey); v != nil {
		if id, ok := v.(string); ok {
			return id
		}
	}
	return ""
}
