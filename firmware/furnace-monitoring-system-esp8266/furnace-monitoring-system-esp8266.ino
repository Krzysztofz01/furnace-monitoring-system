#include <OneWire.h>
#include <LiquidCrystal.h>

#include "lcd-display.hh"
#include "temperature-sensor.hh"

// TODO: Define all const values here

LcdDisplay* p_lcdDisplay = nullptr;

TemperatureSensor* p_sensor1 = nullptr;
TemperatureSensor* p_sensor2 = nullptr;

const int LCD_WIDTH = 16;
const int LCD_HEIGHT = 2;

void print_temperature(int sensorIdentifier, float sensorTemperature) {
  char* printBuffer = new char[LCD_WIDTH];
  snprintf(printBuffer, LCD_WIDTH, "[%d] Temp: %.2f C", sensorIdentifier, sensorTemperature);

  p_lcdDisplay->writeLine(sensorIdentifier, printBuffer);
  delete[] printBuffer;
}

void setup(void) {
  Serial.begin(9600);
  
  try {
    p_lcdDisplay = new LcdDisplay(
      LCD_WIDTH,
      LCD_HEIGHT,
      D1,
      D2,
      D5, D6, D7, D8);

    p_sensor1 = new TemperatureSensor(0, D3);
    p_sensor2 = new TemperatureSensor(1, D4);

  } catch (std::exception& ex) {
    // TODO: Better logging
    Serial.println(ex.what());
  }
}

void loop(void) {
  try {
    auto resultSensor1 = p_sensor1->readTemperature();
    if (resultSensor1.isSuccess()) {
      print_temperature(p_sensor1->getIdentifier(), resultSensor1.getTemperature());
    } else {
      // TODO: Log error
    }

    auto resultSensor2 = p_sensor2->readTemperature();
    if (resultSensor2.isSuccess()) {
      print_temperature(p_sensor2->getIdentifier(), resultSensor2.getTemperature());
    } else {
      // TODO: Log error
    }

  } catch (std::exception& ex) {
    // TODO: Better logging
    Serial.println(ex.what());
  }
}