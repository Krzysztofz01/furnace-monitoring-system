export function GetMeasurementsServerEndpoint(): string | null {
    const server = window.location.host;
    if (!server) {
        console.error("Can not access the server host address");
        return null;
    }

    
    return `http://${server}/api/measurements`
}

export function GetDashboardSocketServerEndpoint(): string | null {
    const server = window.location.host;
    if (!server) {
        console.error("Can not access the server host address");
        return null;
    }

    return `ws://${server}/socket/dashboard`
}