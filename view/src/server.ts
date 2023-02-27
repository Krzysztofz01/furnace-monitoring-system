const server = import.meta.env.VITE_SERVER as string

export function GetMeasurementsServerEndpoint(): string {
    if (!server) console.error("Can not access the server host address")
    return `http://${server}/api/measurements`
}

export function GetDashboardSocketServerEndpoint(): string {
    if (!server) console.error("Can not access the server host address")
    return `ws://${server}/socket/dashboard`
}