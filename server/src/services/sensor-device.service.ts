import { Result } from "@models/result";
import { SensorDeviceMeasurement } from "@models/sensor-device-measurement.model";
import { UnitOfWork } from "src/unit-of-work";
import { Logger } from "winston";

const UNSPECIFIED_ID = -1;

export class SensorDeviceService {
    private readonly _unitOfWork: UnitOfWork;
    private readonly _logger: Logger;

    constructor(unitOfWorkInstance: UnitOfWork, loggerInstance: Logger) {
        if (loggerInstance === undefined) {
            throw new Error("SensorDeviceService: Provided logger instance is undefined.");
        }

        this._logger = loggerInstance;

        if (unitOfWorkInstance === undefined) {
            throw new Error("SensorDeviceService: Provided unitofwork instance is undefined.");
        }

        this._unitOfWork = unitOfWorkInstance;
    }

    public pushMeasurement(encodedMeasurement: string): Result {
        try {
            const measurement = this.decodeMeasurement(encodedMeasurement, new Date(Date.now()));

            this._unitOfWork.sensorDeviceMeasurementRepository.insertMeasurement(measurement);

            return { isSuccess: true };
        } catch (error) {
            return { isSuccess: false };
        }
    }

    private decodeMeasurement(encodedMeasurement: string, timestamp: Date): SensorDeviceMeasurement | undefined {
        // NOTE: The current CSV-like format contains the data in order given below. The data is a ASCII string Base64 encoded.
        // DevideIdentifier;TempSen1;TempSen2;TempSen3;AirConSen;DevideIdentifier
        // MY_DEVICE       ;23.4    ;24.3    ;null    ;10       ;MY_DEVICE

        const asciiMeasurement = Buffer.from(encodedMeasurement, 'base64').toString('ascii');
        const measurementTokens = asciiMeasurement.split(';');

        if (measurementTokens.length !== 6) {
            this._logger.warn("SensorDeviceService: Can not decode measurement. Invalid format.");
            return undefined;
        }

        const deviceIdHead = measurementTokens[0];
        const deviceIdTail = measurementTokens[5];

        if (deviceIdHead !== deviceIdTail) {
            this._logger.warn("SensorDeviceService: Can not decode measurement. Invalid format, the measurement may be corrupted.");
            return undefined;
        }

        let temperatureSensorOne = undefined;
        if (measurementTokens[1] !== 'null') {
            temperatureSensorOne = parseFloat(measurementTokens[1]);
            if (Number.isNaN(temperatureSensorOne)) {
                this._logger.warn("SensorDeviceService: Can not decode measurement. Invalid sensor one temperature value.");
                return undefined;
            }
        }
        
        let temperatureSensorTwo = undefined;
        if (measurementTokens[2] !== 'null') {
            temperatureSensorTwo = parseFloat(measurementTokens[2]);
            if (Number.isNaN(temperatureSensorTwo)) {
                this._logger.warn("SensorDeviceService: Can not decode measurement. Invalid sensor two temperature value.");
                return undefined;
            }
        }

        let temperatureSensorThree = undefined;
        if (measurementTokens[3] !== 'null') {
            temperatureSensorThree = parseFloat(measurementTokens[3]);
            if (Number.isNaN(temperatureSensorThree)) {
                this._logger.warn("SensorDeviceService: Can not decode measurement. Invalid sensor three temperature value.");
                return undefined;
            }
        }

        let airContamination = undefined;
        if (measurementTokens[4] !== 'null') {
            airContamination = parseInt(measurementTokens[4]);
            if (Number.isNaN(airContamination)) {
                this._logger.warn("SensorDeviceService: Can not decode measurement. Invalid air contamination value.");
                return undefined;
            }

            if (airContamination < 0 || airContamination > 100) {
                this._logger.warn("SensorDeviceService: Can not decode measurement. Invalid air contamination value, not percentage range.");
                return undefined;
            }
        }

        return {
            id: UNSPECIFIED_ID,
            temperatureSensorOne: temperatureSensorOne,
            temperatureSensorTwo: temperatureSensorTwo,
            temperatureSensorThree: temperatureSensorThree,
            airContaminationPercentage: airContamination,
            timestamp: timestamp
        };
    }

}