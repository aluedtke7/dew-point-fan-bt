package gpio

// Gpio defines an interface for interacting with general-purpose input/output (GPIO) pins.
// It provides methods for reading the fan sense state and controlling the fan power state.
type Gpio interface {

	// ReadFanSense reads the current state of the fan sensor and returns true if the fan is detected as running.
	ReadFanSense() bool

	// SetFan controls the power state of the fan by turning it on or off based on the provided boolean value.
	SetFan(on bool)
}
