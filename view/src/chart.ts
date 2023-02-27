import { Chart } from "chart.js";
import { format } from "date-fns";
import { Measurement } from "./measurement";

const maxMeasurementCount = 25
const canvasElement = document.getElementById('readings-chart-canvas');

export function CreateChart(): Chart {
    return new Chart(canvasElement as any, {
        type: 'line',
        data: {
            labels: [],
            datasets: [
                {
                    label: "First sensor",
                    data: [],
                    fill: false,
                    borderColor: 'rgb(234, 88, 12)'
                },
                {
                    label: "Second sensor",
                    data: [],
                    fill: false,
                    borderColor: 'rgb(249, 115, 22)'
                },
                {
                    label: "Third sensor",
                    data: [],
                    fill: false,
                    borderColor: 'rgb(251, 146, 60)'
                }
            ]
        },
        options: {
            maintainAspectRatio: false,
            scales: {
                y: {
                    beginAtZero: true
                }
            }
        }
    })
}

export function PushMeasurementIntoChart(chart: Chart, measurement: Measurement): void {
    PushLabelIntoChart(chart, GetTimestampLabel(measurement.Timestamp));

    PushMeasurementDataIntoChart(chart, measurement.TemperatureChannelOne, 0);
    PushMeasurementDataIntoChart(chart, measurement.TemperatureChannelTwo, 1);
    PushMeasurementDataIntoChart(chart, measurement.TemperatureChannelThree, 2);
    
    chart.update();
}

export function PushMeasurementsIntoChart(chart: Chart, measurements: Array<Measurement>): void {
    const labels = measurements.map((measurement) => GetTimestampLabel(measurement.Timestamp));
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
    values = values.slice(-maxMeasurementCount)
    const dataOverflow = values.length + chart.data.datasets[datasetIndex].data.length - maxMeasurementCount
    if (dataOverflow > 0) {
        chart.data.datasets[datasetIndex].data.splice(0, dataOverflow);
    }

    chart.data.datasets[datasetIndex].data.push(...values);
}

function PushMultipleLabelsIntoChart(chart: Chart, values: Array<string>): void {
    values = values.slice(-maxMeasurementCount)
    const dataOverflow = values.length + chart.data.labels.length - maxMeasurementCount;
    if (dataOverflow > 0) {
        chart.data.labels.splice(0, dataOverflow);
    }

    chart.data.labels.push(...values);
}

function GetTimestampLabel(date: Date): string {
    const measurementDate = new Date(date)
    return format(measurementDate, "HH:mm:ss")
}