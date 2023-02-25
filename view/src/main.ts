import { Chart, registerables } from 'chart.js';
import { PushMeasurementIntoCards } from './cards';
import { PushMeasurementIntoChart, PushMeasurementsIntoChart } from './chart';
import './index.css'
import { GetMeasurementFromMeasurementPayload, GetMeasurementsFromApiResponse } from './measurement';

Chart.register(...registerables);

const server = "localhost:5000"

const canvasElement = document.getElementById('readings-chart-canvas');
const chart = new Chart(canvasElement as any, {
  type: 'line',
  data: {
    labels: [],
    datasets: [
      {
        label: "Sensor one",
        data: [],
        fill: false,
        borderColor: 'rgb(255, 0, 0)'
      },
      {
        label: "Sensor one",
        data: [],
        fill: false,
        borderColor: 'rgb(255, 0, 0)'
      },
      {
        label: "Sensor one",
        data: [],
        fill: false,
        borderColor: 'rgb(255, 0, 0)'
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

fetch(`http://${server}/api/measurements`)
  .then((response) => response.json())
  .then((data) => {
    const measurements = GetMeasurementsFromApiResponse(data)
    PushMeasurementsIntoChart(chart, measurements)
  })

var connectionPayloadAccepted = false
const hostId = crypto.randomUUID()
const socket = new WebSocket(`ws://${server}/socket/dashboard`)

socket.addEventListener("open", (_) => {
  if (socket.readyState != WebSocket.OPEN) return
  if (connectionPayloadAccepted) return

  socket.send(`4;${hostId}`)
  connectionPayloadAccepted = true
})

socket.addEventListener("error", (event) => {
  console.log(event)
})

socket.addEventListener("close", (_) => {
  connectionPayloadAccepted = false
})

socket.addEventListener("message", (event) => {
  const measurement = GetMeasurementFromMeasurementPayload(event.data)
  PushMeasurementIntoCards(measurement)
  PushMeasurementIntoChart(chart, measurement)
})
