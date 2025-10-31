package websocket

import (
	"encoding/json"
	"net/http"
	"ride-hail/internal/core/domain/action"
	"ride-hail/pkg/logger"
	"strings"
)

type PassengerWSHandler interface {
	PassengerWebSocketHandler(w http.ResponseWriter, r *http.Request)
}

type PassengerWebSocketHandler struct {
	manager *PassengerWebSocketManager
	log     *logger.Logger
}

func NewPassengerWebSocketHandler(manager *PassengerWebSocketManager, log *logger.Logger) *PassengerWebSocketHandler {
	return &PassengerWebSocketHandler{
		manager: manager,
		log:     log,
	}
}

func (ph *PassengerWebSocketHandler) PassengerWebSocketHandler(w http.ResponseWriter, r *http.Request) {
	log := ph.log.Func("PassengerWebSocketHandler.PassengerWebSocketHandler")
	ctx := r.Context()

	passengerId := getRideID(r)

	if passengerId == "" {
		log.Error(ctx, action.WSPassenger, "invalid id")
		writeJSON(w, http.StatusBadRequest, "invalid ride_id")
		return
	}

	log.Info(r.Context(), "ws_connection_attempt", "passenger attempting WebSocket connection")

	ph.manager.HandlePassengerConnection(w, r, passengerId)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func getRideID(r *http.Request) string {
	path := r.URL.Path
	parts := strings.Split(path, "/")
	return parts[3]
}
