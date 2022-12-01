import { SensorDeviceMeasurement } from "@models/sensor-device-measurement.model";
import { Database, RunResult } from "sqlite3";
import { Logger } from "winston";

export class SensorDeviceMeasurementRepository {
    private readonly _database: Database;
    private readonly _logger: Logger;

    constructor(databaseInstance: Database, loggerInstance: Logger, runMigration: boolean = true) {
        if (loggerInstance === undefined) {
            throw new Error("SensorDeviceMeasurementRepository: Provided logger instance is undefined.");
        }
        
        this._logger = loggerInstance;

        if (databaseInstance === undefined) {
            throw new Error("SensorDeviceMeasurementRepository: Provided database reference is undefined.");
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

        this._logger.info("SensorDeviceMeasurementRepository: Starting the migration process.");

        this._database.exec(migrationQuery, (error: Error) => {
            if (error != null) {
                throw new Error(`SensorDeviceMeasurementRepository: Migration failed. ${error}`);
            }
        });

        this._logger.info("SensorDeviceMeasurementRepository: Migration process finished.");
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
        
        const resultData = new Array<SensorDeviceMeasurement>();

        this._database.each(queryString, (error: Error, row: any) => {
            if (error != null) {
                this._logger.error("SensorDeviceMeasurementRepository: Can not retrieve row for 'getMeasurementsOrderedByTimestamp' query.");
            }

            resultData.push({
                id: row.id,
                temperatureSensorOne: row.temperatureSensorOne,
                temperatureSensorTwo: row.temperatureSensorTwo,
                temperatureSensorThree: row.temperatureSensorThree,
                airContaminationPercentage: row.airContaminationPercentage,
                timestamp: row.timestamp
            });
        }, (error: Error, count: number) => {
            if (error != null) {
                this._logger.error("SensorDeviceMeasurementRepository: 'getMeasurementsOrderedByTimestamp' query failed.");
                resultData.splice(0, resultData.length);
            }

            this._logger.info(`SensorDeviceMeasurementRepository: 'getMeasurementsOrderedByTimestamp' query finished. Retrievied ${count} rows.`);
        });

        return resultData;
    }

    public insertMeasurement(measurement: SensorDeviceMeasurement): void {
        const queryString = `
            INSERT INTO Sensor_Device_Measurements
            (temperature_sensor_one, temperature_sensor_two, temperature_sensor_three, air_contamination_percentage, timestamp)
            VALUES(?, ?, ?, ?, ?);`;

        this._database.run(queryString, [
            measurement.temperatureSensorOne,
            measurement.temperatureSensorTwo,
            measurement.temperatureSensorThree,
            measurement.airContaminationPercentage,
            measurement.timestamp
        ], (_: RunResult, error: Error) => {
            if (error != null) {
                this._logger.error("SensorDeviceMeasurementRepository: 'insertMeasurement' insert query failed.");
            }
        });
    }


}