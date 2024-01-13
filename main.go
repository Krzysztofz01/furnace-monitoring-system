package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Krzysztofz01/furnace-monitoring-system/config"
	"github.com/Krzysztofz01/furnace-monitoring-system/db"
	"github.com/Krzysztofz01/furnace-monitoring-system/handler"
	"github.com/Krzysztofz01/furnace-monitoring-system/log"
	"github.com/Krzysztofz01/furnace-monitoring-system/server"
	"github.com/Krzysztofz01/furnace-monitoring-system/view"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

func main() {
	if err := config.CreateConfig(); err != nil {
		panic(fmt.Errorf("main: failed to create the config instance: %w", err))
	}

	if closeLogger, err := log.CreateLogger(); err != nil {
		panic(fmt.Errorf("main: failed to create the logger instance: %w", err))
	} else {
		defer closeLogger()
	}

	if err := db.CreateDatabase(true); err != nil {
		err = fmt.Errorf("main: failed to create the database driver instance: %s", err)
		log.Instance.Fatal(err)
		panic(err)
	}

	if err := server.CreateWebSocketServer(); err != nil {
		err = fmt.Errorf("main: failed to create the websocket server instance: %s", err)
		log.Instance.Fatal(err)
		panic(err)
	}

	e := echo.New()

	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, values middleware.RequestLoggerValues) error {
			log.Instance.WithFields(logrus.Fields{
				"URI":    values.URI,
				"status": values.Status,
			}).Info("request")

			return nil
		},
	}))

	e.Use(middleware.Gzip())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(view.EmbeddedWebApp())

	e.GET("api/measurements", handler.HandleLatestMeasurements)
	e.GET("api/sensor/config", handler.HandleSensorConfig)
	e.GET("socket/sensor", server.Instance.UpgradeSensorHostConnection)
	e.GET("socket/dashboard", server.Instance.UpgradeDashboardHostConnection)

	go func() {
		if err := e.Start(config.Instance.Host); err != nil && err != http.ErrServerClosed {
			log.Instance.Fatalf("main: shutting down the server due to a runtime error: %s", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Instance.Fatalf("main: http server shutdown failure: %s", err)
	}

	if err := db.Instance.Close(); err != nil {
		log.Instance.Fatalf("main: database closing failure: %s", err)
	}
}
