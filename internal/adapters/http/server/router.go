package server

import (
	"errors"
	"net/http"
	"ride-hail/internal/core/domain/types"
)

func (a *API) setupRoutes(mux *http.ServeMux) error {
	if err := a.setupDefaultRoutes(mux); err != nil {
		return err
	}

	switch a.cfg.Mode {
	case types.ModeAdmin:
	case types.ModeDAL:
	case types.ModeRide:
		if err := a.setupRideRoutes(mux); err != nil {
			return err
		}
	}
	return nil
}

func (a *API) setupDefaultRoutes(mux *http.ServeMux) error {
	if a.h.auth == nil {
		return errors.New("authorization service is request")
	}
	mux.HandleFunc("POST /registration", a.h.auth.Registration)
	mux.HandleFunc("POST /login", a.h.auth.Login)
	return nil
}

func (a *API) setupRideRoutes(mux *http.ServeMux) error {
	if a.h.ride == nil {
		return errors.New("ride service is required")
	}
	mux.HandleFunc("/rides/", a.jwtMiddleware(a.h.ride.CreateNewRide))
	mux.HandleFunc("/rides/{ride_id}/cancel", a.jwtMiddleware(a.h.ride.CancelRide))
	return nil
}
