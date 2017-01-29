package sensor

import (
	"log"
	"time"

	rpi "github.com/stianeikeland/go-rpio"
)

const (
	detectionTimeout  = 5 * time.Second
	detectionInterval = 100 * time.Millisecond
)

type MQ3Response struct {
	AlcoholDetected bool `json:"alcoholDetected"`
}

type MQReader interface {
	DetectAlcohol() MQ3Response
}

type MQ3Reader struct {
	detectionTimeout  time.Duration
	detectionInterval time.Duration
	gpioPin           int
}

func (this *MQ3Reader) init() error {
	this.detectionInterval = detectionInterval
	this.detectionTimeout = detectionTimeout

	return rpi.Open()
}

// DetectAlcohol reads for 5 sec digital out from mq-3 sensor and returns true or false based on readings over time
func (this MQ3Reader) DetectAlcohol() MQ3Response {
	log.Printf("Starting alcohol detection. Detection time: %2.f seconds", this.detectionTimeout.Seconds())

	timeoutChan := time.NewTimer(this.detectionTimeout).C
	intervalChan := time.NewTicker(this.detectionInterval).C
	var result int
	tmp := 1.0
	readings := this.detectionTimeout.Seconds() / this.detectionInterval.Seconds()

	pin := rpi.Pin(this.gpioPin)

	// Set input mode to read from the pin
	pin.Input()

	for {
		select {
		case <-timeoutChan:
			tmp = tmp / float64(readings)
			result = int(tmp + 0.5)

			log.Printf("Finished alcohol detection. Result: %d", result)

			if result == 0 {
				return MQ3Response{AlcoholDetected: false}
			}

			return MQ3Response{AlcoholDetected: true}
		case <-intervalChan:
			out := pin.Read()
			if out == rpi.Low {
				tmp += 1.0
			}
		}
	}
}

func NewMQ3Reader(gpioPin int) (MQReader, error) {
	reader := &MQ3Reader{gpioPin: gpioPin}
	err := reader.init()
	if err != nil {
		return nil, err
	}

	return reader, nil
}
