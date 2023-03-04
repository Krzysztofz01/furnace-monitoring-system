package handler

import (
	"net/http"

	"github.com/Krzysztofz01/furnace-monitoring-system/config"
	"github.com/Krzysztofz01/furnace-monitoring-system/domain"
	"github.com/Krzysztofz01/furnace-monitoring-system/log"
	"github.com/labstack/echo/v4"
)

func HandleSensorConfig(c echo.Context) error {
	log.Instance.Debug("Sensor config API handler invoked")

	sensorConfigSection := config.Instance.SensorGpioConfig
	config, err := domain.CreateSensorConfigFromConfigSection(sensorConfigSection)
	if err != nil {
		log.Instance.Debugf("Failed to parse the sensor config to domain model: %s", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, config)
}
