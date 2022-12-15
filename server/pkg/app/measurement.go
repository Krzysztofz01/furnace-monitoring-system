package app

import (
	"encoding/base64"
	"errors"
	"fmt"
	"furnace-monitoring-system-server/pkg/domain"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type MeasurementService struct {
	measurementRepository domain.MeasurementRepository
}

func CreateMeasurementService(measurementRepository domain.MeasurementRepository) (*MeasurementService, error) {
	if measurementRepository == nil {
		return nil, errors.New("MeasurementService: Provided MeasurementRepository reference is nil")
	}

	return &MeasurementService{
		measurementRepository: measurementRepository,
	}, nil
}

func (ms *MeasurementService) InsertMeasurement(measurementReading string) error {
	measurement, err := ms.parseMeasurementReading(measurementReading)
	if err != nil {
		return fmt.Errorf("MeasurementService: Failed parse the measurement: %w", err)
	}

	if err := ms.measurementRepository.InsertMeasurement(measurement); err != nil {
		return fmt.Errorf("MeasurementService: Failed to insert the measurement: %w", err)
	}

	return nil
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

func (ms *MeasurementService) parseMeasurementReading(measurementReading string) (*domain.Measurement, error) {
	if len(measurementReading) == 0 {
		return nil, errors.New("MeasurementService: The measurement reading is empty")
	}

	decoded, err := base64.StdEncoding.DecodeString(measurementReading)
	if err != nil {
		return nil, errors.New("MeasurementService: Failed to perform the readings base64 decoding")
	}

	readingTokens := strings.Split(string(decoded), ";")
	if len(readingTokens) != 6 {
		return nil, errors.New("MeasurementService: Invalid measurement reading format")
	}

	deviceIdHead := readingTokens[0]
	deviceIdTail := readingTokens[5]
	if deviceIdHead != deviceIdTail {
		return nil, errors.New("MeasurementService: The device id comparision checksum failed")
	}

	temparatureOne, err := ms.parseMeasurementTemperatureReading(readingTokens[1])
	if err != nil {
		return nil, fmt.Errorf("MeasurementService: Failed to parse temperature one reading: %w", err)
	}

	temparatureTwo, err := ms.parseMeasurementTemperatureReading(readingTokens[2])
	if err != nil {
		return nil, fmt.Errorf("MeasurementService: Failed to parse temperature two reading: %w", err)
	}

	temparatureThree, err := ms.parseMeasurementTemperatureReading(readingTokens[3])
	if err != nil {
		return nil, fmt.Errorf("MeasurementService: Failed to parse temperature three reading: %w", err)
	}

	airContamination, err := ms.parseMeasurementAirContaminationReading(readingTokens[4])
	if err != nil {
		return nil, fmt.Errorf("MeasurementService: Failed to parse air contamination reading: %w", err)
	}

	timestamp := time.Now().Unix()

	return &domain.Measurement{
		Id:                         uuid.New(),
		TemperatureSensorOne:       *temparatureOne,
		TemperatureSensorTwo:       *temparatureTwo,
		TemperatureSensorThree:     *temparatureThree,
		AirContaminationPercentage: *airContamination,
		TimestampUnix:              timestamp,
	}, nil
}

func (ms *MeasurementService) parseMeasurementTemperatureReading(temperatureReading string) (*domain.TemperatureReading, error) {
	if temperatureReading == "null" {
		return &domain.TemperatureReading{
			IsDefined: false,
			Value:     0,
		}, nil
	}

	parsedValue, err := strconv.ParseFloat(temperatureReading, 64)
	if err != nil {
		return nil, fmt.Errorf("MeasurementService: Failed to parse temperature value: %w", err)
	}

	return &domain.TemperatureReading{
		IsDefined: true,
		Value:     parsedValue,
	}, nil
}

func (ms *MeasurementService) parseMeasurementAirContaminationReading(airContaminationReading string) (*domain.AirContaminationReading, error) {
	if airContaminationReading == "null" {
		return &domain.AirContaminationReading{
			IsDefined: false,
			Value:     0,
		}, nil
	}

	parsedValue, err := strconv.ParseInt(airContaminationReading, 10, 0)
	if err != nil {
		return nil, fmt.Errorf("MeasurementService: Failed to parse air contamination value: %w", err)
	}

	if parsedValue > 100 || parsedValue < 0 {
		return nil, errors.New("MeasurementService: The givebn air contamination value is not a percentage value")
	}

	return &domain.AirContaminationReading{
		IsDefined: true,
		Value:     parsedValue,
	}, nil
}
