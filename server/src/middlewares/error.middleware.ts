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

    public handle(_1: Request, _2: Response, next: NextFunction): void {
        try {
            next();
        } catch(error: any) {
            this._logger.info("[ErrorMiddleware]: Error caught by the error middleware.");

            if (error instanceof Error) {
                this._logger.error(error.message);
            } else {
                const errorMessage = String(error);
                this._logger.error(errorMessage);
            }

            // TODO: Redirect to some kind of error page
        }
    }
}