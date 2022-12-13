package db

import (
	"errors"
	"fmt"
	"furnace-monitoring-system-server/pkg/domain"

	"github.com/google/uuid"
)

type InMememoryDatabase struct {
	measurementsCollection map[uuid.UUID]*domain.Measurement

	MeasurementsRepository *domain.MeasurementRepository
}

type InMemoryMeasurementsRepository struct {
	measurementsCollection *map[uuid.UUID]*domain.Measurement
}

func CreateInMememoryDatabase() (*InMememoryDatabase, error) {
	database := new(InMememoryDatabase)
	database.measurementsCollection = make(map[uuid.UUID]*domain.Measurement)

	measurementRepository, err := createInMemoryMeasurementRepository(&database.measurementsCollection)
	if err != nil {
		return nil, fmt.Errorf("InMemoryDatabase: Failed to create InMemoryMeasurementRepository: %w", err)
	} else {
		database.MeasurementsRepository = &measurementRepository
	}

	return database, nil
}

func createInMemoryMeasurementRepository(data *map[uuid.UUID]*domain.Measurement) (domain.MeasurementRepository, error) {
	if data == nil {
		return nil, errors.New("InMemoryMeasurementsRepository: Provided data is nil")
	}

	return &InMemoryMeasurementsRepository{
		measurementsCollection: data,
	}, nil
}

func (mr *InMemoryMeasurementsRepository) GetLatestMeasurement() (*domain.Measurement, error) {
	return nil, errors.New("InMemoryMeasurementsRepository: Not implemented")
}

func (mr *InMemoryMeasurementsRepository) GetAllMeasurements() ([]*domain.Measurement, error) {
	measurements := make([]*domain.Measurement, len(*mr.measurementsCollection))
	for _, measurement := range *mr.measurementsCollection {
		measurements = append(measurements, measurement)
	}

	return measurements, nil
}

func (mr *InMemoryMeasurementsRepository) GetCurrentDayMeasurements() ([]*domain.Measurement, error) {
	return nil, errors.New("InMemoryMeasurementsRepository: Not implemented")
}

func (mr *InMemoryMeasurementsRepository) InsertMeasurement(m *domain.Measurement) error {
	if _, exists := (*mr.measurementsCollection)[m.Id]; exists {
		return errors.New("InMememoryMeasurementsRepository: Database already contains a element with given id")
	}

	(*mr.measurementsCollection)[m.Id] = m
	return nil
}
