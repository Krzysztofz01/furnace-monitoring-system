package domain

import "github.com/google/uuid"

type Measurement struct {
	Id                         uuid.UUID               `json:"id"`
	TemperatureSensorOne       TemperatureReading      `json:"temperatureSensorOne"`
	TemperatureSensorTwo       TemperatureReading      `json:"temperatureSensorTwo"`
	TemperatureSensorThree     TemperatureReading      `json:"temperatureSensorThree"`
	AirContaminationPercentage AirContaminationReading `json:"airContaminationPercentage"`
	TimestampUnix              int64                   `json:"timestampUnix"`
}

type TemperatureReading struct {
	Value     float64 `json:"value"`
	IsDefined bool    `json:"isDefined"`
}

type AirContaminationReading struct {
	Value     int64 `json:"value"`
	IsDefined bool  `json:"isDefined"`
}

type MeasurementRepository interface {
	GetLatestMeasurement() (*Measurement, error)
	GetAllMeasurements() ([]*Measurement, error)
	GetCurrentDayMeasurements() ([]*Measurement, error)
	InsertMeasurement(m *Measurement) error
}
