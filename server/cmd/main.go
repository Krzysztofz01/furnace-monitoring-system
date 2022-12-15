package main

import (
	"furnace-monitoring-system-server/pkg/app"
	"furnace-monitoring-system-server/pkg/db"
	fhttp "furnace-monitoring-system-server/pkg/http"
)

func main() {
	_ = startup()
}

func startup() error {
	// NOTE: InMemory for testing purposes, will be replaced in the future with SQLite
	db, err := db.CreateInMememoryDatabase()
	if err != nil {
		// TODO: Wrap
		return err
	}

	measurementService, err := app.CreateMeasurementService(*db.MeasurementsRepository)
	if err != nil {
		// TODO: Wrap
		return err
	}

	server, err := fhttp.CreateServer(measurementService)
	if err != nil {
		// TODO: Wrap
		return err
	}

	return server.Listen()
}
