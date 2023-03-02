package main

import (
	"github.com/Krzysztofz01/furnace-monitoring-system/config"
	"github.com/Krzysztofz01/furnace-monitoring-system/db"
	"github.com/Krzysztofz01/furnace-monitoring-system/handler"
	"github.com/Krzysztofz01/furnace-monitoring-system/log"
	"github.com/Krzysztofz01/furnace-monitoring-system/server"
	"github.com/Krzysztofz01/furnace-monitoring-system/view"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	if err := config.CreateConfig(); err != nil {
		panic("main: failed to create the config instance")
	}

	if disposeFunc, err := log.CreateLogger(); err != nil {
		panic("main: failed to create the logger instance")
	} else {
		defer disposeFunc()
	}

	if err := db.CreateDatabase(true); err != nil {
		log.Instance.Fatalf("Failed to create the database driver instance: %s", err)
		panic("main: failed to create the database driver instance")
	}

	if err := server.CreateWebSocketServer(); err != nil {
		log.Instance.Fatalf("Failed to create the websocket server instance: %s", err)
		panic("main: failed to create the websocket server instance")
	}

	e := echo.New()
	e.Use(middleware.Gzip())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(view.EmbeddedWebApp())

	e.GET("api/measurements", handler.HandleLatestMeasurements)
	e.GET("socket/sensor", server.Instance.UpgradeSensorHostConnection)
	e.GET("socket/dashboard", server.Instance.UpgradeDashboardHostConnection)

	log.Instance.Fatal(e.Start(":5000"))
}
