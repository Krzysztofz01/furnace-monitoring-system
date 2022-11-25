#include "temperature-sensor.hh"

TemperatureSensorResult::TemperatureSensorResult() { }

TemperatureSensorResult::~TemperatureSensorResult() { }

TemperatureSensorResult TemperatureSensorResult::success(float temperature) {
    TemperatureSensorResult result;
    result.mFailureMessage = "";
    result.mIsSuccess = true;
    result.mTemperature = temperature;

    return result
}

TemperatureSensorResult TemperatureSensorResult::failure(String failureMessage) {
    TemperatureSensorResult result;
    result.mFailureMessage = failureMessage;
    result.mIsSuccess = false;
    result.mTemperature = 0.0f;

    return result
}

bool TemperatureSensorResult::isSuccess() {
    return this.mIsSuccess;
}

float TemperatureSensorResult::getTemperature() {
    return this.mTemperature;
}

String TemperatureSensorResult::getFailureMessage() {
    return this.mFailureMessage;
}

TemperatureSensor::TemperatureSensor(const int sensorIndetifier, const int pinSensorCommunication) {
    // TODO: Checking for id duplicates
    this.mSensorIndetifier = sensorIndetifier;
    
    if (pinSensorCommunication <= 0) {
        throw std::runtime_error("TemperatureSensor: Invalid sensor pin sepcified.");
    }

    this.mPinSensorCommunication = pinSensorCommunication;
    this.mpOneWireSensorDevice = new OneWire(this.mPinSensorCommunication);
}

TemperatureSensor::~TemperatureSensor() {
    delete this.mpOneWireSensorDevice;
}

TemperatureSensorResult TemperatureSensor::readTemperature() {
    // TODO: Extract this !
    const int address_buffer_size = 8;
    byte address_buffer[address_buffer_size];
  
    if (!this.mpOneWireSensorDevice->search(address_buffer)) {
        this.mpOneWireSensorDevice->reset_search();
        delay(250);
        return TemperatureResult::failure("No addresses found to read sensor data.");
    }

    if (OneWire::crc8(address_buffer, 7) != address_buffer[7]) {
        this.mpOneWireSensorDevice->reset_search();
        return TemperatureResult::failure("Verification of address CRC checksum failed.");
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
            return TemperatureResult::failure("The sensor circuit is not supported");
    }

    this.mpOneWireSensorDevice->reset();
    this.mpOneWireSensorDevice->select(address_buffer);
    this.mpOneWireSensorDevice->write(0x44, 1);

    // TODO: This value can be fine-tuned down to 750ms?
    delay(1000);

    byte present = 0;
    present = this.mpOneWireSensorDevice->reset();
    this.mpOneWireSensorDevice->select(address_buffer);
    this.mpOneWireSensorDevice->write(0xBE);

    // TODO: Extract this !
    const int value_buffer_size = 12;
    byte value_buffer[value_buffer_size];

    for (int i = 0; i < value_buffer_size - 3; ++i) {
        value_buffer[i] = this.mpOneWireSensorDevice->read();
    }

    this.mpOneWireSensorDevice->read();
    this.mpOneWireSensorDevice->reset_search();

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
    return TemperatureResult::success(celsius_temperature);
}