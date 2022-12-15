package http

import (
	"errors"
	"fmt"
	"furnace-monitoring-system-server/pkg/app"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// TODO: Implement the hadling of all headers
// TODO: Implement the handling of content-type
type EndpointHandler struct {
	measurementService *app.MeasurementService
	assetsHandler      *EmbedadAssetsHandler
	socketUpgrader     websocket.Upgrader
}

func CreateHandler(measurementService *app.MeasurementService, assetsHandler *EmbedadAssetsHandler) (*EndpointHandler, error) {
	if measurementService == nil {
		return nil, errors.New("EndpointHandler: Provided MeasurementService reference is nil")
	}

	if assetsHandler == nil {
		return nil, errors.New("EndpointHandler: Provided EmbededAssetsHandler reference is nil")
	}

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		// TODO: Add handshake timeout
	}

	return &EndpointHandler{
		measurementService: measurementService,
		socketUpgrader:     upgrader,
	}, nil
}

func (eh *EndpointHandler) HandleIndexTemplate(writer http.ResponseWriter, request *http.Request) {
	measurement, err := eh.measurementService.GetLatestMeasurement()
	if err != nil {
		fmt.Println(err)
		if _, err := writer.Write([]byte("[Placeholder]: Service data retrieval failure.")); err != nil {
			fmt.Println(err)
			return
		}
	}

	mTime := time.UnixMilli(measurement.TimestampUnix)
	tDuration := time.Since(mTime)
	tDiff := tDuration.Hours() / 24

	var timestamp = mTime.Format(time.ANSIC)
	if tDiff == 0 {
		timestamp = fmt.Sprintf("%d:%d", mTime.Hour(), mTime.Minute())
	}

	var temperatureOne = "Sensor unavailable."
	if measurement.TemperatureSensorOne.IsDefined {
		temperatureOne = fmt.Sprintf("%f°C", measurement.TemperatureSensorOne.Value)
	}

	var temperatureTwo = "Sensor unavailable."
	if measurement.TemperatureSensorTwo.IsDefined {
		temperatureTwo = fmt.Sprintf("%f°C", measurement.TemperatureSensorTwo.Value)
	}

	var temperatureThree = "Sensor unavailable."
	if measurement.TemperatureSensorThree.IsDefined {
		temperatureThree = fmt.Sprintf("%f°C", measurement.TemperatureSensorThree.Value)
	}

	var airContamination = "Sensor unavailable."
	if measurement.AirContaminationPercentage.IsDefined {
		airContamination = fmt.Sprintf("%d%s", measurement.AirContaminationPercentage.Value, "%")
	}

	templateData := struct {
		Timestamp        string
		Temperature1     string
		Temperature2     string
		Temperature3     string
		AirContamination string
	}{
		timestamp,
		temperatureOne,
		temperatureTwo,
		temperatureThree,
		airContamination,
	}

	template := eh.assetsHandler.GetEmbededViewTemplate("index")
	if err := template.Execute(writer, templateData); err != nil {
		fmt.Println(err)
		return
	}

	if _, err := writer.Write([]byte(fmt.Sprintf("[Placeholder]: Service data retrieved: %s", measurement.Id.String()))); err != nil {
		fmt.Println(err)
		return
	}
}

func (eh *EndpointHandler) HandleErrorTemplate(writer http.ResponseWriter, request *http.Request) {
	if _, err := writer.Write([]byte("[Placeholder]: Error page called.")); err != nil {
		fmt.Println(err)
		return
	}
}

func (eh *EndpointHandler) HandleSensorSocket(writer http.ResponseWriter, request *http.Request) {
	socket, err := eh.socketUpgrader.Upgrade(writer, request, nil)
	if err != nil {
		// TODO: Error handling and socket disposing best practices
		fmt.Println(err)
		return
	}

	for {
		_, messageBuffer, err := socket.ReadMessage()
		if err != nil {
			// TODO: Error handling and socket disposing best practices
			fmt.Println(err)
			return
		}

		// TODO: Implement the handling of the measurements
		fmt.Println(string(messageBuffer))
	}
}
