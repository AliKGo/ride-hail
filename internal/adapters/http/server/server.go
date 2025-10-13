package server

import (
	"net/http"
	"ride-hail/config"
	"ride-hail/internal/adapters/http/handle"
	"strconv"
)

type API struct {
	h    *handle.Handle
	mux  *http.ServeMux
	addr int
}

func New(h *handle.Handle, cfg config.Config) *API {
	api := &API{
		h:    h,
		addr: cfg.Services.RideService,
	}

	api.initRouterAuth()
	return api
}

func (a *API) Run() error {
	return http.ListenAndServe(":"+strconv.Itoa(a.addr), a.Middleware(a.mux))
}
