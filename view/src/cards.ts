import { Measurement } from "./measurement";

const temperatureOneValueElement = document.getElementById("temperature-one-value")
const temperatureTwoValueElement = document.getElementById("temperature-two-value")
const temperatureThreeValueElement = document.getElementById("temperature-three-value")
const airContaminationValueElement = document.getElementById("air-contamination-value")

export function PushMeasurementIntoCards(measurement: Measurement): void {
    temperatureOneValueElement.innerText = measurement.TemperatureChannelOne.toFixed(2)
    temperatureTwoValueElement.innerText = measurement.TemperatureChannelTwo.toFixed(2)
    temperatureThreeValueElement.innerText = measurement.TemperatureChannelThree.toFixed(2)
    
    const acPercentage = Math.max(0, Math.min(100, measurement.AirContaminationPercentage * 100))
    airContaminationValueElement.innerText = acPercentage.toFixed(0)
}