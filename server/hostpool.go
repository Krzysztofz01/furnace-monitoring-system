package server

import (
	"errors"
	"fmt"
	"sync"

	"github.com/Krzysztofz01/furnace-monitoring-system/protocol"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type HostPool struct {
	hosts map[uuid.UUID]*Host
	mutex sync.RWMutex
}

func CreateHostPool() *HostPool {
	return &HostPool{
		hosts: make(map[uuid.UUID]*Host),
		mutex: sync.RWMutex{},
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

	hp.mutex.RLock()
	defer hp.mutex.RUnlock()

	host, hostExists := hp.hosts[hostId]
	if !hostExists {
		return errors.New("server: a host with the given identifier is not stored")
	}

	if err := host.Send(payload.GetBuffer()); err != nil {
		return fmt.Errorf("server: failed to send the payload to given host: %w", err)
	} else {
		return nil
	}
}

// FIXME: This method is using a very dirty workaround to prevent locking on while listeing to incoming socket messages.
// In case of dashboard hosts it is possible to just hold the connection without listening for incoming traffic, but this
// wont work for the sensor hosts. The current approach is unlocking the mutex after host is retrieved from the map,
// and "is hoping" that the current host still exists in the pull while listening for network traffic. Another solution is
// to implement a host-specific mutex used only to control the read operation
func (hp *HostPool) ReadFromHost(hostId uuid.UUID) ([]byte, error) {
	if hostId == uuid.Nil {
		return nil, errors.New("server: invalid unitialized uuid provided as host identifier")
	}

	hp.mutex.RLock()

	host, hostExists := hp.hosts[hostId]
	if !hostExists {
		hp.mutex.RUnlock()
		return nil, errors.New("server: a host with the given identifier is not stored")
	}

	hp.mutex.RUnlock()

	var resultBuffer []byte
	var err error
	defer func() {
		if err := recover(); err != nil {
			resultBuffer = nil
			err = fmt.Errorf("server: most likely the ReadFromHost workaround casued a panic: %s", err)
		}
	}()

	resultBuffer, err = host.Read()
	return resultBuffer, err
}

func (hp *HostPool) GetHostSecondsSinceLastActivity(hostId uuid.UUID) (float64, error) {
	if hostId == uuid.Nil {
		return 0, errors.New("server: invalid unitialized uuid provided as host identifier")
	}

	hp.mutex.RLock()
	defer hp.mutex.RUnlock()

	host, hostExists := hp.hosts[hostId]
	if !hostExists {
		return 0, errors.New("server: a host with the given identifier is not stored")
	}

	return host.GetSecondsSinceLastActivity(), nil
}

func (hp *HostPool) GetHostErrorCount(hostId uuid.UUID) (int, error) {
	if hostId == uuid.Nil {
		return 0, errors.New("server: invalid unitialized uuid provided as host identifier")
	}

	hp.mutex.RLock()
	defer hp.mutex.RUnlock()

	host, hostExists := hp.hosts[hostId]
	if !hostExists {
		return 0, errors.New("server: a host with the given identifier is not stored")
	}

	return host.GetErrorCount(), nil
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
	hp.mutex.RLock()
	defer hp.mutex.RUnlock()

	hostIds := make([]uuid.UUID, 0)
	for hostId := range hp.hosts {
		hostIds = append(hostIds, hostId)
	}

	return hostIds
}
