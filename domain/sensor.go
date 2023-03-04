package domain

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/Krzysztofz01/furnace-monitoring-system/config"
)

type SensorConfig struct {
	TemperatureSensorOneDataGpio   int
	TemperatureSensorTwoDataGpio   int
	TemperatureSensorThreeDataGpio int
	LcdScreenRsGpio                int
	LcdScreenEGpio                 int
	LcdScreenDataBus4Gpio          int
	LcdScreenDataBus5Gpio          int
	LcdScreenDataBus6Gpio          int
	LcdScreenDataBus7Gpio          int
}

func CreateSensorConfigFromConfigSection(c config.SensorGpioConfig) (*SensorConfig, error) {
	config := new(SensorConfig)
	var err error

	if len(c.TempSensorData1) == 0 {
		return nil, errors.New("domain: invalid temperature sensor one gpio address specified")
	}

	if config.TemperatureSensorOneDataGpio, err = strconv.Atoi(c.TempSensorData1); err != nil {
		return nil, fmt.Errorf("domain: failed to parse the temperature sensor one gpio address")
	}

	if len(c.TempSensorData2) == 0 {
		return nil, errors.New("domain: invalid temperature sensor two gpio address specified")
	}

	if config.TemperatureSensorTwoDataGpio, err = strconv.Atoi(c.TempSensorData2); err != nil {
		return nil, fmt.Errorf("domain: failed to parse the temperature sensor two gpio address")
	}

	if len(c.TempSensorData3) == 0 {
		return nil, errors.New("domain: invalid temperature sensor three gpio address specified")
	}

	if config.TemperatureSensorThreeDataGpio, err = strconv.Atoi(c.TempSensorData3); err != nil {
		return nil, fmt.Errorf("domain: failed to parse the temperature sensor three gpio address")
	}

	if len(c.LcdScreenRs) == 0 {
		return nil, errors.New("domain: invalid lcd screen rs gpio address specified")
	}

	if config.LcdScreenRsGpio, err = strconv.Atoi(c.LcdScreenRs); err != nil {
		return nil, fmt.Errorf("domain: failed to parse the lcd screen rs gpio address")
	}

	if len(c.LcdScreenE) == 0 {
		return nil, errors.New("domain: invalid lcd screen e gpio address specified")
	}

	if config.LcdScreenEGpio, err = strconv.Atoi(c.LcdScreenE); err != nil {
		return nil, fmt.Errorf("domain: failed to parse the lcd screen e gpio address")
	}

	if len(c.LcdScreenDataBus4) == 0 {
		return nil, errors.New("domain: invalid lcd screen d4 gpio address specified")
	}

	if config.LcdScreenDataBus4Gpio, err = strconv.Atoi(c.LcdScreenDataBus4); err != nil {
		return nil, fmt.Errorf("domain: failed to parse the lcd screen d4 gpio address")
	}

	if len(c.LcdScreenDataBus5) == 0 {
		return nil, errors.New("domain: invalid lcd screen d5 gpio address specified")
	}

	if config.LcdScreenDataBus5Gpio, err = strconv.Atoi(c.LcdScreenDataBus5); err != nil {
		return nil, fmt.Errorf("domain: failed to parse the lcd screen d5 gpio address")
	}

	if len(c.LcdScreenDataBus6) == 0 {
		return nil, errors.New("domain: invalid lcd screen d6 gpio address specified")
	}

	if config.LcdScreenDataBus6Gpio, err = strconv.Atoi(c.LcdScreenDataBus6); err != nil {
		return nil, fmt.Errorf("domain: failed to parse the lcd screen d6 gpio address")
	}

	if len(c.LcdScreenDataBus7) == 0 {
		return nil, errors.New("domain: invalid lcd screen d7 gpio address specified")
	}

	if config.LcdScreenDataBus7Gpio, err = strconv.Atoi(c.LcdScreenDataBus7); err != nil {
		return nil, fmt.Errorf("domain: failed to parse the lcd screen d7 gpio address")
	}

	return config, nil
}
