import { NextFunction, Request, Response } from "express";
import { Logger } from "winston";

export class ErrorMiddleware {
    private readonly _logger: Logger;

    constructor(loggerInstance: Logger) {
        if (loggerInstance === undefined) {
            throw new Error("[ErrorMiddleware]: Provided logger instance is undefined.");
        }

        this._logger = loggerInstance;
    }

    public handle(error: Error, _1: Request, response: Response, _: NextFunction): void {
        this._logger.info("[ErrorMiddleware]: Error caught by the error middleware.");
        
        if (error instanceof Error) {
            this._logger.error(error.message);
        } else {
            const errorMessage = String(error);
            this._logger.error(errorMessage);
        }

        response.status(500);
        response.redirect('error');
    }
}