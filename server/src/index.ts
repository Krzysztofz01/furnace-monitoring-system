import { Server } from "./server";
import dotenv from 'dotenv';
import { UnitOfWork } from "./unit-of-work";
import winston, { Logger } from "winston";
import path from "path";
import { exit } from "process";
import { SensorDeviceService } from "@services/sensor-device.service";

let _server: Server = undefined;
let _unitOfWork: UnitOfWork = undefined;
let _sensorDeviceService: SensorDeviceService = undefined;
let _logger: Logger = undefined;

try {
    dotenv.config();

    // TOOD: Prepare a better logging format
    _logger = winston.createLogger({
        level: process.env.LOG_LEVEL,
        format: winston.format.json(),
        transports: [
            new winston.transports.Console(),
            new winston.transports.File({ filename: path.join(__dirname, 'log/server.log') })
        ]
    });

    _unitOfWork = new UnitOfWork(
        process.env.DATABASE_NAME,
        Boolean(process.env.DATABASE_MIGRATE),
        _logger
    );

    _sensorDeviceService = new SensorDeviceService(
        _unitOfWork,
        _logger
    );

    _server = new Server(
        _sensorDeviceService,
        _logger
    );

    process.on('SIGINT', () => disposeHandlers());
    process.on('SIGTERM', () => disposeHandlers());
    process.on('exit', () => disposeHandlers());

    _server.listen(Number(process.env.PORT));

} catch (error: unknown) {
    if (_logger != undefined) {
        _logger.error(`Unexpected error occured. ${error}`);
    }

    disposeHandlers();
    exit(1);
}

function disposeHandlers(): void {
    if (_server !== undefined) {
        _server.dispose();
    }

    if (_unitOfWork !== undefined) {
        _unitOfWork.dispose();
    }

    if (_logger !== undefined) {
        _logger.end();
    }
}