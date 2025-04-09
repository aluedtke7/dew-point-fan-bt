package sensor

import "time"

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

type SensorCalibration struct {
	Temperature float64
	Humidity    float64
}

type Sensors struct {
	InsideData         SensorData
	InsideCalibration  SensorCalibration
	OutsideData        SensorData
	OutsideCalibration SensorCalibration
}

type FanConfig struct {
	MinDiff           float64
	Hysteresis        float64
	MinHumidityInside float64
	MinTempInside     float64
	MinTempOutside    float64
}

type Reason int

var ReasonName = map[Reason]string{
	ReasonNone:                 "none",
	ReasonNoData:               "no data",
	ReasonNoEnoughData:         "not enough data",
	ReasonDewPointOverHyst:     "dp over hyst",
	ReasonDewPointUnderHyst:    "dp under hyst",
	ReasonInsideTempTooLow:     "in temp too low",
	ReasonOutsideTempTooLow:    "out temp too low",
	ReasonInsideHumidityTooLow: "in hum too low",
}

const (
	ReasonNone Reason = iota
	ReasonNoData
	ReasonNoEnoughData
	ReasonDewPointOverHyst
	ReasonDewPointUnderHyst
	ReasonInsideTempTooLow
	ReasonOutsideTempTooLow
	ReasonInsideHumidityTooLow
)

type ResultData struct {
	DpDiff     float64
	ShouldBeOn bool
	IsOn       bool
	Outcome    Reason
}
