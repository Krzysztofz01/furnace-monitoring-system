package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	VerboseLogging bool     `mapstructure:"verbose-logging"`
	SensorHostIds  []string `mapstructure:"sensor-host-ids"`
	Host           string   `mapstructure:"host"`
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
	viper.SafeWriteConfig()

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("config: failed to read the config file: %w", err)
	}

	if err := viper.Unmarshal(&Instance); err != nil {
		return fmt.Errorf("config: failed to unmarshal the config file content: %w", err)
	}

	return nil
}
