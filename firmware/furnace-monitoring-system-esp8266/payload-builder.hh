#ifndef PAYLOAD_BUILDER_HH
#define PAYLOAD_BUILDER_HH

#include <Arduino.h>
#include <stdexcept>

#include "measurement.hh"

class PayloadBuilder {
public:
    PayloadBuilder(String hostId);

    String buildConnectedPayload();
    String buildMeasurementPayload(Measurement measurement);

    ~PayloadBuilder() {}
private:
    PayloadBuilder() {}
    
    const char mPayloadSeparator = ';';
    const int mConnectedEventPayloadType = 1;
    const int mMeasurementEventPayloadType = 3;
    
    String mHostId;
};

#endif // PAYLOAD_BUILDER_HH