#include <OneWire.h>

#include "lcd-display.hh"
#include "temperature-sensor.hh"
#include "config.hh"
#include "server-handler.hh"
#include "measurement.hh"

// #define WEBSOCKETS_TCP_TIMEOUT (15000)

const int LCD_WIDTH = 16;
const int LCD_HEIGHT = 2;
const int MEASUREMENT_SENDING_INTERVAL_SECONDS = 10;

static unsigned long s_lastCycle;

LcdDisplay* p_lcdDisplay = nullptr;
ServerHandler* p_serverHandler = nullptr;
TemperatureSensor* p_sensor1 = nullptr;
TemperatureSensor* p_sensor2 = nullptr;

void print_temperature(int sensorIdentifier, float sensorTemperature);

void setup(void) {
  Serial.begin(9600);
  
  try {
    p_lcdDisplay = new LcdDisplay(
      LCD_WIDTH,
      LCD_HEIGHT,
      D1,
      D2,
      D5, D6, D7, D8);

    p_serverHandler = new ServerHandler(
      FMS_HOSTID,
      FMS_NETWORK_SSID,
      FMS_NETWORK_PASSWORD,
      FMS_SERVER_ADDRESS,
      FMS_SERVER_PORT);

    p_sensor1 = new TemperatureSensor(0, D3);
    p_sensor2 = new TemperatureSensor(1, D4);
  } catch (std::exception& ex) {
    // TODO: Better logging
    Serial.println(ex.what());
  }
}

void loop(void) {
  try {
    unsigned long currentCycle = millis();
    p_serverHandler->handleCycle();

    if (currentCycle - s_lastCycle < MEASUREMENT_SENDING_INTERVAL_SECONDS * 1000) return;

    Measurement measurement;

    auto resultSensor1 = p_sensor1->readTemperature();
    if (resultSensor1.isSuccess()) {
      float temperature = resultSensor1.getTemperature();

      measurement.TemperatureSensorOne = temperature;
      print_temperature(p_sensor1->getIdentifier(), temperature);
    } else {
      // TODO: Log error
    }

    auto resultSensor2 = p_sensor2->readTemperature();
    if (resultSensor2.isSuccess()) {
      float temperature = resultSensor2.getTemperature();

      measurement.TemperatureSensorTwo = temperature;
      print_temperature(p_sensor2->getIdentifier(), temperature);
    } else {
      // TODO: Log error
    }

    p_serverHandler->sendMeasurement(measurement);
    s_lastCycle = millis();
    
  } catch (std::exception& ex) {
    // TODO: Better logging
    Serial.println(ex.what());
  }
}

void print_temperature(int sensorIdentifier, float sensorTemperature) {
  char* printBuffer = new char[LCD_WIDTH];
  snprintf(printBuffer, LCD_WIDTH, "[%d] Temp: %.2f C", sensorIdentifier, sensorTemperature);

  p_lcdDisplay->writeLine(sensorIdentifier, printBuffer);
  delete[] printBuffer;
}