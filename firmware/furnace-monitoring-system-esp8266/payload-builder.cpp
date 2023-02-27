#include "payload-builder.hh"

PayloadBuilder::PayloadBuilder(String hostId) {
    if (hostId.length() == 0) {
        throw std::runtime_error("PayloadBuilder: Invalid sensor identifier provided");
    }
    
    mHostId = hostId;
}

String PayloadBuilder::BuildConnectedPayload() {
    String payload = "1;";
    payload.concat(mHostId);

    return payload;
}

String PayloadBuilder::BuildMeasurementPayload(float tempOne, float tempTwo, float tempThree, int airContamination) {
    String payload = "3;";
    payload.concat(mHostId);
    payload.concat(mPayloadSeparator);

    payload.concat(tempOne);
    payload.concat(mPayloadSeparator);
    
    payload.concat(tempTwo);
    payload.concat(mPayloadSeparator);

    payload.concat(tempThree);
    payload.concat(mPayloadSeparator);

    payload.concat(airContamination);
    return payload;
}
