import './index.css'
import { Chart, registerables } from 'chart.js';
import { PushMeasurementIntoCards } from './cards';
import { CreateChart, PushMeasurementIntoChart, PushMeasurementsIntoChart } from './chart';
import { GetMeasurementFromMeasurementPayload, GetMeasurementsFromApiResponse } from './measurement';
import { GetDashboardSocketServerEndpoint, GetMeasurementsServerEndpoint } from './server';
import { GetConnectedEventPayload, IsSocketStateConnected, ResetSocketState, SetSocketStateConnected } from './socket';

Chart.register(...registerables);
const chart = CreateChart();

fetch(GetMeasurementsServerEndpoint())
  .then((response) => response.json())
  .then((data) => {
    const measurements = GetMeasurementsFromApiResponse(data);
    PushMeasurementIntoCards(measurements.at(-1));
    PushMeasurementsIntoChart(chart, measurements);
  })

const socket = new WebSocket(GetDashboardSocketServerEndpoint());

socket.addEventListener("open", (_) => {
  if (socket.readyState != WebSocket.OPEN || IsSocketStateConnected()) return;

  socket.send(GetConnectedEventPayload());
  SetSocketStateConnected();
});

socket.addEventListener("error", (event) => {
  console.error(event);
});

socket.addEventListener("close", (_) => {
  ResetSocketState();
});

socket.addEventListener("message", (event) => {
  const measurement = GetMeasurementFromMeasurementPayload(event.data);
  PushMeasurementIntoCards(measurement);
  PushMeasurementIntoChart(chart, measurement);
});
