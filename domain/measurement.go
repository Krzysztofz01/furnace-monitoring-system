package domain

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Krzysztofz01/furnace-monitoring-system/protocol"
	"github.com/google/uuid"
)

type Measurement struct {
	Id                         uuid.UUID
	HostId                     uuid.UUID
	TemperatureChannelOne      float64
	TemperatureChannelTwo      float64
	TemperatureChannelThree    float64
	AirContaminationPercentage float64
	TimestampUnix              int64
}

// The structure of the measurement event
// event_type;device_id;temp1;temp2;temp3;air_contamination

func CreateMeasurementFromEventPayload(ep string) (*Measurement, error) {
	payloadParts := strings.Split(ep, ";")
	if len(payloadParts) != 6 {
		return nil, errors.New("domain: failed to create the measurement due to invalid event payload format")
	}

	eventType, err := protocol.GetEventTypeFromEventPayload(ep)
	if err != nil {
		return nil, fmt.Errorf("domain: failed to parse the event payload type: %w", err)
	}

	if eventType != protocol.SensorMeasurementEvent {
		return nil, errors.New("domain: invalid event payload type provided")
	}

	hostId, err := protocol.GetHostIdFromEventPayload(ep)
	if err != nil {
		return nil, fmt.Errorf("domain: failed to parse the host id: %w", err)
	}

	measurement := new(Measurement)
	measurement.Id = uuid.New()
	measurement.HostId = hostId
	measurement.TimestampUnix = time.Now().Unix()

	if measurement.TemperatureChannelOne, err = strconv.ParseFloat(payloadParts[2], 64); err != nil {
		return nil, fmt.Errorf("domain: failed to parse the channel one temperature value: %w", err)
	}

	if measurement.TemperatureChannelTwo, err = strconv.ParseFloat(payloadParts[3], 64); err != nil {
		return nil, fmt.Errorf("domain: failed to parse the channel two temperature value: %w", err)
	}

	if measurement.TemperatureChannelThree, err = strconv.ParseFloat(payloadParts[4], 64); err != nil {
		return nil, fmt.Errorf("domain: failed to parse the channel three temperature value: %w", err)
	}

	if measurement.AirContaminationPercentage, err = strconv.ParseFloat(payloadParts[5], 64); err != nil {
		return nil, fmt.Errorf("domain: failed to parse the air contamination value: %w", err)
	}

	return measurement, nil
}
