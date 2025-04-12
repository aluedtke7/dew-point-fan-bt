//go:build !linux && !arm

package gpio

import "github.com/d2r2/go-logger"

var lgGp = logger.NewPackageLogger("gpio", logger.InfoLevel)

type gpioDummyData struct {
	fanState bool
}

func (g gpioDummyData) ReadFanSense() bool {
	return false
}

func (g gpioDummyData) SetFan(on bool) {
	if on != g.fanState {
		lgGp.Infof("Switching Fan to %v", on)
	}
	g.fanState = on
}

func New() (_ Gpio, err error) {
	err = nil
	gpio := &gpioDummyData{fanState: false}

	return *gpio, err
}
