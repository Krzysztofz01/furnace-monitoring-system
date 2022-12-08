import dotenv from 'dotenv';
import winston, { Logger } from "winston";
import path from "path";
import { SensorDeviceService } from "@services/sensor-device.service";
import { exit } from "process";
import { UnitOfWork } from "@server/unit-of-work";
import { Server } from '@server/server';

let _server: Server = undefined;
let _unitOfWork: UnitOfWork = undefined;
let _sensorDeviceService: SensorDeviceService = undefined;
let _logger: Logger = undefined;

let _disposing = false;

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

    process.on('SIGINT', () => disposeAneClose(0));
    process.on('SIGTERM', () => disposeAneClose(0));
    process.on('exit', () => disposeAneClose(0));

    _server.listen(Number(process.env.PORT));

} catch (error: unknown) {
    if (_logger != undefined) {
        _logger.error(`[Program] Unexpected error occured. ${error}`);
    }

    disposeAneClose(1);
}

function disposeAneClose(returnCode: number): void {
    if (_disposing) return;
    _disposing = true;
    
    if (_logger !== undefined) {
        _logger.info("[Program] Releaseing resources and closing the application.");
    }
    
    if (_server !== undefined) {
        _server.dispose();
        _server = undefined;
    }

    if (_unitOfWork !== undefined) {
        _unitOfWork.dispose();
        _unitOfWork = undefined;
    }

    if (_logger !== undefined) {
        _logger.end();
        _logger = undefined;
    }

    // NOTE: Workaround for the file logger flushing problem 
    setTimeout(() => {}, 1000 * 4);
    exit(returnCode);
}