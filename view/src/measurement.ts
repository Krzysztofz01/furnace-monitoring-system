export interface Measurement {
    TemperatureChannelOne: number;
    TemperatureChannelTwo: number;
    TemperatureChannelThree: number;
    AirContaminationPercentage: number;
}

export function GetMeasurementFromMeasurementPayload(eventPayload: string): Measurement {
    return null
}

export function GetMeasurementsFromApiResponse(measuremnets: Array<object>): Array<Measurement> {
    return null
}