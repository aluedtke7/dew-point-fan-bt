package gpio

type gpioDummyData struct {
	fanState bool
}

func (g gpioDummyData) ReadFanSense() bool {
	return false
}

func (g gpioDummyData) SetFan(on bool) {
	g.fanState = on
}

func New() (_ Gpio, err error) {
	err = nil
	gpio := &gpioDummyData{fanState: false}

	return *gpio, err
}
