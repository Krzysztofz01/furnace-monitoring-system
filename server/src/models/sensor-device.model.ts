import { SensorDeviceMeasurement } from "./sensor-device-measurement.model";

export interface SensorDevice {
    identifier: string;
    measurements: Array<SensorDeviceMeasurement>;
}