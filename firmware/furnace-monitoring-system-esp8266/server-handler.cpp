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