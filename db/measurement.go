package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Krzysztofz01/furnace-monitoring-system/domain"
	"github.com/Krzysztofz01/furnace-monitoring-system/log"
)

func PerformMeasurementTableMigration(database *sql.DB) error {
	log.Instance.Debugln("Running the measurement domain database migration.")

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

	ctx, cancelfunc := context.WithTimeout(context.Background(), contextTimeout)
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

	ctx, cancelfunc := context.WithTimeout(context.Background(), contextTimeout)
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
		return fmt.Errorf("db: measurement insert query execution failed: %w", err)
	}

	return nil
}

func GetMeasurementsFromLastHours(database *sql.DB, hours int) ([]domain.Measurement, error) {
	if database == nil {
		return nil, errors.New("db: the provided database connection instance is not initialzied")
	}

	if hours < 0 {
		return nil, errors.New("db: the provided hours value must be postive")
	}

	targetTimestamp := time.Now().Add(time.Duration(-hours) * time.Hour).Unix()

	measurementLastHoursRowsSelect := `
	SELECT 
		id,
		host_id,
		temperature_sensor_one,
		temperature_sensor_two,
		temperature_sensor_three,
		air_contamination_percentage,
		timestamp_unix
	FROM fms_measurements
	WHERE timestamp_unix > ?
	ORDER BY timestamp_unix ASC;
	`

	ctx, cancelfunc := context.WithTimeout(context.Background(), contextTimeout)
	defer cancelfunc()

	stmt, err := database.PrepareContext(ctx, measurementLastHoursRowsSelect)
	if err != nil {
		return nil, fmt.Errorf("db: failed to prepare context for measurement select last hours rows query: %w", err)
	}

	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, targetTimestamp)
	if err != nil {
		return nil, fmt.Errorf("db: measurement select last hours query execution failed: %w", err)
	}

	defer rows.Close()

	measurements := make([]domain.Measurement, 0)
	for rows.Next() {
		var measurement domain.Measurement
		err = rows.Scan(
			&measurement.Id,
			&measurement.HostId,
			&measurement.TemperatureChannelOne,
			&measurement.TemperatureChannelTwo,
			&measurement.TemperatureChannelThree,
			&measurement.AirContaminationPercentage,
			&measurement.TimestampUnix)

		if err != nil {
			return nil, fmt.Errorf("db: failed to map the query result row to measurement struc: %w", err)
		} else {
			measurements = append(measurements, measurement)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("db: failed to obtain measurement last hours query results: %w", err)
	}

	return measurements, nil
}
