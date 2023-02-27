#include <OneWire.h>
#include <LiquidCrystal.h>
#include <ESP8266WiFi.h>
#include <ESP8266WiFiMulti.h>
#include <Hash.h>
#include <WebSocketsClient.h>

#include "lcd-display.hh"
#include "temperature-sensor.hh"
#include "payload-builder.hh"
#include "config.hh"

// TODO: Define all const values here

LcdDisplay* p_lcdDisplay = nullptr;

TemperatureSensor* p_sensor1 = nullptr;
TemperatureSensor* p_sensor2 = nullptr;

ESP8266WiFiMulti WiFiMulti;
WebSocketsClient webSocket;

PayloadBuilder* p_payloadBuilder = nullptr;

const int LCD_WIDTH = 16;
const int LCD_HEIGHT = 2;

bool socketIsConnected = false;
bool socketIsEstablished = false;

void handleWebSocketEvent(WStype_t type, uint8_t* payload, size_t length) {
  hexdump(payload, length);
  
  switch (type) {
    case WStype_DISCONNECTED: {
      Serial.println("Disconnected");
      socketIsConnected = false;
      socketIsEstablished = false;
      break;      
    }
    case WStype_CONNECTED: {
      Serial.println("Connected");
      socketIsConnected = true;
      socketIsEstablished = false;
      break;
    }
    default: {
      Serial.println(type);
      break;
    }
  }
}

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


    WiFiMulti.addAP(FMS_NETWORK_SSID.c_str(), FMS_NETWORK_PASSWORD.c_str());

	  //WiFi.disconnect();
	  while(WiFiMulti.run() != WL_CONNECTED) {
		  Serial.println("Connecting to the network...");
      delay(200);
	  }



    // WiFi.begin(FMS_NETWORK_SSID, FMS_NETWORK_PASSWORD);
    // while (WiFi.status() != WL_CONNECTED) {
    //   Serial.println("Connecting to the network...");
    //   delay(200);
    // }
    Serial.println("Wifi connected");

    p_payloadBuilder = new PayloadBuilder(FMS_HOSTID);

    webSocket.begin(FMS_SERVER_ADDRESS, FMS_SERVER_PORT, "/socket/sensor");
    webSocket.onEvent(handleWebSocketEvent);
    webSocket.setReconnectInterval(5000);
    webSocket.enableHeartbeat(15000, 3000, 2);

  } catch (std::exception& ex) {
    // TODO: Better logging
    Serial.println(ex.what());
  }
}

static unsigned long last_run;

void loop(void) {
  try {
    Serial.println("Runing websocket loop");
    unsigned long now = millis();
    webSocket.loop();
    Serial.println("After websocket loop");

    if (now - last_run < 6000) return;

    last_run = millis();  

    Serial.println("Program loop start");
    float temperatureOne = 0.0;
    float temperatureTwo = 0.0;
    float temperatureThree = 0.0;
    int airContamination = 0;

    auto resultSensor1 = p_sensor1->readTemperature();
    if (resultSensor1.isSuccess()) {
      temperatureOne = resultSensor1.getTemperature();
      print_temperature(p_sensor1->getIdentifier(), resultSensor1.getTemperature());
    } else {
      // TODO: Log error
    }

    auto resultSensor2 = p_sensor2->readTemperature();
    if (resultSensor2.isSuccess()) {
      temperatureTwo = resultSensor2.getTemperature();
      print_temperature(p_sensor2->getIdentifier(), resultSensor2.getTemperature());
    } else {
      // TODO: Log error
    }

    
    if (socketIsConnected && socketIsEstablished) {
      Serial.println("Connected and established");
      // TODO: Add error counting for auto restart
      String payload = p_payloadBuilder->BuildMeasurementPayload(temperatureOne, temperatureTwo, temperatureThree, airContamination);
      if (webSocket.sendTXT(payload)) {
        Serial.println("Sending measurement payload returned true");
      } else {
        Serial.println("Sending measurement payload returned false");
      }
    } else if (socketIsConnected && !socketIsEstablished) {
      String payload = p_payloadBuilder->BuildConnectedPayload();
      if (webSocket.sendTXT(payload)) {
        Serial.println("Sending connected payload returned true");
        socketIsEstablished = true;
      } else {
        Serial.println("Sending connected payload returned false");
      }
    }
  } catch (std::exception& ex) {
    // TODO: Better logging
    Serial.println(ex.what());
  }
}