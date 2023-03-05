#include "server-handler.hh"

bool ServerHandler::smIsNetworkConnectionEstablished = false;
bool ServerHandler::smIsProtocolConnectionEstablished = false;

ServerHandler::ServerHandler(const String hostId, const String networkSsid, const String networkPassword, const String serverAddress, const int serverPort) {
    // TODO: Param validation
    mpWifi = new ESP8266WiFiMulti();
    mpWebSocket = new WebSocketsClient();
    mpPayloadBuilder = new PayloadBuilder(hostId);

    smIsNetworkConnectionEstablished = false;
    smIsProtocolConnectionEstablished = false;
    mErrorCount = 0;

    mpWifi->addAP(networkSsid.c_str(), networkPassword.c_str());
    while (mpWifi->run() != WL_CONNECTED) {
        Serial.println("Trying to connected to the network...");
        delay(300);
    }

    Serial.println("Connected to the network");

    mpWebSocket->begin(serverAddress, serverPort, mServerSensorSocketEndpoint);
    mpWebSocket->onEvent(handleWebSocketEvent);
    mpWebSocket->setReconnectInterval(mReconnectionInterval);
    mpWebSocket->enableHeartbeat(mHeartbeatPingTimeout, mHeartbeatPongTimeout, mHeartbeatThreshold);
}

void ServerHandler::handleCycle() {
    mpWebSocket->loop();

    if (smIsNetworkConnectionEstablished && !smIsProtocolConnectionEstablished) {
        String payload = mpPayloadBuilder->buildConnectedPayload();
        if (mpWebSocket->sendTXT(payload)) {
            smIsProtocolConnectionEstablished = true;
        } else {
            mErrorCount += 1;
        }
    }
}

SensorConfig ServerHandler::pullSensorConfig() {
    Serial.println("Trying to pull the sensor configuration from the server...");
    
    if (mpHttpClient == nullptr) {
        mpHttpClient = new HTTPClient();
    }

    auto requestAddress = FMS_SERVER_ADDRESS + mServerSensorConfigEndpoint;

    WiFiClient client;
    mpHttpClient->begin(client, requestAddress);
    int requestResponseCode = mpHttpClient->GET();
    if (requestResponseCode != 200) {
        Serial.println("Failed to pull the configuration from the server.");
        throw std::runtime_error("ServerHandler: Failed to make a sensor config request to the server");
    }

    Serial.println("Configuration pulled from the server successful.");

    String responsePayload = mpHttpClient->getString();
    SensorConfig config = SensorConfig::parseFromJson(responsePayload);

    mpHttpClient->end();
    return config;
}

void ServerHandler::sendMeasurement(Measurement measurement) {
    // TODO: Measurement validation
    if (smIsNetworkConnectionEstablished && smIsProtocolConnectionEstablished) {
        String payload = mpPayloadBuilder->buildMeasurementPayload(measurement);
        if (!mpWebSocket->sendTXT(payload)) {
            mErrorCount += 1;
        }
    }
}

bool ServerHandler::isErrorCountExceeded() {
    return mErrorCount > mMaxErrorCount;
}

ServerHandler::~ServerHandler() {
    mpWebSocket->disconnect();
    delete mpWebSocket;
    mpWebSocket = nullptr;

    delete mpWifi;
    mpWifi = nullptr;

    delete mpPayloadBuilder;
    mpPayloadBuilder = nullptr;

    smIsNetworkConnectionEstablished = false;
    smIsProtocolConnectionEstablished = false;
}