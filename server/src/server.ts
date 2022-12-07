import { Server as HttpServer } from "http";
import express, { Express, NextFunction, Request, Response } from 'express';
import path from "path";
import { Logger } from "winston";
import { SensorDeviceService } from "@services/sensor-device.service";
import { ViewController } from "@controllers/view.controlller";
import { WebSocket, WebSocketServer } from "ws";
import { ErrorMiddleware } from "./middlewares/error.middleware";

const DEFAULT_APP_PORT = 80;

export class Server {
    private readonly _application: Express;
    private readonly _server: HttpServer;
    private readonly _socket: WebSocketServer;

    private readonly _logger: Logger;
    private readonly _sensorDeviceService: SensorDeviceService;

    private readonly _errorMiddleware: ErrorMiddleware;

    private readonly _viewController: ViewController;

    constructor(sensorDeviceServiceInstance: SensorDeviceService, loggerInstance: Logger) {
        if (loggerInstance === undefined) {
            throw new Error("[Server]: Provided logger instance is undefined.");
        }
        
        this._logger = loggerInstance;

        if (sensorDeviceServiceInstance === undefined) {
            throw new Error("[Server]: Provided sensordeviceservice instance is undefined.");
        }

        this._sensorDeviceService = sensorDeviceServiceInstance;

        this._application = express();
        this._server = new HttpServer(this._application);
        this._socket = new WebSocketServer({ server: this._server });

        this._errorMiddleware = new ErrorMiddleware(this._logger);

        this._viewController = new ViewController(this._sensorDeviceService, this._logger);

        this.configureExpress();
        this.configureExpressEndpoints();

        this.configureSocket();
    }

    private configureExpress(): void {
        // NOTE: Define EJS as the view engine for server-side view rendering
        this._application.set("views", path.join(__dirname, 'views'));
        this._application.set("view engine", "ejs");

        // NOTE: Add custom global error handling middleware
        this._application.use((_: any, request: Request, response: Response, next: NextFunction): void => {
            this._errorMiddleware.handle(request, response, next);
        });

        // NOTE: Define static files location
        this._application.use('/static', express.static(path.join(__dirname, 'static')));

        // NOTE: CORS configuration
        this._application.use((_: Request, response: Response, next: NextFunction): void => {
            response.setHeader("Access-Control-Allow-Origin", "*");
            response.setHeader("Access-Control-Allow-Credentials", "true");
            response.setHeader("Access-Control-Allow-Methods", "GET,HEAD,OPTIONS,POST,PUT");
            response.setHeader("Access-Control-Allow-Headers", "Access-Control-Allow-Origin,Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers,Authorization");
            next();
        });
    }

    private configureExpressEndpoints(): void {
        this._application.get('/', (request: Request, response: Response) => {
            this._viewController.handleIndex(request, response);
        });

        this._application.get('*', (_: Request, response: Response) => {
            response.statusCode = 301;
            response.redirect('/');
        });
    }

    private configureSocket(): void {
        this._socket.on('connection', (socket: WebSocket) => {
            
            // TODO: This is a dirty work-around
            const ipAddress = (socket as any)._socket.remoteAddress;

            this._logger.info(`[Server]: SensorDevice with address: ${ ipAddress } connected.`);

            socket.on('measurement', (measurementData: any) => {
                if (measurementData === undefined) {
                    this._logger.info("[Server]: SensorDevice measurement is undefined.");
                    return;
                }

                const result = this._sensorDeviceService.pushMeasurement(measurementData);
                if (!result.isSuccess) {
                    this._logger.warn("[Server]: SensorDevice measurement push failed.");
                }
            });

            socket.on('close', () => {
                this._logger.info(`[Server]: SensorDevice with address: ${ ipAddress } disconnected.`);
            });
        });
    }

    public listen(port: number | undefined = undefined): void {
        if (port === undefined) port = DEFAULT_APP_PORT;

        this._server.listen(port, () => {
            this._logger.info("[Server]: Server started to listening for requets.");
        });
    }

    public dispose(): void {
        this._logger.info("[Server]: Disposing the HTTP/Socket server.");

        this._server.close((error: Error) => {
            this._logger.error(`[Server]: Some problem occured while disposing. ${error}`);
        });
    }
}