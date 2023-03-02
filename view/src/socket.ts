import { v4 as uuidv4 } from 'uuid';

let connectedEventPayloadAccepted = false;
let hostId: string = null;

export function GetHostId(): string {
    if (hostId == null || hostId.length == 0) {
        hostId = uuidv4();
    }

    return hostId;
}

export function SetSocketStateConnected(): void {
    connectedEventPayloadAccepted = true;
}

export function ResetSocketState(): void {
    connectedEventPayloadAccepted = false;
    hostId = uuidv4();
}

export function IsSocketStateConnected(): boolean {
    return connectedEventPayloadAccepted;
}

export function GetConnectedEventPayload(): string {
    return `4;${GetHostId()}`;
}