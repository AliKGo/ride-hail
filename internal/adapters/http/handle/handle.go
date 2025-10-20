package handle

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"ride-hail/config"
	"ride-hail/internal/adapters/http/handle/dto"
	"ride-hail/internal/core/domain/action"
	"ride-hail/internal/core/domain/models"
	"ride-hail/internal/core/domain/types"
	"ride-hail/internal/core/ports"
	"ride-hail/pkg/logger"
)

type Handle struct {
	svc ports.AuthService
	log *logger.Logger
	cfg config.Config
}

func New(cfg config.Config, svc ports.AuthService, log *logger.Logger) *Handle {
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
	log := h.log.Func("Registration")
	ctx := r.Context()
	reqID := logger.GetRequestID(ctx)

	log.Info(
		ctx,
		action.Registration,
		"registration request started",
		"requestID", reqID,
	)

	var req models.User
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error(
			ctx,
			action.Registration,
			"error parsing request body",
			"requestID", reqID,
			"error", err,
		)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if er, msg := dto.ValidateLogin(req.Email, req.Password); !er {
		log.Error(
			ctx,
			action.Registration, msg,
			"requestID", reqID,
			"error", ErrorInValidateLogin,
		)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	dto.GetRole(h.cfg.Mode, &req)

	err := h.svc.CreateNewUser(ctx, req)
	if err != nil {
		if errors.Is(err, types.ErrUserAlreadyExists) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"message": "registration successful"})
}

func (h *Handle) Login(w http.ResponseWriter, r *http.Request) {
	log := h.log.Func("Login")
	ctx := r.Context()

	log.Info(ctx, action.Login, "login request started")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error(
			ctx,
			action.Login,
			"error reading body",
			"error", err,
		)

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(body) == 0 {
		log.Warn(
			ctx,
			action.Login,
			"empty request body",
		)

		http.Error(w, "empty request body", http.StatusBadRequest)
		return
	}

	var user models.User
	if err = json.Unmarshal(body, &user); err != nil {
		log.Error(
			ctx,
			action.Login,
			"invalid JSON",
			"error", err,
		)
		return
	}

	token, err := h.svc.Login(ctx, user)
	if err != nil {
		switch {
		case errors.Is(err, types.ErrUserNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		case errors.Is(err, types.ErrIncorrectPassword):
			http.Error(w, err.Error(), http.StatusUnauthorized)
		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "Authorization",
		Value:    token,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
		MaxAge:   h.cfg.JWT.ExpireHours * 60 * 60,
	})

	log.Info(ctx, action.Login, "user successfully logged in")
	writeJSON(w, http.StatusOK, map[string]string{"message": "login successful"})
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
