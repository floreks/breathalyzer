package main

import (
	"github.com/CarlosRA97/wpi"
	"fmt"
)

const deviceID = 0x6e

func main() {
	err := wpi.I2CSetup(deviceID)
	fmt.Printf(err)
}