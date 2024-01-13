package log

import (
	"fmt"
	"io"
	"os"

	"github.com/Krzysztofz01/furnace-monitoring-system/config"
	"github.com/sirupsen/logrus"
)

const logFilePath string = "furnace-monitoring-system-server.log"

var logFile *os.File

var Instance *logrus.Logger

func CreateLogger() (func(), error) {
	// TODO: Implement custom formatter
	var err error
	logFile, err = os.OpenFile(logFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("log: failed to access the log file: %w", err)
	}

	level := logrus.InfoLevel
	if config.Instance.VerboseLogging {
		level = logrus.DebugLevel
	}

	out := io.MultiWriter(logFile, os.Stdout)
	Instance = &logrus.Logger{
		Out:       out,
		Formatter: &logrus.TextFormatter{},
		Level:     level,
	}

	return func() {
		logFile.Close()
	}, nil
}
