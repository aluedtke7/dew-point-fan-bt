package bluetooth

import (
	"dpf-bt/sensor"
	"dpf-bt/utility"
	"encoding/binary"
	"fmt"
	"github.com/d2r2/go-logger"
	"slices"
	"strings"
	"time"
	bt "tinygo.org/x/bluetooth"
)

var lg = logger.NewPackageLogger("bt", logger.InfoLevel)

// ProcessAdvertisement processes the Bluetooth advertisement payload to extract sensor data and update the
// Sensors object. It determines whether the data belongs to the inside or outside sensor and logs the
// sensor's details if valid.
func ProcessAdvertisement(scanResult bt.ScanResult,
	sensors *sensor.Sensors,
	sensorStore *sensor.SensorStore) {
	payload := scanResult.AdvertisementPayload.ManufacturerData()[0].Data
	if len(payload) == 18 {
		sensorData := parseWS02Data(payload, scanResult.RSSI, *sensors)
		if sensorData.Name != "" {
			if sensorData.Name == "Inside" {
				sensors.InsideData = *sensorData
				sensorStore.Inside.AddSensorData(*sensorData)
			} else {
				sensors.OutsideData = *sensorData
				sensorStore.Outside.AddSensorData(*sensorData)
			}
			lg.Infof("%8s Temp: %.1f°C - Hum: %.1f%% - Bat: %d - RSSI: %d - Uptime: %s",
				sensorData.Name, sensorData.Temperature, sensorData.Humidity, sensorData.BatLevel,
				sensorData.RSSI, formatUptime(sensorData.Uptime))
		}
	}
}

// parseWS02Data parses WS02 sensor advertisement payload to extract sensor data including temperature,
// humidity, signal strength, and uptime.
// It also calculates calibrated values using sensor-specific offsets and formats the identifier and name
// of the sensor. The input includes the raw payload, RSSI, and a sensor calibration object, and the function
// returns structured sensor data.
func parseWS02Data(payload []byte, rssi int16, sensors sensor.Sensors) *sensor.SensorData {
	// The WS02 advertisement contains temperature and humidity in specific locations.
	// The mac address starts at offset 2, and the 16-bit value of the battery level starts at offset 8.
	// The temperature is a 16-bit value starting at offset 10, and humidity is a 16-bit value starting at offset 12.
	// The uptime in seconds since the last reset is a 32-bit value starting at offset 14.
	const macOffset = 2
	const batOffset = 8
	const tempOffset = 10
	const humidityOffset = 12
	const uptimeOffset = 14

	macAdr := ""
	for _, c := range slices.Backward(payload[macOffset : macOffset+6]) {
		macAdr = macAdr + fmt.Sprintf("%02X:", c)
	}
	macAdr = strings.TrimSuffix(macAdr, ":")
	name := ""
	tempCal := 0.0
	humCal := 0.0
	if macAdr == sensors.InsideData.MacAddress {
		name = "Inside"
		tempCal = sensors.InsideCalibration.Temperature
		humCal = sensors.InsideCalibration.Humidity
	} else if macAdr == sensors.OutsideData.MacAddress {
		name = "Outside"
		tempCal = sensors.OutsideCalibration.Temperature
		humCal = sensors.OutsideCalibration.Humidity
	}
	batLevel := binary.LittleEndian.Uint16(payload[batOffset : batOffset+2])
	uptime := binary.LittleEndian.Uint32(payload[uptimeOffset : uptimeOffset+4])

	temperatureInt := binary.LittleEndian.Uint16(payload[tempOffset : tempOffset+2])
	temperatureRaw := float64(temperatureInt) / 16.0
	if temperatureRaw > 4000 {
		temperatureRaw -= 4096
	}
	temperature := temperatureRaw + tempCal

	humidityInt := binary.LittleEndian.Uint16(payload[humidityOffset : humidityOffset+2])
	humidityRaw := float64(humidityInt) / 16.0
	if humidityRaw > 4000 {
		humidityRaw -= 4096
	}
	humidity := humidityRaw + humCal

	roundedTemperature := utility.RoundDouble(temperature, 1)
	roundedHumidity := utility.RoundDouble(humidity, 1)

	return &sensor.SensorData{
		MacAddress:  macAdr,
		Name:        name,
		BatLevel:    batLevel,
		RSSI:        rssi,
		Uptime:      uptime,
		Temperature: roundedTemperature,
		Humidity:    roundedHumidity,
		DewPoint:    utility.CalcDewPoint(roundedTemperature, roundedHumidity),
		Scanned:     time.Now(),
	}
}

// formatUptime converts uptime in seconds to a human-readable format as a string in the form "Xd Yh Zm".
func formatUptime(seconds uint32) string {
	days := seconds / (24 * 3600)
	hours := (seconds % (24 * 3600)) / 3600
	minutes := (seconds % 3600) / 60
	return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
}
