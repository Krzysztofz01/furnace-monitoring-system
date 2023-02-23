package main

import (
	"net/http"

	"github.com/Krzysztofz01/furnace-monitoring-system/config"
	"github.com/Krzysztofz01/furnace-monitoring-system/db"
	"github.com/Krzysztofz01/furnace-monitoring-system/server"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	if err := config.CreateConfig(); err != nil {
		panic("main: failed to create the config instance")
	}

	if err := db.CreateDatabase(true); err != nil {
		panic("main: failed to create the database driver instance")
	}

	if err := server.CreateWebSocketServer(); err != nil {
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

	e.Start(":5000")
}
