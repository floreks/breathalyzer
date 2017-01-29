package service

import (
	"net/http"

	"github.com/emicklei/go-restful"
	"github.com/floreks/breathalyzer/device/sensor"
)

const defaultSensorPin = 5

// TODO add doc
type MQ3Service struct {
	reader sensor.MQReader
}

// TODO add doc
func (d MQ3Service) Handler() *restful.WebService {
	ws := new(restful.WebService)
	ws.
		Path("/mq3").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/measure").To(d.detectAlcohol).
		Doc("Reads temperature from DS18B20 sensor").
		Writes(sensor.MQ3Response{}))

	return ws
}

// TODO add doc
func (d MQ3Service) detectAlcohol(request *restful.Request, response *restful.Response) {
	result := d.reader.DetectAlcohol()
	response.WriteHeaderAndEntity(http.StatusOK, result)
}

// TODO add doc
func NewMQ3Service() (*MQ3Service, error) {
	reader, err := sensor.NewMQ3Reader(defaultSensorPin)
	if err != nil {
		return nil, err
	}

	return &MQ3Service{reader: reader}, nil
}
