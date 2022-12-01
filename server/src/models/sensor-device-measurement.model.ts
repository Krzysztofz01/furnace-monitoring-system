export interface SensorDeviceMeasurement {
    id: number;
    temperatureSensorOne: number | undefined;
    temperatureSensorTwo: number | undefined;
    temperatureSensorThree: number | undefined;
    airContaminationPercentage: number | undefined;
    timestamp: Date;
}