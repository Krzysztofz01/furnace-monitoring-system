package domain

import "github.com/google/uuid"

type Measurement struct {
	Id                         uuid.UUID               `json:"id"`
	TemperatureSensorOne       TemperatureReading      `json:"temperatureSensorOne"`
	TemperatureSensorTwo       TemperatureReading      `json:"temperatureSensorTwo"`
	TemperatureSensorThree     TemperatureReading      `json:"temperatureSensorThree"`
	AirContaminationPercentage AirContaminationReading `json:"airContaminationPercentage"`
}

type TemperatureReading struct {
	value     float64 `json:"value"`
	isDefined bool    `json:"isDefined"`
}

type AirContaminationReading struct {
	value     int  `json:"value"`
	isDefined bool `json:"isDefined"`
}

type MeasurementRepository interface {
	GetLatesstMeasurement() (*Measurement, error)
	GetAllMeasurements() ([]*Measurement, error)
	GetCurrentDayMeasurements() ([]*Measurement, error)
	InsertMeasurement(m *Measurement) error
}
