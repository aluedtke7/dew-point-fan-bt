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

type FanData struct {
	ShouldBeOn bool
	IsOn       bool
	Reason     string
}
