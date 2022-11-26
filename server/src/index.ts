import { Request, Response } from "express";
import { Server } from "./server";

const server = new Server();

server.app.get('/', (_: Request, response: Response) => {
    response.send("Hello World! This is the furnace-monitoring-system server!");
});

server.start();