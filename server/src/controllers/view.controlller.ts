import { SensorDeviceService } from "@services/sensor-device.service";
import { isDateCurrentDay } from "@utilities/date.utility";
import { Request, Response } from "express";
import { Logger } from "winston";

export class ViewController {
    private readonly _sensorDeviceService: SensorDeviceService;
    private readonly _logger: Logger;
    
    constructor(sensorDeviceServiceInstance: SensorDeviceService, loggerInstance: Logger) {
        if (loggerInstance === undefined) {
            throw new Error("[ViewController]: Provided logger instance is undefined.");
        }

        this._logger = loggerInstance;

        if (sensorDeviceServiceInstance === undefined) {
            throw new Error("[ViewController]: Provided sensordeviceservice instance is undefined.");
        }

        this._sensorDeviceService = sensorDeviceServiceInstance;
    }

    public handleIndex(_: Request, response: Response): void {
        const serviceResult = this._sensorDeviceService.popMeasurement();
        if (!serviceResult.isSuccess) {
            throw new Error("[ViewController]: Failed to obtain service data required to render index view.");
        }

        this._logger.info("[ViewController]: Data required to render index view obtained successful.");
     
        const latestMeasurement = serviceResult.value;

        const formatedTimestamp = this.formatDisplayTimestamp(latestMeasurement.timestamp);
        const formatedTemperatureOne = (latestMeasurement.temperatureSensorOne == undefined) ? "Sensor unavailable" : `${latestMeasurement.temperatureSensorOne}°C`;
        const formatedTemperatureTwo = (latestMeasurement.temperatureSensorTwo == undefined) ? "Sensor unavailable" : `${latestMeasurement.temperatureSensorTwo}°C`;
        const formatedTemperatureThree = (latestMeasurement.temperatureSensorThree == undefined) ? "Sensor unavailable" : `${latestMeasurement.temperatureSensorThree}°C`;
        const formatedAir = (latestMeasurement.airContaminationPercentage == undefined) ? "Sensor unavailable" : `${latestMeasurement.airContaminationPercentage}%`;

        response.render("index", {
            time: formatedTimestamp,
            temp1: formatedTemperatureOne,
            temp2: formatedTemperatureTwo,
            temp3: formatedTemperatureThree,
            air: formatedAir
        });
    }

    public handleError(_: Request, response: Response): void {
        let errorMessage = "Unknown error occured.";
                  
        response.render("error", {
            message: errorMessage
        });
    } 

    private formatDisplayTimestamp(date: Date): string {
        const dateObj = new Date(date);

        const minutes = (dateObj.getMinutes().toString().length === 1) ?
            `0${dateObj.getMinutes().toString()}` : dateObj.getMinutes().toString();

        const hours = (dateObj.getHours().toString().length === 1) ?
            `0${dateObj.getHours().toString().length}` : dateObj.getHours().toString().length;

        if (isDateCurrentDay(dateObj)) {
            return `${ hours }:${ minutes }`;
        }

        const day = (dateObj.getDate().toString().length === 1) ?
            `0${ dateObj.getDate().toString() }` : dateObj.getDate().toString();
        
        const month = (dateObj.getMonth().toString().length === 1) ?
            `0${ dateObj.getMonth().toString() }` : dateObj.getMonth().toString();
        const year = dateObj.getFullYear().toString();
        
        return `${ hours }:${ minutes } ${ day }-${ month }-${ year }`;
    }
}