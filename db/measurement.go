package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Krzysztofz01/furnace-monitoring-system/domain"
)

func PerformMeasurementTableMigration(database *sql.DB) error {
	if database == nil {
		return errors.New("db: the provided database connection instance is not initialzied")
	}

	measurementTableCreate := `
	CREATE TABLE IF NOT EXISTS fms_measurements (
		id TEXT NOT NULL PRIMARY KEY,
		host_id TEXT NOT NULL,
		temperature_sensor_one FLOAT NULL,
		temperature_sensor_two FLOAT NULL,
		temperature_sensor_three FLOAT NULL,
		air_contamination_percentage INTEGER NULL,
		timestamp_unix DATE NOT NULL);
	`

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	stmt, err := database.PrepareContext(ctx, measurementTableCreate)
	if err != nil {
		return fmt.Errorf("db: failed to preapre context for measurement table migration query: %w", err)
	}

	defer stmt.Close()

	if _, err = stmt.ExecContext(ctx); err != nil {
		return fmt.Errorf("db: measurement table migration quiery execution failed: %w", err)
	}

	return nil
}

func InsertMeasurement(database *sql.DB, measurement *domain.Measurement) error {
	if database == nil {
		return errors.New("db: the provided database connection instance is not initialzied")
	}

	measurementRowInsert := `
	INSERT INTO fms_measurements (
		id,
		host_id,
		temperature_sensor_one,
		temperature_sensor_two,
		temperature_sensor_three,
		air_contamination_percentage,
		timestamp_unix)
	VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	stmt, err := database.PrepareContext(ctx, measurementRowInsert)
	if err != nil {
		return fmt.Errorf("db: failed to preapre context for measurement insert query: %w", err)
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		measurement.Id,
		measurement.HostId,
		measurement.TemperatureChannelOne,
		measurement.TemperatureChannelTwo,
		measurement.TemperatureChannelThree,
		measurement.AirContaminationPercentage,
		measurement.TimestampUnix)

	if err != nil {
		return fmt.Errorf("db: measurement insert quiery execution failed: %w", err)
	}

	return nil
}
