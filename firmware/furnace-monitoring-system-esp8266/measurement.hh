#ifndef MEASUREMENT_HH
#define MEASUREMENT_HH

class Measurement {
public:
    Measurement() {
        TemperatureSensorOne = 0.0;
        TemperatureSensorTwo = 0.0;
        TemperatureSensorThree = 0.0;
        AirContamination = 0;
    }

    float TemperatureSensorOne;
    float TemperatureSensorTwo;
    float TemperatureSensorThree;
    int AirContamination;

    ~Measurement() {}
};

#endif // MEASUREMENT_HH