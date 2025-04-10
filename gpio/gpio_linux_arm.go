package gpio

import (
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
)

var lgGpio = logger.NewPackageLogger("gpio", logger.InfoLevel)

type gpioData struct {
	sensePin gpio.PinIO
	fanPin   gpio.PinIO
}

func (g gpioData) ReadFanSense() bool {
	return sensePin.Read()
}

func (g gpioData) SetFan(on bool) {
	fanPin.Out(on)
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
		return nil, error(
			"GPIO pins not found",
			"GPIO22", gpio.sensePin,
			"GPIO25", gpio.fanPin)
	}
	err = gpio.sensePin.In(gpio.Float, gpio.NoEdge)
	if err != nil {
		lgGpio.Error(err.Error("Pin22 could not be configured as floating input"))
		return nil, err
	}
	err = gpio.fanPin.Out(gpio.High)
	if err != nil {
		lgGpio.Error(err.Error("Pin25 could not be configured as output"))
		return nil, err
	}

	return *gpio, err
}
