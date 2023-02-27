package server

import (
	"net/http"
	"time"

	"github.com/Krzysztofz01/furnace-monitoring-system/config"
	"github.com/Krzysztofz01/furnace-monitoring-system/db"
	"github.com/Krzysztofz01/furnace-monitoring-system/domain"
	"github.com/Krzysztofz01/furnace-monitoring-system/log"
	"github.com/Krzysztofz01/furnace-monitoring-system/protocol"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

const (
	timeOffsetInactivityCheckMinutes int     = 60
	timeOffsetErrorCountCheckMinutes int     = 2
	maxInactivitySeconds             float64 = 60
	maxErrorCount                    int     = 5
)

type WebsocketServer struct {
	dashboardHostPool *HostPool
	sensorHostPool    *HostPool

	sensorMeasurementChannel chan protocol.EventPayload
}

var Instance *WebsocketServer

func CreateWebSocketServer() error {
	Instance = new(WebsocketServer)
	Instance.dashboardHostPool = CreateHostPool()
	Instance.sensorHostPool = CreateHostPool()
	Instance.sensorMeasurementChannel = make(chan protocol.EventPayload)

	go Instance.handleSensorMeasurements()
	go Instance.handleHostInactivityCheck()
	go Instance.handleHostErrorCountCheck()

	return nil
}

var upgrader = websocket.Upgrader{
	CheckOrigin:     func(r *http.Request) bool { return true },
	WriteBufferSize: 1024,
	ReadBufferSize:  1024,
}

func (wss *WebsocketServer) UpgradeSensorHostConnection(c echo.Context) error {
	socket, err := upgrader.Upgrade(c.Response().Writer, c.Request(), nil)
	if err != nil {
		log.Instance.Debugf("Failed to upgrade the connection to websocket communication: %s\n", err)
		return nil
	}

	_, connectionPayloadBuffer, err := socket.ReadMessage()
	if err != nil {
		log.Instance.Debugf("Failed to retrieve the connection payload event message from the socket: %s\n", err)
		socket.Close()
		return nil
	}

	connectionEventPayload, err := protocol.ParseEventPayloadFromBuffer(connectionPayloadBuffer)
	if err != nil {
		log.Instance.Debugf("Failed to parse the connection payload event: %s\n", err)
		socket.Close()
		return nil
	} else {
		if connectionEventPayload.GetEventType() != protocol.SensorConnectedEvent {
			log.Instance.Debugf("The retrieved event payload is not of the expected type: %d\n", connectionEventPayload.GetEventType())
			socket.Close()
			return nil
		}
	}

	hostId := connectionEventPayload.GetHostId()
	if !wss.ValidateSensorHostId(hostId) {
		log.Instance.Debug("The host with the given id is not permited to be a sensor\n")
		socket.Close()
		return nil
	}

	if err := wss.sensorHostPool.InsertHost(hostId, socket); err != nil {
		log.Instance.Debugf("Failed to store the host connection: %s\n", err)
		socket.Close()
		return nil
	}

	log.Instance.Infof("Sensor host connection upgraded for host: %s with address: %s\n", hostId, socket.RemoteAddr().String())

	go wss.handleSensorHostConnection(hostId)
	return nil
}

func (wss *WebsocketServer) handleSensorHostConnection(hostId uuid.UUID) {
	defer func() {
		if deleted, err := wss.sensorHostPool.RemoveHost(hostId); !deleted {
			log.Instance.Debug("The sensor host has not been deleted, but it might be deleted previously\n")
			return
		} else if err != nil {
			log.Instance.Debugf("The sensor host has been deleted, but some errors occured: %s\n", err)
			return
		}

		log.Instance.Debug("The sensor host has been disposed and deleted\n")
	}()

	for {
		eventPayloadBuffer, err := wss.sensorHostPool.ReadFromHost(hostId)
		if err != nil {
			log.Instance.Debugf("Failed to retrieve the payload event message from the socket: %s\n", err)
			return
		}

		eventPayload, err := protocol.ParseEventPayloadFromBuffer(eventPayloadBuffer)
		if err != nil {
			log.Instance.Debugf("server: failed to parse the received event payload: %w", err)
			wss.sensorHostPool.BumpHostErrorCount(hostId)
			continue
		}

		switch eventPayload.GetEventType() {
		case protocol.SensorConnectedEvent:
			{
				log.Instance.Debug("Sensor connected event received on listening loop. Possible protocol error\n")
				wss.sensorHostPool.BumpHostErrorCount(hostId)
				continue
			}
		case protocol.SensorMeasurementEvent:
			{
				log.Instance.Debug("Sensor measurement event received. Pushing the payload to measurement channel\n")
				wss.sensorMeasurementChannel <- eventPayload
				continue
			}
		case protocol.SensorDisconnectedEvent:
			{
				log.Instance.Debug("Sensor disconnected event received. Breaking the listening loop\n")
				return
			}
		default:
			{
				log.Instance.Debug("Undefined event received on listening loop. Possible protocol error\n")
				wss.sensorHostPool.BumpHostErrorCount(hostId)
				continue
			}
		}
	}
}

func (wss *WebsocketServer) UpgradeDashboardHostConnection(c echo.Context) error {
	socket, err := upgrader.Upgrade(c.Response().Writer, c.Request(), nil)
	if err != nil {
		log.Instance.Debugf("Failed to upgrade the connection to websocket communication: %s\n", err)
		return nil
	}

	_, connectionPayloadBuffer, err := socket.ReadMessage()
	if err != nil {
		log.Instance.Debugf("Failed to retrieve the connection payload event message from the socket: %s\n", err)
		socket.Close()
		return nil
	}

	connectionEventPayload, err := protocol.ParseEventPayloadFromBuffer(connectionPayloadBuffer)
	if err != nil {
		log.Instance.Debugf("Failed to parse the connection payload event: %s\n", err)
		socket.Close()
		return nil
	} else {
		if connectionEventPayload.GetEventType() != protocol.DashboardConnectedEvent {
			log.Instance.Debugf("The retrieved event payload is not of the expected type: %d\n", connectionEventPayload.GetEventType())
			socket.Close()
			return nil
		}
	}

	hostId := connectionEventPayload.GetHostId()
	if err := wss.dashboardHostPool.InsertHost(hostId, socket); err != nil {
		log.Instance.Debugf("Failed to store the host connection: %s\n", err)
		socket.Close()
		return nil
	}

	log.Instance.Infof("Dashboard host connection upgraded for host: %s with address: %s", hostId, socket.RemoteAddr().String())

	go wss.handleDashboardHostConnection(hostId)
	return nil
}

func (wss *WebsocketServer) handleDashboardHostConnection(hostId uuid.UUID) {
	defer func() {
		if deleted, err := wss.dashboardHostPool.RemoveHost(hostId); !deleted {
			log.Instance.Debug("The dashboard host has not been deleted, but it might be deleted previously\n")
			return
		} else if err != nil {
			log.Instance.Debugf("The dashboard host has been deleted, but some errors occured: %s\n", err)
			return
		}

		log.Instance.Debug("The dashboard host has been disposed and deleted\n")
	}()

	for {
		eventPayloadBuffer, err := wss.dashboardHostPool.ReadFromHost(hostId)
		if err != nil {
			log.Instance.Debugf("Failed to retrieve the payload event message from the socket: %s\n", err)
			return
		}

		eventPayload, err := protocol.ParseEventPayloadFromBuffer(eventPayloadBuffer)
		if err != nil {
			log.Instance.Debugf("server: failed to parse the received event payload: %w", err)
			wss.dashboardHostPool.BumpHostErrorCount(hostId)
			continue
		}

		switch eventPayload.GetEventType() {
		case protocol.DashboardConnectedEvent:
			{
				log.Instance.Debug("Dashboard connected event received on listening loop. Possible protocol error\n")
				wss.dashboardHostPool.BumpHostErrorCount(hostId)
				continue
			}
		case protocol.DashboardDisconnectedEvent:
			{
				log.Instance.Debug("Dashboard disconnected event received. Breaking the listening loop\n")
				return
			}
		default:
			{
				log.Instance.Debug("Undefined event received on listening loop. Possible protocol error\n")
				wss.dashboardHostPool.BumpHostErrorCount(hostId)
				continue
			}
		}
	}
}

func (wss *WebsocketServer) handleSensorMeasurements() {
	for measurementPayload := range wss.sensorMeasurementChannel {
		measurement, err := domain.CreateMeasurementFromEventPayload(measurementPayload)
		if err != nil {
			log.Instance.Debugf("Failed to parse the measurement payload in order to create the domain measurement representation")
			continue
		}

		if err := db.InsertMeasurement(db.Instance, measurement); err != nil {
			log.Instance.Debugf("Failed to store the measurement in the database: %w", err)
		}

		hostIds := wss.dashboardHostPool.GetAllHostIds()
		log.Instance.Debugf("About to pass the measurement payload to: %d dashboard hosts", len(hostIds))

		for _, hostId := range hostIds {
			if err := wss.dashboardHostPool.SendToHost(hostId, measurementPayload); err != nil {
				log.Instance.Debugf("Failed to pass the measurement payload to dashboard host: %s", hostId)
			} else {
				log.Instance.Debugf("Successful passed the measurement payload to dashboard host: %s", hostId)
			}
		}
	}
}

func (wss *WebsocketServer) handleHostInactivityCheck() {
	for range time.Tick(time.Minute * time.Duration(timeOffsetInactivityCheckMinutes)) {
		dashboardHosts := wss.dashboardHostPool.GetAllHostIds()
		log.Instance.Debugf("About to check inactivity time for: %d dashboard hosts", len(dashboardHosts))

		for _, hostId := range dashboardHosts {
			seconds, err := wss.dashboardHostPool.GetHostSecondsSinceLastActivity(hostId)
			if err != nil {
				log.Instance.Debugf("Failed to perform inactivity time validation for host: %s", hostId)
				wss.dashboardHostPool.BumpHostErrorCount(hostId)
				continue
			}

			if seconds > maxInactivitySeconds {
				log.Instance.Debugf("Inactivity time of the host: %s exceeded, about to dispose the host.", hostId)

				if deleted, err := wss.dashboardHostPool.RemoveHost(hostId); !deleted {
					log.Instance.Debug("The dashboard host has not been deleted, but it might be deleted previously\n")
					continue
				} else if err != nil {
					log.Instance.Debugf("The dashboard host has been deleted, but some errors occured: %s\n", err)
					continue
				}
			}
		}

		sensorHosts := wss.sensorHostPool.GetAllHostIds()
		log.Instance.Debugf("About to check inactivity time for: %d sensor hosts", len(sensorHosts))

		for _, hostId := range sensorHosts {
			seconds, err := wss.sensorHostPool.GetHostSecondsSinceLastActivity(hostId)
			if err != nil {
				log.Instance.Debugf("Failed to perform inactivity time validation for host: %s", hostId)
				wss.sensorHostPool.BumpHostErrorCount(hostId)
				continue
			}

			if seconds > maxInactivitySeconds {
				log.Instance.Debugf("Inactivity time of the host: %s exceeded, about to dispose the host.", hostId)

				if deleted, err := wss.sensorHostPool.RemoveHost(hostId); !deleted {
					log.Instance.Debug("The sensor host has not been deleted, but it might be deleted previously\n")
					continue
				} else if err != nil {
					log.Instance.Debugf("The sensor host has been deleted, but some errors occured: %s\n", err)
					continue
				}
			}
		}
	}
}

func (wss *WebsocketServer) handleHostErrorCountCheck() {
	for range time.Tick(time.Minute * time.Duration(timeOffsetErrorCountCheckMinutes)) {
		dashboardHosts := wss.dashboardHostPool.GetAllHostIds()
		log.Instance.Debugf("About to check error count for: %d dashboard hosts", len(dashboardHosts))

		for _, hostId := range dashboardHosts {
			errorCount, err := wss.dashboardHostPool.GetHostErrorCount(hostId)
			if err != nil {
				log.Instance.Debugf("Failed to perform error count validation for host: %s", hostId)
				continue
			}

			if errorCount > maxErrorCount {
				log.Instance.Debugf("Error count of the host: %s exceeded, about to dispose the host.", hostId)

				if deleted, err := wss.dashboardHostPool.RemoveHost(hostId); !deleted {
					log.Instance.Debug("The dashboard host has not been deleted, but it might be deleted previously\n")
					continue
				} else if err != nil {
					log.Instance.Debugf("The dashboard host has been deleted, but some errors occured: %s\n", err)
					continue
				}
			}
		}

		sensorHosts := wss.sensorHostPool.GetAllHostIds()
		log.Instance.Debugf("About to check error count for: %d sensor hosts", len(sensorHosts))

		for _, hostId := range sensorHosts {
			errorCount, err := wss.sensorHostPool.GetHostErrorCount(hostId)
			if err != nil {
				log.Instance.Debugf("Failed to perform error count validation for host: %s", hostId)
				continue
			}

			if errorCount > maxErrorCount {
				log.Instance.Debugf("Error count of the host: %s exceeded, about to dispose the host.", hostId)

				if deleted, err := wss.sensorHostPool.RemoveHost(hostId); !deleted {
					log.Instance.Debug("The sensor host has not been deleted, but it might be deleted previously\n")
					continue
				} else if err != nil {
					log.Instance.Debugf("The sensor host has been deleted, but some errors occured: %s\n", err)
					continue
				}
			}
		}
	}
}

func (wss *WebsocketServer) ValidateSensorHostId(hostId uuid.UUID) bool {
	if config.Instance.SensorHostIds == nil {
		return false
	}

	for _, allowedHostId := range config.Instance.SensorHostIds {
		if hostId == uuid.MustParse(allowedHostId) {
			return true
		}
	}

	return false
}
