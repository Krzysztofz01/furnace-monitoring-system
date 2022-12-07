import { SensorDeviceMeasurementRepository } from "@repositories/sensor-device-measurement.repository";
import sqlite3 from 'better-sqlite3';
import { Logger } from "winston";

export class UnitOfWork {
    private readonly _database: sqlite3.Database;
    private readonly _logger: Logger;
    private _sensorDeviceMeasurementRepository: SensorDeviceMeasurementRepository;

    constructor(databaseName: string, migrateDatabase: boolean, loggerInstance: Logger) {
        if (loggerInstance === undefined) {
            throw new Error("[UnitOfWork]: Provided logger instance is undefined.");
        }
        
        this._logger = loggerInstance;

        if (databaseName === undefined) {
            throw new Error("[UnitOfWork]: Provided database name is undefined.");
        }

        this._database = new sqlite3(databaseName);
        if (this._database === undefined) {
            throw new Error("[UnitOfWork]: Failed to open the database.");
        }
                
        this._logger.info("[UnitOfWork]: Database access gained successful.");

        this._sensorDeviceMeasurementRepository = new SensorDeviceMeasurementRepository(this._database, this._logger, migrateDatabase);
    }

    public dispose(): void {
        this._logger.info("[UnitOfWork]: Disposing the database access.");

        try {
            this._database.close();
        } catch (error) {
            this._logger.info(`[UnitOfWork]: Some problem occured while disposing. ${error}`);
        }
    }

    public get sensorDeviceMeasurementRepository(): SensorDeviceMeasurementRepository {
        return this._sensorDeviceMeasurementRepository;
    }
}