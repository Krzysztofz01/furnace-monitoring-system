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
        response.render("index");
    }
}