package main

import (
	"fmt"
	"net/http"

	"github.com/Krzysztofz01/furnace-monitoring-system/db"
	"github.com/Krzysztofz01/furnace-monitoring-system/server"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	if err := db.CreateDatabase(true); err != nil {
		// TODO: Panic here!
		fmt.Println(err)
	}

	if err := server.CreateWebSocketServer(); err != nil {
		// TODO: Panic here!
		fmt.Println(err)
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
		fmt.Println("Panel endpoint hit")
		server.Instance.UpgradeDashboardHostConnection(c.Request(), c.Response().Writer)
		return nil
	})

	e.Logger.Fatal(e.Start(":5000"))
}
