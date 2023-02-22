package domain

import (
	"errors"
	"fmt"
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

func CreateMeasurementFromEventPayload(ep protocol.EventPayload) (*Measurement, error) {
	if ep.GetEventType() != protocol.SensorMeasurementEvent {
		return nil, errors.New("domain: failed to create the measurement due to invalid event payload format")
	}

	var err error = nil
	measurement := new(Measurement)
	measurement.Id = uuid.New()
	measurement.HostId = ep.GetHostId()
	measurement.TimestampUnix = time.Now().Unix()

	if measurement.TemperatureChannelOne, err = ep.GetFieldFloat64(2); err != nil {
		return nil, fmt.Errorf("domain: failed to parse the channel one temperature value: %w", err)
	}

	if measurement.TemperatureChannelTwo, err = ep.GetFieldFloat64(3); err != nil {
		return nil, fmt.Errorf("domain: failed to parse the channel two temperature value: %w", err)
	}

	if measurement.TemperatureChannelThree, err = ep.GetFieldFloat64(4); err != nil {
		return nil, fmt.Errorf("domain: failed to parse the channel three temperature value: %w", err)
	}

	if measurement.AirContaminationPercentage, err = ep.GetFieldFloat64(5); err != nil {
		return nil, fmt.Errorf("domain: failed to parse the air contamination value: %w", err)
	}

	return measurement, nil
}
