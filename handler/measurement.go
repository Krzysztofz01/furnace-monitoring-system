package handler

import (
	"net/http"

	"github.com/Krzysztofz01/furnace-monitoring-system/db"
	"github.com/Krzysztofz01/furnace-monitoring-system/domain"
	"github.com/Krzysztofz01/furnace-monitoring-system/log"
	"github.com/labstack/echo/v4"
)

func HandleLatestMeasurements(c echo.Context) error {
	log.Instance.Debug("Latest measurement API handler invoked")

	measurements, err := db.GetMeasurementsFromLastHours(db.Instance, 12)
	if err != nil {
		log.Instance.Errorf("Failed to retrievie latest measurements: %s", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	measurementDtos := make([]domain.MeasurementDto, 0, len(measurements))
	for _, measurement := range measurements {
		measurementDtos = append(measurementDtos, measurement.ToDto())
	}

	return c.JSON(http.StatusOK, measurementDtos)
}
