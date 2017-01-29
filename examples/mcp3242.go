package main

import (
	"log"

	"github.com/floreks/breathalyzer/device/adc"
)

func main() {
	mcp3424 := adc.NewMCP342X(adc.MCP3424_ADDRESS, adc.MCP3424_DEFAULT_BUS, adc.MCP342X_GAIN_X1, adc.MCP342X_12_BIT)

	err := mcp3424.Init(adc.MCP342X_MODE_ONE_SHOOT)
	if err != nil {
		log.Fatalf(err.Error())
	}

	defer mcp3424.Close()
	log.Println("Device initialized.")

	err = mcp3424.StartConversion(adc.MCP342X_CHANNEL_4)
	if err != nil {
		log.Fatalf(err.Error())
	}

	mV, err := mcp3424.GetMeasurement()
	if err != nil {
		log.Fatalf(err.Error())
	}

	log.Printf("Measured: %d mV", mV)
}
