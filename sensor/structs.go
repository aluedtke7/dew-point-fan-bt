package sensor

import "time"

// SensorData represents data collected from a sensor, including environmental measurements and metadata.
type SensorData struct {
	MacAddress  string
	Name        string
	BatLevel    uint16
	RSSI        int16
	Uptime      uint32
	Temperature float64
	Humidity    float64
	DewPoint    float64
	Scanned     time.Time
}

// SensorCalibration represents calibration data for sensors, including adjustments for temperature and humidity.
type SensorCalibration struct {
	Temperature float64
	Humidity    float64
}

// Sensors represent a collection of sensor data and calibration for both inside and outside environments.
type Sensors struct {
	InsideData         SensorData
	InsideCalibration  SensorCalibration
	OutsideData        SensorData
	OutsideCalibration SensorCalibration
}

// FanConfig is a configuration structure for controlling fan behavior based on environmental parameters and thresholds.
type FanConfig struct {
	MinDiff           float64
	Hysteresis        float64
	MinHumidityInside float64
	MinTempInside     float64
	MinTempOutside    float64
}

// Reason represents a categorized outcome or state as an integer constant.
type Reason int

const (
	// ReasonNone specifies no reason or default state.
	ReasonNone Reason = iota

	// ReasonNoData indicates the absence of data.
	ReasonNoData

	// ReasonNoEnoughData indicates insufficient data.
	ReasonNoEnoughData

	// ReasonDewPointOverHyst indicates the dew point is over the hysteresis limit.
	ReasonDewPointOverHyst

	// ReasonDewPointUnderHyst indicates the dew point is under the hysteresis limit.
	ReasonDewPointUnderHyst

	// ReasonDewPointInBetween indicates the dew point is within the hysteresis range.
	ReasonDewPointInBetween

	// ReasonInsideTempTooLow indicates the inside temperature is too low.
	ReasonInsideTempTooLow

	// ReasonOutsideTempTooLow indicates the outside temperature is too low.
	ReasonOutsideTempTooLow

	// ReasonInsideHumidityTooLow indicates the inside humidity is too low.
	ReasonInsideHumidityTooLow

	// ReasonSoftOverrideOn indicates a state where an override is set to 'on' via REST api.
	ReasonSoftOverrideOn

	// ReasonSoftOverrideOff indicates a state where an override is set to 'off' via REST API.
	ReasonSoftOverrideOff

	// ReasonUnknown indicates an undefined or unknown reason.
	ReasonUnknown
)

// ReasonName maps Reason constants to their corresponding string representations for descriptive purposes.
var ReasonName = map[Reason]string{
	ReasonNone:                 "none",
	ReasonNoData:               "no data",
	ReasonNoEnoughData:         "not enough data",
	ReasonDewPointOverHyst:     "dp > hysteresis",
	ReasonDewPointUnderHyst:    "dp < hysteresis",
	ReasonDewPointInBetween:    "dp in between",
	ReasonInsideTempTooLow:     "inside temp too low",
	ReasonOutsideTempTooLow:    "outside temp too low",
	ReasonInsideHumidityTooLow: "inside hum too low",
	ReasonSoftOverrideOn:       "soft override on",
	ReasonSoftOverrideOff:      "soft override off",
	ReasonUnknown:              "unknown reason",
}

// ResultData represents the computational output for fan control based on sensor data and configuration thresholds.
type ResultData struct {
	DpDiff     float64
	ShouldBeOn bool
	IsOn       bool
	Reason     Reason
}

// InfluxDbConfig represents the configuration settings for connecting to an InfluxDB instance.
type InfluxDbConfig struct {
	Enabled bool
	Url     string
	Token   string
	Org     string
	Bucket  string
}
