package server

import "net/http"

func (a *API) initRouterAuth() {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /registration", a.h.Registration)
	mux.HandleFunc("POST /login", a.h.Login)

	mux.HandleFunc("/health", func(writer http.ResponseWriter, request *http.Request) {

	})
	a.mux = mux
}
