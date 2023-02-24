package server

import (
	"errors"
	"fmt"
	"sync"

	"github.com/Krzysztofz01/furnace-monitoring-system/protocol"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// TODO: Optimize using RWMutex
type HostPool struct {
	hosts map[uuid.UUID]*Host
	mutex sync.Mutex
}

func CreateHostPool() *HostPool {
	return &HostPool{
		hosts: make(map[uuid.UUID]*Host),
		mutex: sync.Mutex{},
	}
}

func (hp *HostPool) InsertHost(hostId uuid.UUID, socket *websocket.Conn) error {
	if hostId == uuid.Nil {
		return errors.New("server: invalid unitialized uuid provided as host identifier")
	}

	if socket == nil {
		return errors.New("server: the socket instance is not initialized")
	}

	hp.mutex.Lock()
	defer hp.mutex.Unlock()

	if _, hostExists := hp.hosts[hostId]; hostExists {
		return errors.New("server: a host with the given identifier is already stored")
	}

	hp.hosts[hostId] = CreateHost(hostId, socket)
	return nil
}

func (hp *HostPool) RemoveHost(hostId uuid.UUID) (bool, error) {
	if hostId == uuid.Nil {
		return false, errors.New("server: invalid unitialized uuid provided as host identifier")
	}

	hp.mutex.Lock()
	defer hp.mutex.Unlock()

	host, hostExists := hp.hosts[hostId]
	if !hostExists {
		return false, errors.New("server: a host with the given identifier is not stored")
	}

	err := host.Dispose()
	delete(hp.hosts, hostId)

	return true, err
}

func (hp *HostPool) SendToHost(hostId uuid.UUID, payload protocol.EventPayload) error {
	if hostId == uuid.Nil {
		return errors.New("server: invalid unitialized uuid provided as host identifier")
	}

	hp.mutex.Lock()
	defer hp.mutex.Unlock()

	host, hostExists := hp.hosts[hostId]
	if !hostExists {
		return errors.New("server: a host with the given identifier is not stored")
	}

	// TODO: Implement payload to []byte conversion
	if err := host.Send([]byte{}); err != nil {
		return fmt.Errorf("server: failed to send the payload to given host: %w", err)
	} else {
		return nil
	}
}

func (hp *HostPool) ReadFromHost(hostId uuid.UUID) ([]byte, error) {
	if hostId == uuid.Nil {
		return nil, errors.New("server: invalid unitialized uuid provided as host identifier")
	}

	hp.mutex.Lock()
	defer hp.mutex.Unlock()

	host, hostExists := hp.hosts[hostId]
	if !hostExists {
		return nil, errors.New("server: a host with the given identifier is not stored")
	}

	return host.Read()
}

func (hp *HostPool) HasHostInactivityTimeExceeded(hostId uuid.UUID) (bool, error) {
	if hostId == uuid.Nil {
		return false, errors.New("server: invalid unitialized uuid provided as host identifier")
	}

	hp.mutex.Lock()
	defer hp.mutex.Unlock()

	host, hostExists := hp.hosts[hostId]
	if !hostExists {
		return false, errors.New("server: a host with the given identifier is not stored")
	}

	return host.HasInactivityTimeExceeded(), nil
}

func (hp *HostPool) HasHostErrorCountExceeded(hostId uuid.UUID) (bool, error) {
	if hostId == uuid.Nil {
		return false, errors.New("server: invalid unitialized uuid provided as host identifier")
	}

	hp.mutex.Lock()
	defer hp.mutex.Unlock()

	host, hostExists := hp.hosts[hostId]
	if !hostExists {
		return false, errors.New("server: a host with the given identifier is not stored")
	}

	return host.HasErrorCountExceeded(), nil
}

func (hp *HostPool) BumpHostErrorCount(hostId uuid.UUID) error {
	if hostId == uuid.Nil {
		return errors.New("server: invalid unitialized uuid provided as host identifier")
	}

	hp.mutex.Lock()
	defer hp.mutex.Unlock()

	host, hostExists := hp.hosts[hostId]
	if !hostExists {
		return errors.New("server: a host with the given identifier is not stored")
	}

	host.BumpErrorCount()
	return nil
}

func (hp *HostPool) GetAllHostIds() []uuid.UUID {
	hostIds := make([]uuid.UUID, 0, len(hp.hosts))
	for hostId := range hp.hosts {
		hostIds = append(hostIds, hostId)
	}

	return hostIds
}
