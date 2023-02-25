package server

import (
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// TODO: Implement aditional UUID and socket connection validation
// TODO: Perform tests to verify if the hostpool mutex control is enough.
type Host struct {
	id           uuid.UUID
	socket       *websocket.Conn
	errorCount   int
	lastActivity time.Time
}

func CreateHost(hostId uuid.UUID, socket *websocket.Conn) *Host {
	return &Host{
		id:           hostId,
		socket:       socket,
		errorCount:   0,
		lastActivity: time.Now(),
	}
}

func (h *Host) Send(buffer []byte) error {
	err := h.socket.WriteMessage(1, buffer)
	if err != nil {
		h.errorCount += 1
	} else {
		h.lastActivity = time.Now()
	}

	return err
}

func (h *Host) Read() ([]byte, error) {
	_, buffer, err := h.socket.ReadMessage()
	if err != nil {
		h.errorCount += 1
	} else {
		h.lastActivity = time.Now()
	}

	return buffer, err
}

func (h *Host) BumpErrorCount() bool {
	h.errorCount += 1
	return h.errorCount < maxErrorCount
}

func (h *Host) GetSecondsSinceLastActivity() float64 {
	return time.Since(h.lastActivity).Seconds()
}

func (h *Host) GetErrorCount() int {
	return h.errorCount
}

func (h *Host) Dispose() error {
	return h.socket.Close()
}
