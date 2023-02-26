import { Chart } from "chart.js";
import { Measurement } from "./measurement";

const maxMeasurementCount = 50

export function PushMeasurementIntoChart(chart: Chart, measurement: Measurement): void {
    // TODO: Display timestamp as label
    PushLabelIntoChart(chart, "");
    
    PushMeasurementDataIntoChart(chart, measurement.TemperatureChannelOne, 0);
    PushMeasurementDataIntoChart(chart, measurement.TemperatureChannelTwo, 1);
    PushMeasurementDataIntoChart(chart, measurement.TemperatureChannelThree, 2);
    
    chart.update();
}

export function PushMeasurementsIntoChart(chart: Chart, measurements: Array<Measurement>): void {
    // TODO: Display timestamp as label
    const labels = measurements.map((_) => "");
    PushMultipleLabelsIntoChart(chart, labels);
    
    const firstTemperatureValues = measurements.map((m) => m.TemperatureChannelOne);
    PushMultipleMeasurementDataIntoChart(chart, firstTemperatureValues, 0);
    
    const secondsTemperatureValues = measurements.map((m) => m.TemperatureChannelTwo);
    PushMultipleMeasurementDataIntoChart(chart, secondsTemperatureValues, 1);
    
    const thirdTemperatureValues = measurements.map((m) => m.TemperatureChannelThree);
    PushMultipleMeasurementDataIntoChart(chart, thirdTemperatureValues, 2);

    chart.update();
}

function PushMeasurementDataIntoChart(chart: Chart, value: number, datasetIndex: number): void {
    if (chart.data.datasets[datasetIndex].data.length == maxMeasurementCount) {
        chart.data.datasets[datasetIndex].data.shift();
    }

    chart.data.datasets[datasetIndex].data.push(value);
}

function PushLabelIntoChart(chart: Chart, value: string): void {
    if (chart.data.labels.length == maxMeasurementCount) {
        chart.data.labels.shift();
    }

    chart.data.labels.push(value);
}

function PushMultipleMeasurementDataIntoChart(chart: Chart, values: Array<number>, datasetIndex: number): void {
    const dataOverflow = values.length + chart.data.datasets[datasetIndex].data.length - maxMeasurementCount
    if (dataOverflow > 0) {
        chart.data.datasets[datasetIndex].data.splice(0, dataOverflow);
    }

    chart.data.datasets[datasetIndex].data.push(...values);
}

function PushMultipleLabelsIntoChart(chart: Chart, values: Array<string>): void {
    const dataOverflow = values.length + chart.data.labels.length - maxMeasurementCount;
    if (dataOverflow > 0) {
        chart.data.labels.splice(0, dataOverflow);
    }

    chart.data.labels.push(...values);
}