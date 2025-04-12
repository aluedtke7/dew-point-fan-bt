package gpio

type Gpio interface {
	ReadFanSense() bool
	SetFan(on bool)
}
