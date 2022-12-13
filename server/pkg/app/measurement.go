package app

import (
	"errors"
	"fmt"
	"furnace-monitoring-system-server/pkg/domain"
)

type MeasurementService struct {
	measurementRepository domain.MeasurementRepository
}

func CreateMeasurementService(measurementRepository domain.MeasurementRepository) (*MeasurementService, error) {
	if measurementRepository == nil {
		return nil, errors.New("MeasurementService: Provided measurement repository instance is nil")
	}

	return &MeasurementService{
		measurementRepository: measurementRepository,
	}, nil
}

func (ms *MeasurementService) InsertMeasurement(encodedMeasurement string) error {
	// TODO: Implement the decoding algorithm
	return errors.New("MeasurementService: Not implemented")
}

func (ms *MeasurementService) GetLatestMeasurement() (*domain.Measurement, error) {
	measurement, err := ms.measurementRepository.GetLatestMeasurement()
	if err != nil {
		return nil, fmt.Errorf("MeasurementService: Failed to retrieve latest measurement: %w", err)
	}

	return measurement, nil
}

func (ms *MeasurementService) GetAllMeasurements() ([]*domain.Measurement, error) {
	measurements, err := ms.measurementRepository.GetAllMeasurements()
	if err != nil {
		return nil, fmt.Errorf("MeasurementService: Failed to retrieve all measurements: %w", err)
	}

	return measurements, nil
}

func (ms *MeasurementService) GetCurrentDayMeasurements() ([]*domain.Measurement, error) {
	measurements, err := ms.measurementRepository.GetCurrentDayMeasurements()
	if err != nil {
		return nil, fmt.Errorf("MeasurementService: Failed to retrieve current day measurements: %w", err)
	}

	return measurements, nil
}
