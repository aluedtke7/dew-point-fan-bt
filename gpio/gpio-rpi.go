//go:build linux && arm

package gpio

import (
	"errors"
	"github.com/d2r2/go-logger"
	gp "periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/host/v3"
)

var lgGpio = logger.NewPackageLogger("gpio", logger.InfoLevel)

type gpioData struct {
	sensePin gp.PinIO
	fanPin   gp.PinIO
}

func (g gpioData) ReadFanSense() bool {
	return g.sensePin.Read() == gp.High
}

func (g gpioData) SetFan(on bool) {
	if on {
		_ = g.fanPin.Out(gp.High)
	} else {
		_ = g.fanPin.Out(gp.Low)
	}
}

func New() (_ Gpio, err error) {
	_, err = host.Init()
	if err != nil {
		return nil, err
	}
	gpio := &gpioData{
		sensePin: gpioreg.ByName("GPIO22"),
		fanPin:   gpioreg.ByName("GPIO25"),
	}
	if gpio.sensePin == nil || gpio.fanPin == nil {
		lgGpio.Error("GPIO pins not found")
		return nil, errors.New("GPIO pins not found")
	}
	err = gpio.sensePin.In(gp.Float, gp.NoEdge)
	if err != nil {
		lgGpio.Error("Pin22 could not be configured as floating input")
		return nil, err
	}
	err = gpio.fanPin.Out(gp.High)
	if err != nil {
		lgGpio.Error("Pin25 could not be configured as output")
		return nil, err
	}

	return *gpio, err
}
