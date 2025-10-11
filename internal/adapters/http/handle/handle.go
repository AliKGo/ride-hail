package handle

import (
	"encoding/json"
	"errors"
	"net/http"
	"ride-hail/config"
	"ride-hail/internal/adapters/http/handle/dto"
	"ride-hail/internal/core/domain/action"
	"ride-hail/internal/core/domain/models"
	"ride-hail/internal/core/ports"
	"ride-hail/pkg/logger"
)

type ctxKey string

const reqIDKey ctxKey = "reqID"

type Handle struct {
	svc ports.AuthService
	log logger.Logger
	cfg config.Config
}

func New(cfg config.Config, svc ports.AuthService, log logger.Logger) *Handle {
	return &Handle{
		svc: svc,
		log: log,
		cfg: cfg,
	}
}

var (
	ErrorInValidateLogin = errors.New("validation error")
)

func (h *Handle) Registration(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqID := GetReqID(r)

	var req models.User
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error(action.Registration, "error in parsing request", reqID, "", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if er, msg := dto.ValidateLogin(req.Email, req.Password); !er {
		h.log.Error(action.Registration, msg, reqID, "", ErrorInValidateLogin)
		http.Error(w, msg, http.StatusUnauthorized)
		return
	}

	err := h.svc.CreateNewUser(ctx, reqID, req)
	if err != nil {
		if errors.Is(err, models.ErrUserAlreadyExists) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *Handle) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqID := GetReqID(r)
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		h.log.Error(action.Login, "error in parsing request", reqID, "", err)
		return
	}

	token, err := h.svc.Login(ctx, reqID, user)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		} else if errors.Is(err, models.ErrIncorrectPassword) {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "auth",
		Value:    token,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
		MaxAge:   h.cfg.JWT.ExpireHours * 60 * 60,
	})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "login successful"})
}

func GetReqID(r *http.Request) string {
	if v := r.Context().Value(reqIDKey); v != nil {
		if id, ok := v.(string); ok {
			return id
		}
	}
	return ""
}
