package server

import "net/http"

func (a *API) initRouterAuth() {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /registration", a.h.Registration)
	mux.HandleFunc("POST /login", a.h.Login)

	a.mux = mux
}
