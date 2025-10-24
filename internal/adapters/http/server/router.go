package server

import (
	"ride-hail/internal/core/domain/types"
)

func (a *API) setupRoutes() {
	a.setupDefaultRoutes()

	switch a.cfg.Mode {
	case types.ModeAdmin:
	case types.ModeDAL:
	case types.ModeRide:
		a.setupRideRoutes()
	}

}
func (a *API) setupDefaultRoutes() {
	a.mux.HandleFunc("POST /registration", a.h.auth.Registration)
	a.mux.HandleFunc("POST /login", a.h.auth.Login)
}

func (a *API) setupRideRoutes() {
	a.mux.HandleFunc("/rides/", a.h.ride.CreateNewRide)
	a.mux.HandleFunc("/rides/{ride_id}/cancel", a.h.ride.CancelRide)
}
