export interface Measurement {
    TemperatureChannelOne: number;
    TemperatureChannelTwo: number;
    TemperatureChannelThree: number;
    AirContaminationPercentage: number;
    Timestamp: Date;
}

export function GetMeasurementFromMeasurementPayload(eventPayload: string): Measurement | null {
    if (eventPayload.length == 0) return null;

    const payloadParts = eventPayload.split(";");
    if (payloadParts.length != 6) return null;

    return {
        TemperatureChannelOne: Number(payloadParts[2]),
        TemperatureChannelTwo: Number(payloadParts[3]),
        TemperatureChannelThree: Number(payloadParts[4]),
        AirContaminationPercentage: Number(payloadParts[5]),
        Timestamp: new Date()
    }    
}

export function GetMeasurementsFromApiResponse(requestMeasurements: Array<any>): Array<Measurement> {
    return requestMeasurements.map((measurement) => {        
        const parsedMeasurement: Measurement = {
            TemperatureChannelOne: measurement.temperatureOne,
            TemperatureChannelTwo: measurement.temperatureTwo,
            TemperatureChannelThree: measurement.temperatureThree,
            AirContaminationPercentage: measurement.airContamination,
            Timestamp: new Date(measurement.timestampUnix)
        }
    
        return parsedMeasurement
    })
}