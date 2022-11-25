#include "temperature-sensor.hh"

TemperatureSensorResult::TemperatureSensorResult() { }

TemperatureSensorResult::~TemperatureSensorResult() { }

TemperatureSensorResult TemperatureSensorResult::success(float temperature) {
    TemperatureSensorResult result;
    result.mFailureMessage = "";
    result.mIsSuccess = true;
    result.mTemperature = temperature;

    return result;
}

TemperatureSensorResult TemperatureSensorResult::failure(String failureMessage) {
    TemperatureSensorResult result;
    result.mFailureMessage = failureMessage;
    result.mIsSuccess = false;
    result.mTemperature = 0.0f;

    return result;
}

bool TemperatureSensorResult::isSuccess() const {
    return mIsSuccess;
}

float TemperatureSensorResult::getTemperature() const {
    return mTemperature;
}

String TemperatureSensorResult::getFailureMessage() const {
    return mFailureMessage;
}

TemperatureSensor::TemperatureSensor(const int sensorIndetifier, const int pinSensorCommunication) {
    // TODO: Checking for id duplicates
    mSensorIndetifier = sensorIndetifier;
    
    if (pinSensorCommunication <= 0) {
        // FIXME: How to validate pin numbers?
        //throw std::runtime_error("TemperatureSensor: Invalid sensor pin sepcified.");
    }

    mPinSensorCommunication = pinSensorCommunication;
    mpOneWireSensorDevice = new OneWire(mPinSensorCommunication);
}

TemperatureSensor::~TemperatureSensor() {
    delete mpOneWireSensorDevice;
}

TemperatureSensorResult TemperatureSensor::readTemperature() {
    // TODO: Extract this !
    const int address_buffer_size = 8;
    byte address_buffer[address_buffer_size];
  
    if (!mpOneWireSensorDevice->search(address_buffer)) {
        mpOneWireSensorDevice->reset_search();
        delay(250);
        return TemperatureSensorResult::failure("No addresses found to read sensor data.");
    }

    if (OneWire::crc8(address_buffer, 7) != address_buffer[7]) {
        mpOneWireSensorDevice->reset_search();
        return TemperatureSensorResult::failure("Verification of address CRC checksum failed.");
    }

    // NOTE: Identification of the sensor circuit based on the first communication byte
    byte ic_identifier;
    switch (address_buffer[0]) {
        // NOTE: DS118S20 or old DS1820
        case 0x10:
            ic_identifier = 1;
            break;

        // NOTE: DS18B20
        case 0x28:
            ic_identifier = 0;
            break;

        // NOTE: DS1822
        case 0x22:
            ic_identifier = 0;
            break;

        // NOTE: Undefined sensor circuit
        default:
            return TemperatureSensorResult::failure("The sensor circuit is not supported");
    }

    mpOneWireSensorDevice->reset();
    mpOneWireSensorDevice->select(address_buffer);
    mpOneWireSensorDevice->write(0x44, 1);

    // TODO: This value can be fine-tuned down to 750ms?
    delay(1000);

    byte present = 0;
    present = mpOneWireSensorDevice->reset();
    mpOneWireSensorDevice->select(address_buffer);
    mpOneWireSensorDevice->write(0xBE);

    // TODO: Extract this !
    const int value_buffer_size = 12;
    byte value_buffer[value_buffer_size];

    for (int i = 0; i < value_buffer_size - 3; ++i) {
        value_buffer[i] = mpOneWireSensorDevice->read();
    }

    mpOneWireSensorDevice->read();
    mpOneWireSensorDevice->reset_search();

    // NOTE: Parsing the actual temperature for the data buffer. The result data size
    //       is always a 16 bit signed integer, and the type MUST be explicit in order
    //       to make the firmware work also for 32-bit processor's.
    int16_t raw_data = (value_buffer[1] << 8) | value_buffer[0];

    if (ic_identifier) {
        raw_data = raw_data << 3;

        if (value_buffer[7] == 0x10) {
            raw_data = (raw_data & 0xFFF0) + value_buffer_size - value_buffer[6];
        }
    } else {
        byte cfg = (value_buffer[4] & 0x60);

        if (cfg == 0x00) {
            // NOTE: 9bit resolution 94ms
            raw_data = raw_data & ~7;
        } else if (cfg == 0x20) {
            // NOTE: 10bit resolution 188ms
            raw_data = raw_data & ~3;
        } else if (cfg == 0x40) {
            // NOTE: 11bit resolution 375ms
            raw_data = raw_data & ~1;
        }

        // NOTE: Default resolution 750ms
    }

    float celsius_temperature = (float)raw_data / 16.0;
    return TemperatureSensorResult::success(celsius_temperature);
}

int TemperatureSensor::getIdentifier() const {
    return mSensorIndetifier;
}