package protocol

type EventType int

const (
	SensorConnectedEvent EventType = 1 + iota
	SensorDisconnectedEvent
	SensorMeasurementEvent
	DashboardConnectedEvent
	DashboardDisconnectedEvent
)
