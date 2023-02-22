package server

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Host struct {
	id           uuid.UUID
	socket       *websocket.Conn
	lastActivity time.Time
	mu           sync.Mutex
}

// FIXME: Socket nil check
func CreateHost(hostId uuid.UUID, socket *websocket.Conn) *Host {
	host := new(Host)
	host.id = hostId
	host.socket = socket
	host.lastActivity = time.Now()
	host.mu = sync.Mutex{}

	return host
}

func (h *Host) Send(e []byte) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if err := h.socket.WriteMessage(1, e); err != nil {
		return err
	} else {
		h.lastActivity = time.Now()
		return nil
	}
}

func (h *Host) SendAsJson(v interface{}) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if err := h.socket.WriteJSON(v); err != nil {
		return err
	} else {
		h.lastActivity = time.Now()
		return nil
	}
}

func (h *Host) BumpActivity() {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.lastActivity = time.Now()
}

func (h *Host) GetSecondsSinceLastActivity() float64 {
	h.mu.Lock()
	defer h.mu.Unlock()

	return time.Since(h.lastActivity).Seconds()
}

func (h *Host) Dispose() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	return h.socket.Close()
}
