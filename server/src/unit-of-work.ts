import { SensorDeviceMeasurementRepository } from "@repositories/sensor-device-measurement.repository";
import { Database, OPEN_CREATE, OPEN_FULLMUTEX, OPEN_READWRITE } from "sqlite3";
import { Logger } from "winston";

export class UnitOfWork {
    private readonly _database: Database;
    private readonly _logger: Logger;
    private _sensorDeviceMeasurementRepository: SensorDeviceMeasurementRepository;

    constructor(databaseName: string, migrateDatabase: boolean, loggerInstance: Logger) {
        if (loggerInstance === undefined) {
            throw new Error("UnitOfWork: Provided logger instance is undefined.");
        }
        
        this._logger = loggerInstance;

        if (databaseName === undefined) {
            throw new Error("UnitOfWork: Provided database name is undefined.");
        }

        this._database = new Database(databaseName, OPEN_READWRITE | OPEN_CREATE | OPEN_FULLMUTEX, (error: Error) => {
            this._logger.error(`UnitOfWork: Failed to open the database. ${error}`);
        });

        this._logger.info("UnitOfWork: Database access gained successful.");

        this._sensorDeviceMeasurementRepository = new SensorDeviceMeasurementRepository(this._database, this._logger, migrateDatabase);
    }

    public dispose(): void {
        this._logger.info("UnitOfWork: Disposing the database access.");

        this._database.close((error: Error) => {
            this._logger.info(`UnitOfWork: Some problem occured while disposing. ${error}`);
        });
    }

    public get sensorDeviceMeasurementRepository(): SensorDeviceMeasurementRepository {
        return this._sensorDeviceMeasurementRepository;
    }
}