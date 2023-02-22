package protocol

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	PayloadSeparatorSymbolRune   rune   = ';'
	PayloadSeparatorSymbolString string = string(PayloadSeparatorSymbolRune)
)

type EventType int

const (
	SensorConnectedEvent EventType = 1 + iota
	SensorDisconnectedEvent
	SensorMeasurementEvent
	DashboardConnectedEvent
	DashboardDisconnectedEvent
)

func GetEventTypeFromEventPayload(ep string) (EventType, error) {
	epParts := strings.Split(ep, PayloadSeparatorSymbolString)
	if len(epParts) < 2 {
		return 0, errors.New("protocol: invalid event payload format")
	}

	eventTypeValue, err := strconv.Atoi(epParts[0])
	if err != nil {
		return 0, errors.New("protocol: failed to parse the event type")
	}

	switch eventTypeValue {
	case 1, 2, 3, 4, 5:
		return EventType(eventTypeValue), nil
	default:
		return 0, errors.New("protocol: invalid event type identifier provided")
	}
}

func EventTypeToString(et EventType) string {
	return fmt.Sprintf("%d", et)
}
