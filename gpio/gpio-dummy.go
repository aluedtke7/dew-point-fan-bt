//go:build !linux && !arm

package gpio

var lgTerm = logger.NewPackageLogger("gpio", logger.InfoLevel)

type gpioDummyData struct {
	fanState bool
}

func (g gpioDummyData) ReadFanSense() bool {
	return false
}

func (g gpioDummyData) SetFan(on bool) {
	if on != g.fanState {
		lg.infof("Switching Fan to %v", on)
	}
	g.fanState = on
}

func New() (_ Gpio, err error) {
	err = nil
	gpio := &gpioDummyData{fanState: false}

	return *gpio, err
}
