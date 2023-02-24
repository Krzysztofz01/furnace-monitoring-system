package protocol

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

const (
	PayloadSeparatorSymbolRune   rune   = ';'
	PayloadSeparatorSymbolString string = string(PayloadSeparatorSymbolRune)
)

type EventPayload interface {
	GetEventType() EventType
	GetHostId() uuid.UUID
	GetBuffer() []byte
	GetSize() int
	GetFieldString(index int) (string, error)
	GetFieldInt(index int) (int, error)
	GetFieldFloat64(index int) (float64, error)
}

type eventPayload struct {
	eventType    EventType
	hostId       uuid.UUID
	payloadParts []string
}

func ParseEventPayloadFromBuffer(epb []byte) (EventPayload, error) {
	ep := new(eventPayload)
	epParts := strings.Split(string(epb), PayloadSeparatorSymbolString)
	if len(epParts) < 2 {
		return nil, errors.New("protocol: invalid event payload format")
	}

	ep.payloadParts = make([]string, 0)
	ep.payloadParts = append(ep.payloadParts, epParts...)

	eventTypeValue, err := strconv.Atoi(epParts[0])
	if err != nil {
		return nil, errors.New("protocol: failed to parse the event type")
	} else {
		switch eventTypeValue {
		case 1, 2, 3, 4, 5:
			ep.eventType = EventType(eventTypeValue)
		default:
			return nil, errors.New("protocol: invalid event type identifier provided")
		}
	}

	hostId, err := uuid.Parse(epParts[1])
	if err != nil {
		return nil, errors.New("protocol: failed to parse the provided host id")
	} else {
		ep.hostId = hostId
	}

	return ep, nil
}

// TODO: DRY-out the implementation...
func ParseEventPayloadFromValueSlice(values []string) (EventPayload, error) {
	lastValueIndex := len(values) - 1
	bufferBuilder := strings.Builder{}

	for index, value := range values {
		bufferBuilder.WriteString(value)
		if index != lastValueIndex {
			bufferBuilder.WriteRune(PayloadSeparatorSymbolRune)
		}
	}

	buffer := []byte(bufferBuilder.String())
	return ParseEventPayloadFromBuffer(buffer)
}

func (ep *eventPayload) GetEventType() EventType {
	return ep.eventType
}

func (ep *eventPayload) GetHostId() uuid.UUID {
	return ep.hostId
}

func (ep *eventPayload) GetSize() int {
	return len(ep.payloadParts)
}

func (ep *eventPayload) GetFieldString(index int) (string, error) {
	if index >= ep.GetSize() {
		return "", errors.New("protocol: provided payload part index is out of bound")
	}

	return ep.payloadParts[index], nil
}

func (ep *eventPayload) GetFieldInt(index int) (int, error) {
	if index >= ep.GetSize() {
		return 0, errors.New("protocol: provided payload part index is out of bound")
	}

	value, err := strconv.Atoi(ep.payloadParts[index])
	if err != nil {
		return 0, fmt.Errorf("protocol: failed to parse the payload part value to int: %w", err)
	}

	return value, nil
}

func (ep *eventPayload) GetFieldFloat64(index int) (float64, error) {
	if index >= ep.GetSize() {
		return 0, errors.New("protocol: provided payload part index is out of bound")
	}

	value, err := strconv.ParseFloat(ep.payloadParts[index], 64)
	if err != nil {
		return 0, fmt.Errorf("protocol: failed to parse the payload part value to float64: %w", err)
	}

	return value, nil
}

func (ep *eventPayload) GetBuffer() []byte {
	lastPartIndex := len(ep.payloadParts) - 1
	bufferBuilder := strings.Builder{}

	for index, part := range ep.payloadParts {
		bufferBuilder.WriteString(part)
		if index != lastPartIndex {
			bufferBuilder.WriteRune(PayloadSeparatorSymbolRune)
		}
	}

	return []byte(bufferBuilder.String())
}
