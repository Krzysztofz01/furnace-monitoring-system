#ifndef SENSOR_CONFIG_HH
#define SENSOR_CONFIG_HH

#include <ArduinoJson.h>
#include <stdexcept>

class SensorConfig {
public:
    int temp1;
    int temp2;
    int temp3;
    int rs;
    int e;
    int d4;
    int d5;
    int d6;
    int d7;

    static SensorConfig parseFromJson(String sensorConfigJsonString);
private:
    static const int sJsonBufferSize;
};

#endif // SENSOR_CONFIG_HH