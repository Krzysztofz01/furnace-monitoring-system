import { SensorDeviceMeasurement } from "@models/sensor-device-measurement.model";
import sqlite3 from 'better-sqlite3';
import { Logger } from "winston";

export class SensorDeviceMeasurementRepository {
    private readonly _database: sqlite3.Database;
    private readonly _logger: Logger;

    constructor(databaseInstance: sqlite3.Database, loggerInstance: Logger, runMigration: boolean = true) {
        if (loggerInstance === undefined) {
            throw new Error("[SensorDeviceMeasurementRepository]: Provided logger instance is undefined.");
        }
        
        this._logger = loggerInstance;

        if (databaseInstance === undefined) {
            throw new Error("[SensorDeviceMeasurementRepository]: Provided database reference is undefined.");
        }
        
        this._database = databaseInstance;

        if (runMigration) this.performMigration();
    }

    private performMigration(): void {
        const migrationQuery = `
            CREATE TABLE IF NOT EXISTS Sensor_Device_Measurements (
            id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
            temperature_sensor_one FLOAT NULL,
            temperature_sensor_two FLOAT NULL,
            temperature_sensor_three FLOAT NULL,
            air_contamination_percentage INTEGER NULL,
            timestamp DATE NOT NULL);`;

        this._logger.info("[SensorDeviceMeasurementRepository]: Starting the migration process.");

        try {
            this._database.exec(migrationQuery);
        } catch (error) {
            throw new Error(`[SensorDeviceMeasurementRepository]: Migration failed. ${error}`);
        }
        
        this._logger.info("[SensorDeviceMeasurementRepository]: Migration process finished.");
    }

    public getLatestMeasurement(): SensorDeviceMeasurement | undefined {
        const queryString = `
        SELECT
        id as 'id',
        temperature_sensor_one as 'temperatureSensorOne',
        temperature_sensor_two as 'temperatureSensorTwo',
        temperature_sensor_three as 'temperatureSensorThree',
        air_contamination_percentage as 'airContaminationPercentage',
        timestamp as 'timestamp'
        FROM Sensor_Device_Measurements
        ORDER BY timestamp DESC
        LIMIT 1;`;

        const statement = this._database.prepare(queryString);

        const result = statement.get();
        if (result === undefined) return undefined;

        return {
            id: result.id,
            temperatureSensorOne: result.temperatureSensorOne,
            temperatureSensorTwo: result.temperatureSensorTwo,
            temperatureSensorThree: result.temperatureSensorThree,
            airContaminationPercentage: result.airContaminationPercentage,
            timestamp: result.timestamp
        };
    }

    public getMeasurementsOrderedByTimestamp(): Array<SensorDeviceMeasurement> {
        const queryString = `
            SELECT
            id as 'id',
            temperature_sensor_one as 'temperatureSensorOne',
            temperature_sensor_two as 'temperatureSensorTwo',
            temperature_sensor_three as 'temperatureSensorThree',
            air_contamination_percentage as 'airContaminationPercentage',
            timestamp as 'timestamp'
            FROM Sensor_Device_Measurements
            ORDER BY timestamp DESC;`;
        
        const statement = this._database.prepare(queryString);
        
        const mappedResult: Array<SensorDeviceMeasurement> = statement.all().map((row) => {
            return {
                id: row.id,
                temperatureSensorOne: row.temperatureSensorOne,
                temperatureSensorTwo: row.temperatureSensorTwo,
                temperatureSensorThree: row.temperatureSensorThree,
                airContaminationPercentage: row.airContaminationPercentage,
                timestamp: row.timestamp
            }
        })

        return mappedResult;
    }

    public insertMeasurement(measurement: SensorDeviceMeasurement): void {
        const queryString = `
            INSERT INTO Sensor_Device_Measurements
            (temperature_sensor_one, temperature_sensor_two, temperature_sensor_three, air_contamination_percentage, timestamp)
            VALUES(?, ?, ?, ?, ?);`;

        const statement = this._database.prepare(queryString);

        const results = statement.run(
            measurement.temperatureSensorOne,
            measurement.temperatureSensorTwo,
            measurement.temperatureSensorThree,
            measurement.airContaminationPercentage,
            measurement.timestamp.getTime());

        if (results.changes !== 1) {
            throw new Error("[SensorDeviceMeasurementRepository]: Failed to insert measurement into the database.");
        }
    }
}