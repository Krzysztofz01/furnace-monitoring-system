import { SensorDeviceService } from "@services/sensor-device.service";
import { Request, Response } from "express";
import { Logger } from "winston";

export class ViewController {
    private readonly _sensorDeviceService: SensorDeviceService;
    private readonly _logger: Logger;
    
    constructor(sensorDeviceServiceInstance: SensorDeviceService, loggerInstance: Logger) {
        if (loggerInstance === undefined) {
            throw new Error("ViewController: Provided logger instance is undefined.");
        }

        this._logger = loggerInstance;

        if (sensorDeviceServiceInstance === undefined) {
            throw new Error("ViewController: Provided sensordeviceservice instance is undefined.");
        }

        this._sensorDeviceService = sensorDeviceServiceInstance;
    }

    public handleIndex(_: Request, response: Response): void {
        const serviceResult = this._sensorDeviceService.popMeasurement();
        if (!serviceResult.isSuccess) {
            this._logger.warn("ViewController: Failed to obtain service data required to render index view.");
            // TODO: Error view needs to be implemented here
            return; 
        }

        const latestMeasurement = serviceResult.value;
   
        // FIXME: Wrong timestamp, corrupted test data?
        const timestamp = new Date(latestMeasurement.timestamp);
        const formatedTimestamp = `${timestamp.getHours()}:${timestamp.getMinutes()}:${timestamp.getSeconds()}`;
        const formatedTemperatureOne = (latestMeasurement.temperatureSensorOne === undefined) ? "Sensor unavailable" : `${latestMeasurement.temperatureSensorOne}°C`;
        const formatedTemperatureTwo = (latestMeasurement.temperatureSensorTwo === undefined) ? "Sensor unavailable" : `${latestMeasurement.temperatureSensorTwo}°C`;
        const formatedTemperatureThree = (latestMeasurement.temperatureSensorThree === undefined) ? "Sensor unavailable" : `${latestMeasurement.temperatureSensorThree}°C`;
        const formatedAir = (latestMeasurement.airContaminationPercentage === undefined) ? "Sensor unavailable" : `${latestMeasurement.airContaminationPercentage}%`;

        response.render("index", {
            time: formatedTimestamp,
            temp1: formatedTemperatureOne,
            temp2: formatedTemperatureTwo,
            temp3: formatedTemperatureThree,
            air: formatedAir
        });
    }
}