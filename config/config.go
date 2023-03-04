package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	VerboseLogging   bool             `mapstructure:"verbose-logging"`
	SensorHostIds    []string         `mapstructure:"sensor-host-ids"`
	SensorGpioConfig SensorGpioConfig `mapstructure:"sensor-gpio-config"`
}

type SensorGpioConfig struct {
	TempSensorData1   string `mapstructure:"temp-1"`
	TempSensorData2   string `mapstructure:"temp-2"`
	TempSensorData3   string `mapstructure:"temp-3"`
	LcdScreenRs       string `mapstructure:"lcd-rs"`
	LcdScreenE        string `mapstructure:"lcd-e"`
	LcdScreenDataBus4 string `mapstructure:"lcd-d4"`
	LcdScreenDataBus5 string `mapstructure:"lcd-d5"`
	LcdScreenDataBus6 string `mapstructure:"lcd-d6"`
	LcdScreenDataBus7 string `mapstructure:"lcd-d7"`
}

const (
	configFilePath string = "config/config.json"
)

var Instance *Config

func CreateConfig() error {
	viper.SetConfigFile(configFilePath)

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("config: failed to read the config file: %w", err)
	}

	if err := viper.Unmarshal(&Instance); err != nil {
		return fmt.Errorf("config: failed to unmarshal the config file content: %w", err)
	}

	return nil
}
