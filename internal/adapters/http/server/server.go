package server

import (
	"net/http"
	"ride-hail/config"
	"ride-hail/internal/adapters/http/handle"
	"ride-hail/internal/core/domain/types"
	"strconv"
)

type API struct {
	h    *handlers
	mux  *http.ServeMux
	addr int
	cfg  config.Config
}

type handlers struct {
	auth *handle.Handle
	ride *handle.RideHandle
}

func New(authH *handle.Handle, cfg config.Config) *API {
	h := &handlers{
		auth: authH,
	}

	api := &API{
		h:   h,
		cfg: cfg,
		mux: http.NewServeMux(),
	}
	api.initAddr()
	api.setupRoutes()
	return api
}

func (a *API) Run() error {
	return http.ListenAndServe(":"+strconv.Itoa(a.addr), a.Middleware(a.mux))
}

func (a *API) initAddr() {
	switch a.cfg.Mode {
	case types.ModeAdmin:
		a.addr = a.cfg.Services.AdminService
	case types.ModeDAL:
		a.addr = a.cfg.Services.DriverLocationService
	case types.ModeRide:
		a.addr = a.cfg.Services.RideService
	}
}
