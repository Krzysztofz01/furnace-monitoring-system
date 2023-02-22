package server

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/Krzysztofz01/furnace-monitoring-system/db"
	"github.com/Krzysztofz01/furnace-monitoring-system/domain"
	"github.com/Krzysztofz01/furnace-monitoring-system/protocol"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type WebsocketServer struct {
	sensorHosts         map[uuid.UUID]*Host
	mutexSensorHosts    sync.Mutex
	dashboardHosts      map[uuid.UUID]*Host
	mutexDashboardHosts sync.Mutex

	sensorMeasurementChannel chan string
}

var Instance *WebsocketServer

func CreateWebSocketServer() error {
	Instance = new(WebsocketServer)
	Instance.sensorHosts = make(map[uuid.UUID]*Host)
	Instance.mutexSensorHosts = sync.Mutex{}
	Instance.dashboardHosts = make(map[uuid.UUID]*Host)
	Instance.mutexDashboardHosts = sync.Mutex{}
	Instance.sensorMeasurementChannel = make(chan string)

	go Instance.handleSensorMeasurements()
	go Instance.handleHostDisposal()

	return nil
}

var upgrader = websocket.Upgrader{
	CheckOrigin:     func(r *http.Request) bool { return true },
	WriteBufferSize: 1024,
	ReadBufferSize:  1024,
}

func (wss *WebsocketServer) UpgradeSensorClientConnection(r *http.Request, w http.ResponseWriter) {
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(fmt.Errorf("server: failed to upgrade the connection to websocket: %w", err))
		return
	}

	// TODO: Solution for hostid preservation
	hostId := uuid.New()
	connectionEventPayload, err := protocol.CreateDashboardConnectedEventPayload(hostId)
	if err != nil {
		fmt.Println(fmt.Errorf("server: failed to create the sensor connected event payload: %w", err))
		return
	}

	if err := socket.WriteMessage(1, connectionEventPayload); err != nil {
		fmt.Println(fmt.Errorf("server: failed to send the connection event payload: %w", err))
		return
	}

	defer func() {
		// NOTE: Defering the handling of host disconnection
		wss.mutexSensorHosts.Lock()
		defer wss.mutexSensorHosts.Unlock()

		host, isStored := wss.sensorHosts[hostId]
		if isStored {
			if err := host.Dispose(); err != nil {
				fmt.Println(fmt.Errorf("server: problem occured while disposing the host connection: %w", err))
			}
			delete(wss.sensorHosts, hostId)
		}
	}()

	for {
		_, messagePayload, err := socket.ReadMessage()
		if err != nil {
			// TODO: Handle this situation. Introduce failure count, to dispose host that generates a lot of failures
			fmt.Println(fmt.Errorf("server: failed to retrievie the message from the socket: %w", err))
			continue
		}

		// TODO: Should we lock this with RWMutex
		wss.mutexSensorHosts.Lock()
		host, isStored := wss.sensorHosts[hostId]
		if isStored {
			host.BumpActivity()
		} else {
			wss.sensorHosts[hostId] = CreateHost(hostId, socket)
		}
		wss.mutexSensorHosts.Unlock()

		eventType, err := protocol.GetEventTypeFromEventPayload(string(messagePayload))
		if err != nil {
			// TODO: Handle this situation. Introduce failure count, to dispose host that generates a lot of failures
			fmt.Println(fmt.Errorf("server: failed to retrieve the event type from the event payload: %w", err))
			continue
		}

		switch eventType {
		case protocol.SensorMeasurementEvent:
			{
				wss.sensorMeasurementChannel <- string(messagePayload)
			}

		case protocol.SensorDisconnectedEvent:
			{
				wss.mutexSensorHosts.Lock()
				host := wss.sensorHosts[hostId]
				if err := host.Dispose(); err != nil {
					fmt.Println(fmt.Errorf("server: problem occured while disposing the host connection: %w", err))
				}

				delete(wss.sensorHosts, hostId)
				wss.mutexSensorHosts.Unlock()
			}
		default:
			{
				// TODO: Handle this situation. Log and continue
				fmt.Println(fmt.Errorf("server: invalid event type payload provided: %w", err))
				continue
			}
		}
	}
}

func (wss *WebsocketServer) UpgradeDashboardHostConnection(r *http.Request, w http.ResponseWriter) {
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(fmt.Errorf("server: failed to upgrade the connection to websocket: %w", err))
		return
	}

	hostId := uuid.New()
	connectionEventPayload, err := protocol.CreateDashboardConnectedEventPayload(hostId)
	if err != nil {
		fmt.Println(fmt.Errorf("server: failed to crate the dashboard connected event payload: %w", err))
		return
	}

	if err := socket.WriteMessage(1, connectionEventPayload); err != nil {
		fmt.Println(fmt.Errorf("server: failed to send the connection event payload: %w", err))
		return
	}

	defer func() {
		// NOTE: Defering the handling of host disconnection
		wss.mutexDashboardHosts.Lock()
		defer wss.mutexDashboardHosts.Unlock()

		host, isStored := wss.dashboardHosts[hostId]
		if isStored {
			if err := host.Dispose(); err != nil {
				fmt.Println(fmt.Errorf("server: problem occured while disposing the host connection: %w", err))
			}
			delete(wss.dashboardHosts, hostId)
		}
	}()

	for {
		_, messagePayload, err := socket.ReadMessage()
		if err != nil {
			// TODO: Handle this situation. Introduce failure count, to dispose host that generates a lot of failures
			fmt.Println(fmt.Errorf("server: failed to retrievie the message from the socket: %w", err))
			continue
		}

		// TODO: Should we lock this with RWMutex
		wss.mutexDashboardHosts.Lock()
		host, isStored := wss.dashboardHosts[hostId]
		if isStored {
			host.BumpActivity()
		} else {
			wss.dashboardHosts[hostId] = CreateHost(hostId, socket)
		}
		wss.mutexDashboardHosts.Unlock()

		eventType, err := protocol.GetEventTypeFromEventPayload(string(messagePayload))
		if err != nil {
			// TODO: Handle this situation. Introduce failure count, to dispose host that generates a lot of failures
			fmt.Println(fmt.Errorf("server: failed to retrieve the event type from the event payload: %w", err))
			continue
		}

		switch eventType {
		case protocol.DashboardDisconnectedEvent:
			{
				wss.mutexDashboardHosts.Lock()
				host := wss.dashboardHosts[hostId]
				if err := host.Dispose(); err != nil {
					fmt.Println(fmt.Errorf("server: problem occured while disposing the host connection: %w", err))
				}

				delete(wss.dashboardHosts, hostId)
				wss.mutexDashboardHosts.Unlock()
			}
		default:
			{
				// TODO: Handle this situation. Log and continue
				fmt.Println(fmt.Errorf("server: invalid event type payload provided: %w", err))
				continue
			}
		}
	}
}

func (wss *WebsocketServer) handleSensorMeasurements() {
	for measurementPayload := range wss.sensorMeasurementChannel {
		measurement, err := domain.CreateMeasurementFromEventPayload(string(measurementPayload))
		if err != nil {
			// TODO: Handle this situation. Drop the measurement and log
			fmt.Println(fmt.Errorf("server: failed to create a measurement instance: %w", err))
			continue
		}

		db.InsertMeasurement(db.Instance, measurement)

		// TODO: Should we lock this with RWMutex
		for _, host := range wss.dashboardHosts {
			// TODO: Implement the sending of the measurements
			if err := host.Send([]byte{}); err != nil {
				// TODO: Handle this situation. The failure can indicate that there is no connection and the host can be removed
				fmt.Println(fmt.Errorf("server: failed to create a measurement instance: %w", err))
				continue
			}
		}
	}
}

func (wss *WebsocketServer) handleHostDisposal() {
	// TODO: Fine-tune the time of this task
	for range time.Tick(time.Second * 60) {
		for hostId, host := range wss.dashboardHosts {
			// TODO: Fine-tune the time-to-life of the hosts
			if host.GetSecondsSinceLastActivity() < 30 {
				wss.mutexDashboardHosts.Lock()
				if err := host.Dispose(); err != nil {
					fmt.Println(fmt.Errorf("server: problem occured while disposing the host connection: %w", err))
				}

				delete(wss.dashboardHosts, hostId)
				wss.mutexDashboardHosts.Unlock()
			}
		}

		for hostId, host := range wss.sensorHosts {
			// TODO: Fine-tune the time-to-life of the hosts
			if host.GetSecondsSinceLastActivity() < 30 {
				wss.mutexSensorHosts.Lock()
				if err := host.Dispose(); err != nil {
					fmt.Println(fmt.Errorf("server: problem occured while disposing the host connection: %w", err))
				}

				delete(wss.sensorHosts, hostId)
				wss.mutexSensorHosts.Unlock()
			}
		}
	}
}
