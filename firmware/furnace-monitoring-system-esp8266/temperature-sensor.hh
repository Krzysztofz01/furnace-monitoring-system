#ifndef TEMPERATURE_SENSOR_HH
#define TEMPERATURE_SENSOR_HH

#include <Onewire.h>
#include <stdexcept>

class TemperatureSensorResult {
public:
    static TemperatureSensorResult success(float temperature);
    static TemperatureSensorResult failure(String failureMessage);
    
    bool isSuccess() const;
    float getTemperature() const;
    String getFailureMessage() const;

    ~TemperatureSensorResult();
private:
    TemperatureSensorResult();

    String mFailureMessage;
    bool mIsSuccess;
    float mTemperature;
};

class TemperatureSensor {
public:
    TemperatureSensorResult readTemperature();
    int getIdentifier() const;

    TemperatureSensor(const int sensorIndetifier, const int pinSensorCommunication);
    ~TemperatureSensor();
private:
    OneWire* mpOneWireSensorDevice;
    int mPinSensorCommunication;
    int mSensorIndetifier;
};

#endif // TEMPERATURE_SENSOR_HH