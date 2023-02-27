#ifndef PAYLOAD_BUILDER_HH
#define PAYLOAD_BUILDER_HH

#include <Arduino.h>
#include <stdexcept>

class PayloadBuilder {
public:
    PayloadBuilder(String hostId);

    String BuildConnectedPayload();
    String BuildMeasurementPayload(float tempOne, float tempTwo, float tempThree, int airContamination);

    ~PayloadBuilder() {}
private:
    PayloadBuilder() {}
    
    const char mPayloadSeparator = ';';
    String mHostId;
};

#endif // PAYLOAD_BUILDER_HH