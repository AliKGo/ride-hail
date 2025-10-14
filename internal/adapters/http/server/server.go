package server

import (
	"net/http"
	"ride-hail/config"
	"ride-hail/internal/core/ports"
	"strconv"
)

type API struct {
	h    ports.AuthHandler
	mux  *http.ServeMux
	addr int
	cfg  config.Config
}

func New(h ports.AuthHandler, cfg config.Config) *API {
	api := &API{
		h:    h,
		cfg:  cfg,
		addr: cfg.Services.RideService,
	}

	api.initRouterAuth()
	return api
}

func (a *API) Run() error {
	return http.ListenAndServe(":"+strconv.Itoa(a.addr), a.Middleware(a.mux))
}
