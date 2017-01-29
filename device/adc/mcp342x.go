package adc

import (
	"fmt"
	"log"

	"github.com/golang/exp/io/i2c"
)

const MCP3424_ADDRESS = 0x6a
const MCP3424_DEFAULT_BUS = 0x1

type MCP342XResolution byte
type MCP342XGain byte
type MCP342XChannel byte
type MCP342XMode byte

const (
	MCP342X_CHANNEL_1 MCP342XChannel = 0x00
	MCP342X_CHANNEL_2                = 0x20
	MCP342X_CHANNEL_3                = 0x40
	MCP342X_CHANNEL_4                = 0x60

	MCP342X_GAIN_X1 MCP342XGain = 0x00
	MCP342X_GAIN_X2             = 0x04
	MCP342X_GAIN_X3             = 0x08
	MCP342X_GAIN_X4             = 0x0C

	MCP342X_12_BIT MCP342XResolution = 0x00
	MCP342X_14_BIT                   = 0x04
	MCP342X_16_BIT                   = 0x08
	MCP342X_18_BIT                   = 0x0C

	MCP342X_MODE_ONE_SHOOT  MCP342XMode = 0x00
	MCP342X_MODE_CONTINUOUS             = 0x10

	MCP342X_START byte = 0X80 // write: start a conversion
)

type MCP342X struct {
	address int // I2C address
	bus     int // dev-x (bus)
	value   int // Last measured result, in uV

	gain       MCP342XGain       // Gain
	resolution MCP342XResolution // ADC resolution

	config byte

	i2c *i2c.Device
}

func (this *MCP342X) Init(mode MCP342XMode) (err error) {
	this.i2c, err = i2c.Open(&i2c.Devfs{Dev: fmt.Sprintf("/dev/i2c-%d", this.bus)}, this.address)
	if err != nil {
		return err
	}

	this.config = byte(mode) |
		this.getResolution(this.resolution) |
		this.getGain(this.gain)

	err = this.writeConfig(this.config)
	if err != nil {
		return err
	}

	return nil
}

func (this *MCP342X) Close() {
	this.i2c.Close()
}

func (this *MCP342X) getResolution(resolution MCP342XResolution) byte {
	return byte(resolution & 0x0C)
}

func (this *MCP342X) getGain(gain MCP342XGain) byte {
	return byte(gain & 0x03)
}

func (this *MCP342X) getChannel(channel MCP342XChannel) byte {
	return byte(channel & 0x60)
}

func (this *MCP342X) writeConfig(value byte) error {
	err := this.i2c.Write([]byte{value})
	if err != nil {
		return err
	}

	log.Printf("Written config (%d) to the device.\n", value)
	return nil
}

func (this *MCP342X) StartConversion(channel MCP342XChannel) error {
	log.Println("Starting conversion.")

	config := MCP342X_START |
		this.getChannel(channel) |
		this.config

	err := this.writeConfig(config)
	if err != nil {
		return err
	}

	return nil
}

// Returns measured value in mV
func (this *MCP342X) GetMeasurement() (int, error) {
	log.Print("Starting polling for the result.")
	result := make([]byte, 3)

	for {
		err := this.i2c.Read(result)
		if err != nil {
			return 0, err
		}

		if (result[2] & MCP342X_START) == 0x00 {
			break
		}
	}

	// TODO: support all resolutions, it is only for 12 bit
	measurement := ((int(result[0]) & 0x3F) << 8) | int(result[1])

	if measurement > 2048-1 {
		measurement = measurement - 4096 - 1
	}

	return measurement, nil
}

func NewMCP342X(address, bus int, gain MCP342XGain, resolution MCP342XResolution) *MCP342X {
	return &MCP342X{address: address, bus: bus, gain: gain, resolution: resolution}
}
