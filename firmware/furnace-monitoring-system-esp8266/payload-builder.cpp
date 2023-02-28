#include "payload-builder.hh"

PayloadBuilder::PayloadBuilder(String hostId) {
    if (hostId.length() == 0) {
        throw std::runtime_error("PayloadBuilder: Invalid sensor identifier provided");
    }
    
    mHostId = hostId;
}

String PayloadBuilder::buildConnectedPayload() {
    String payload = String(mConnectedEventPayloadType);
    payload.concat(mPayloadSeparator);
    payload.concat(mHostId);
    return payload;
}

String PayloadBuilder::buildMeasurementPayload(Measurement measurement) {
    String payload = String(mMeasurementEventPayloadType);
    payload.concat(mPayloadSeparator);
    payload.concat(mHostId);
    payload.concat(mPayloadSeparator);
    payload.concat(measurement.TemperatureSensorOne);
    payload.concat(mPayloadSeparator);
    payload.concat(measurement.TemperatureSensorTwo);
    payload.concat(mPayloadSeparator);
    payload.concat(measurement.TemperatureSensorThree);
    payload.concat(mPayloadSeparator);
    payload.concat(measurement.AirContamination);
    return payload;
}
