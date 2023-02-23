package main

import (
	"net/http"

	"github.com/Krzysztofz01/furnace-monitoring-system/config"
	"github.com/Krzysztofz01/furnace-monitoring-system/db"
	"github.com/Krzysztofz01/furnace-monitoring-system/log"
	"github.com/Krzysztofz01/furnace-monitoring-system/server"
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

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello world!")
	})

	e.GET("socket/sensor", func(c echo.Context) error {
		server.Instance.UpgradeSensorHostConnection(c.Request(), c.Response().Writer)
		return nil
	})

	e.GET("socket/dashboard", func(c echo.Context) error {
		server.Instance.UpgradeDashboardHostConnection(c.Request(), c.Response().Writer)
		return nil
	})

	log.Instance.Info("Furnace monitoring system server started")
	log.Instance.Fatal(e.Start(":5000"))
}
