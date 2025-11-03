package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"ride-hail/internal/core/domain/action"
	"ride-hail/pkg/logger"
	"sync"
	"time"
)

type PassengerWebSocketManager struct {
	connections map[string]*Passenger
	mu          sync.RWMutex
	log         *logger.Logger
	ctx         context.Context
	cancel      context.CancelFunc
	wg          sync.WaitGroup
}

type Passenger struct {
	id            string
	conn          *websocket.Conn
	authenticated bool
	authTimeout   time.Time
	send          chan []byte
	lastPing      time.Time
	cancel        context.CancelFunc
}

type PassengerWSMessage struct {
	Type  string      `json:"type"`
	Token string      `json:"token,omitempty"`
	Data  interface{} `json:"data,omitempty"`
}

func NewPassengerWebSocketManager(ctx context.Context, log *logger.Logger) *PassengerWebSocketManager {
	ctx, cancel := context.WithCancel(ctx)
	return &PassengerWebSocketManager{
		connections: make(map[string]*Passenger),
		log:         log,
		ctx:         ctx,
		cancel:      cancel,
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (m *PassengerWebSocketManager) HandlePassengerConnection(w http.ResponseWriter, r *http.Request, passengerID string) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		m.log.Func("HandlePassengerConnection").Error(r.Context(), action.WSPassenger, "upgrade failed", "error", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to upgrade"})
		return
	}

	ctx, cancel := context.WithCancel(m.ctx)
	passenger := &Passenger{
		id:            passengerID,
		conn:          conn,
		authenticated: false,
		send:          make(chan []byte, 10),
		authTimeout:   time.Now().Add(5 * time.Second),
		cancel:        cancel,
	}

	m.mu.Lock()
	m.connections[passenger.id] = passenger
	m.mu.Unlock()

	m.log.Func("HandlePassengerConnection").Info(ctx, action.WSPassenger, "new passenger connected", "id", passenger.id)

	m.wg.Add(2)
	go m.writePump(ctx, passenger)
	go m.readPump(ctx, passenger)

	go func() {
		<-ctx.Done()
		conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "server shutdown"))
		conn.Close()
	}()
}

func (m *PassengerWebSocketManager) readPump(ctx context.Context, p *Passenger) {
	defer func() {
		p.cancel()
		m.removePassenger(p.id)
		m.wg.Done()
	}()

	log := m.log.Func("PassengerWebSocketManager.readPump")
	p.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	p.conn.SetPongHandler(func(string) error {
		p.lastPing = time.Now()
		p.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		var msg PassengerWSMessage
		if err := p.conn.ReadJSON(&msg); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Error(ctx, action.WSPassenger, "unexpected close", "error", err)
			}
			return
		}

		if !p.authenticated {
			if msg.Type != "auth" {
				log.Warn(ctx, action.WSPassenger, "unauthenticated message", "type", msg.Type)
				return
			}
			if time.Now().After(p.authTimeout) {
				log.Warn(ctx, action.WSPassenger, "auth timeout")
				return
			}
			m.handleAuth(ctx, p, msg)
			continue
		}

		m.handleMessage(ctx, p, msg)
	}
}

func (m *PassengerWebSocketManager) writePump(ctx context.Context, p *Passenger) {
	defer func() {
		p.cancel()
		m.removePassenger(p.id)
		m.wg.Done()
	}()

	log := m.log.Func("PassengerWebSocketManager.writePump")
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Debug(ctx, action.WSPassenger, "context done -> closing writePump")
			return
		case msg, ok := <-p.send:
			if !ok {
				log.Debug(ctx, action.WSPassenger, "send channel closed")
				return
			}
			p.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := p.conn.WriteJSON(msg); err != nil {
				log.Error(ctx, action.WSPassenger, "write error", "error", err)
				return
			}
		case <-ticker.C:
			p.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := p.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Error(ctx, action.WSPassenger, "ping error", "error", err)
				return
			}
		}
	}
}

func (m *PassengerWebSocketManager) handleMessage(ctx context.Context, p *Passenger, msg PassengerWSMessage) {
	log := m.log.Func("handleMessage")
	switch msg.Type {
	default:
		log.Warn(ctx, action.WSPassenger, "unknown message", "type", msg.Type)
	}
}

func (m *PassengerWebSocketManager) handleAuth(ctx context.Context, p *Passenger, msg PassengerWSMessage) {
	log := m.log.Func("handleAuth")

	ctxToken := logger.GetToken(ctx)
	if ctxToken == "" {
		p.send <- m.marshalMessage(PassengerWSMessage{Type: "auth_error", Data: "missing token"})
		return
	}

	if ctxToken != msg.Token {
		p.send <- m.marshalMessage(PassengerWSMessage{Type: "auth_error", Data: "invalid token"})
		return
	}

	p.authenticated = true
	p.authTimeout = time.Time{}
	log.Info(ctx, action.WSPassenger, "authenticated", "id", p.id)

	resp := PassengerWSMessage{Type: "auth_success"}
	data, _ := json.Marshal(resp)
	select {
	case p.send <- data:
	default:
		log.Warn(ctx, action.WSPassenger, "send channel full -> closing connection")
		p.cancel()
	}
}

func (m *PassengerWebSocketManager) removePassenger(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if p, ok := m.connections[id]; ok {
		close(p.send)
		delete(m.connections, id)
	}
}

func (m *PassengerWebSocketManager) Shutdown() {
	m.log.Func("Shutdown").Info(context.Background(), action.WSPassenger, "closing all WS connections")
	m.cancel()

	m.mu.Lock()
	for _, p := range m.connections {
		p.cancel()
		p.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "server shutdown"))
		p.conn.Close()
	}
	m.mu.Unlock()

	m.wg.Wait()
}

func (m *PassengerWebSocketManager) marshalMessage(msg interface{}) []byte {
	data, _ := json.Marshal(msg)
	return data
}

func (m *PassengerWebSocketManager) SendRide(ctx context.Context, passengerID string, data []byte) error {
	m.mu.RLock()
	conn, exists := m.connections[passengerID]
	m.mu.RUnlock()

	if !exists || !conn.authenticated {
		return fmt.Errorf("passenger %s not connected", passengerID)
	}

	select {
	case <-ctx.Done():
		return fmt.Errorf("timeout sending message to passenger %s", passengerID)
	case conn.send <- data:
		return nil
	case <-time.After(5 * time.Second):
		return fmt.Errorf("timeout sending message to passenger %s", passengerID)
	}
}
