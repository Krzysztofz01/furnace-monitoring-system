package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	VerboseLogging   bool             `mapstructure:"verbose-logging"`
	SensorHostIds    []string         `mapstructure:"sensor-host-ids"`
	Host             string           `mapstructure:"host"`
	SensorGpioConfig SensorGpioConfig `mapstructure:"sensor-gpio-config"`
}

type SensorGpioConfig struct {
	TempSensorData1      string `mapstructure:"temp-1"`
	TempSensorData2      string `mapstructure:"temp-2"`
	TempSensorData3      string `mapstructure:"temp-3"`
	AirContaminationData string `mapstructure:"air-1"`
	LcdScreenRs          string `mapstructure:"lcd-rs"`
	LcdScreenE           string `mapstructure:"lcd-e"`
	LcdScreenDataBus4    string `mapstructure:"lcd-d4"`
	LcdScreenDataBus5    string `mapstructure:"lcd-d5"`
	LcdScreenDataBus6    string `mapstructure:"lcd-d6"`
	LcdScreenDataBus7    string `mapstructure:"lcd-d7"`
}

const (
	configFileName string = "config"
	configFileType string = "json"
)

var Instance *Config

func CreateConfig() error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("config: failed to read the current working directory: %w", err)
	}

	viper.AddConfigPath(cwd)
	viper.SetConfigName(configFileName)
	viper.SetConfigType(configFileType)

	viper.SetDefault("verbose-logging", false)
	viper.SetDefault("sensor-host-ids", []string{})
	viper.SetDefault("host", ":5000")
	viper.SetDefault("temp-1", "")
	viper.SetDefault("temp-2", "")
	viper.SetDefault("temp-3", "")
	viper.SetDefault("air-1", "")
	viper.SetDefault("lcd-rs", "")
	viper.SetDefault("lcd-e", "")
	viper.SetDefault("lcd-d4", "")
	viper.SetDefault("lcd-d5", "")
	viper.SetDefault("lcd-d6", "")
	viper.SetDefault("lcd-d7", "")
	viper.SafeWriteConfig()

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("config: failed to read the config file: %w", err)
	}

	if err := viper.Unmarshal(&Instance); err != nil {
		return fmt.Errorf("config: failed to unmarshal the config file content: %w", err)
	}

	return nil
}
