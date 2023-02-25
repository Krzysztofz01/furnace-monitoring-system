package server

import (
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Host struct {
	id           uuid.UUID
	socket       *websocket.Conn
	errorCount   int
	lastActivity time.Time
	// mu           sync.Mutex
}

// TODO: Socket nil check
func CreateHost(hostId uuid.UUID, socket *websocket.Conn) *Host {
	return &Host{
		id:           hostId,
		socket:       socket,
		errorCount:   0,
		lastActivity: time.Now(),
		// mu:           sync.Mutex{},
	}
}

func (h *Host) Send(buffer []byte) error {
	// h.mu.Lock()
	// defer h.mu.Unlock()

	err := h.socket.WriteMessage(1, buffer)
	if err != nil {
		h.errorCount += 1
	} else {
		h.lastActivity = time.Now()
	}

	return err
}

func (h *Host) Read() ([]byte, error) {
	// h.mu.Lock()
	// defer h.mu.Unlock()

	_, buffer, err := h.socket.ReadMessage()
	if err != nil {
		h.errorCount += 1
	} else {
		h.lastActivity = time.Now()
	}

	return buffer, err
}

func (h *Host) BumpErrorCount() bool {
	// h.mu.Lock()
	// defer h.mu.Unlock()

	h.errorCount += 1
	return h.errorCount < maxErrorCount
}

// func (h *Host) HasInactivityTimeExceeded() bool {
// 	// h.mu.Lock()
// 	// defer h.mu.Unlock()

// 	return time.Since(h.lastActivity).Seconds() > float64(maxInactivitySeconds)
// }

func (h *Host) GetSecondsSinceLastActivity() float64 {
	// h.mu.Lock()
	// defer h.mu.Unlock()

	return time.Since(h.lastActivity).Seconds()
}

func (h *Host) HasErrorCountExceeded() bool {
	// h.mu.Lock()
	// defer h.mu.Unlock()

	return h.errorCount > maxErrorCount
}

func (h *Host) GetErrorCount() int {
	// h.mu.Lock()
	// defer h.mu.Unlock()

	return h.errorCount
}

func (h *Host) Dispose() error {
	// h.mu.Lock()
	// defer h.mu.Unlock()

	return h.socket.Close()
}
