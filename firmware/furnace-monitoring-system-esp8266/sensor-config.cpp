#include "sensor-config.hh"

const int SensorConfig::sJsonBufferSize = 192; 

SensorConfig SensorConfig::parseFromJson(String sensorConfigJsonString) {
    DynamicJsonBuffer jsonBuffer(sJsonBufferSize);
    JsonObject& parsedJson = jsonBuffer.parseObject(sensorConfigJsonString);
    if (!parsedJson.success()) {
        throw std::runtime_error("SensorConfig: Failed to parse the sensor gpio config server response.");
    }

    SensorConfig config;
    config.temp1 = parsedJson["temp1"].as<int>();
    config.temp2 = parsedJson["temp2"].as<int>();
    config.temp3 = parsedJson["temp3"].as<int>();
    config.rs = parsedJson["rs"].as<int>();
    config.e = parsedJson["e"].as<int>();
    config.d4 = parsedJson["d4"].as<int>();
    config.d5 = parsedJson["d5"].as<int>();
    config.d6 = parsedJson["d6"].as<int>();
    config.d7 = parsedJson["d7"].as<int>();
    
    return config;
}