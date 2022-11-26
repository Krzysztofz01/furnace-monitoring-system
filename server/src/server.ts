import * as http from "http";
import express, { Express, NextFunction, Request, Response } from 'express';

export class Server {
    private readonly _app: Express;
    private _server!: http.Server;

    constructor() {
        this._app = express();
        this._app.set("port", process.env.PORT || 5000);

        this.configureMiddleware();
    }

    private configureMiddleware(): void {
        // Cross-Origin-Resource-Sharing middleware config
        this._app.use((_: Request, response: Response, next: NextFunction) => {
            response.setHeader("Access-Control-Allow-Origin", "*");
            response.setHeader("Access-Control-Allow-Credentials", "true");
            response.setHeader("Access-Control-Allow-Methods", "GET,HEAD,OPTIONS,POST,PUT");
            response.setHeader("Access-Control-Allow-Headers", "Access-Control-Allow-Origin,Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers,Authorization");
            next();
        });
    }

    public start(): void {
        this._server = this._app.listen(this._app.get("port"));
    }

    public get app(): Express {
        return this._app;
    }

    public get server(): http.Server {
        return this._server;
    }
}