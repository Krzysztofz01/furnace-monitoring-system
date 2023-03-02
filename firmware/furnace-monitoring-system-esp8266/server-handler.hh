#ifndef SERVER_HANDLER_HH
#define SERVER_HANDLER_HH

#include <Arduino.h>
#include <ESP8266WiFi.h>
#include <ESP8266WiFiMulti.h>
#include <WebSocketsClient.h>

#include "measurement.hh"
#include "payload-builder.hh"

class ServerHandler {
public:
    ServerHandler(const String hostId, const String networkSsid, const String networkPassword, const String serverAddress, const int serverPort);

    void handleCycle();
    void sendMeasurement(Measurement measurement);
    bool isErrorCountExceeded();

    ~ServerHandler();
private:
    ServerHandler() {}
    
    static void handleWebSocketEvent(WStype_t type, uint8_t* payload, size_t length) {
      switch (type) {
        case WStype_DISCONNECTED: {
            smIsNetworkConnectionEstablished = false;
            smIsProtocolConnectionEstablished = false;
            break;
        }
        case WStype_CONNECTED: {
            smIsNetworkConnectionEstablished = true;
            smIsProtocolConnectionEstablished = false;
            break;
        }
        default: {
            break;
        }
      }
    }

    const String mServerSensorSocketEndpoint = "/socket/sensor"; 
    const int mReconnectionInterval = 5000;
    const int mHeartbeatPingTimeout = 15000;
    const int mHeartbeatPongTimeout = 3000;
    const int mHeartbeatThreshold = 2;
    const int mMaxErrorCount = 5;

    ESP8266WiFiMulti* mpWifi;
    WebSocketsClient* mpWebSocket;
    PayloadBuilder* mpPayloadBuilder;

    static bool smIsNetworkConnectionEstablished;
    static bool smIsProtocolConnectionEstablished;
    int mErrorCount;
};

#endif // SERVER_HANDLER_HH