package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	VerboseLogging bool `mapstructure:"verbose-logging"`
}

const (
	configFilePath string = "config.json"
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
