package bluetooth

import (
	"dpf-bt/sensor"
	"dpf-bt/utility"
	"encoding/binary"
	"testing"
	"time"
)

func TestParseWS02Data(t *testing.T) {
	tests := []struct {
		name           string
		payload        []byte
		rssi           int16
		sensors        sensor.Sensors
		expectedResult *sensor.SensorData
	}{
		{
			name: "Valid inside sensor data",
			payload: func() []byte {
				payload := make([]byte, 18)
				binary.LittleEndian.PutUint16(payload[8:10], 72)    // battery level
				binary.LittleEndian.PutUint32(payload[14:18], 3600) // uptime
				binary.LittleEndian.PutUint16(payload[10:12], 1280) // temperature
				binary.LittleEndian.PutUint16(payload[12:14], 960)  // humidity
				copy(payload[2:8], []byte{0x12, 0x34, 0x56, 0x78, 0x9A, 0xBC})
				return payload
			}(),
			rssi: -50,
			sensors: sensor.Sensors{
				InsideData: sensor.SensorData{
					MacAddress: "BC:9A:78:56:34:12",
				},
				InsideCalibration: sensor.SensorCalibration{
					Temperature: 0.5,
					Humidity:    1.0,
				},
			},
			expectedResult: &sensor.SensorData{
				MacAddress:  "BC:9A:78:56:34:12",
				Name:        "Inside",
				BatLevel:    72,
				RSSI:        -50,
				Uptime:      3600,
				Temperature: 80.5,
				Humidity:    61.0,
				DewPoint:    utility.CalcDewPoint(80.5, 61.0),
				Scanned:     time.Now(), // Should not be compared directly in test
			},
		},
		{
			name: "Valid outside sensor data",
			payload: func() []byte {
				payload := make([]byte, 18)
				binary.LittleEndian.PutUint16(payload[8:10], 55)    // battery level
				binary.LittleEndian.PutUint32(payload[14:18], 7200) // uptime
				binary.LittleEndian.PutUint16(payload[10:12], 1500) // temperature
				binary.LittleEndian.PutUint16(payload[12:14], 500)  // humidity
				copy(payload[2:8], []byte{0xAB, 0xCD, 0xEF, 0x12, 0x34, 0x56})
				return payload
			}(),
			rssi: -60,
			sensors: sensor.Sensors{
				OutsideData: sensor.SensorData{
					MacAddress: "56:34:12:EF:CD:AB",
				},
				OutsideCalibration: sensor.SensorCalibration{
					Temperature: -0.5,
					Humidity:    -1.0,
				},
			},
			expectedResult: &sensor.SensorData{
				MacAddress:  "56:34:12:EF:CD:AB",
				Name:        "Outside",
				BatLevel:    55,
				RSSI:        -60,
				Uptime:      7200,
				Temperature: 93.3,
				Humidity:    30.3,
				DewPoint:    utility.CalcDewPoint(93.3, 30.3),
				Scanned:     time.Now(), // Should not be compared directly in test
			},
		},
		{
			name: "Invalid MAC, no matching sensor",
			payload: func() []byte {
				payload := make([]byte, 18)
				binary.LittleEndian.PutUint16(payload[8:10], 100)   // battery level
				binary.LittleEndian.PutUint32(payload[14:18], 1800) // uptime
				binary.LittleEndian.PutUint16(payload[10:12], 2000) // temperature
				binary.LittleEndian.PutUint16(payload[12:14], 3000) // humidity
				copy(payload[2:8], []byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xAB})
				return payload
			}(),
			rssi: -70,
			sensors: sensor.Sensors{
				InsideData: sensor.SensorData{
					MacAddress: "BC:9A:78:56:34:12",
				},
				OutsideData: sensor.SensorData{
					MacAddress: "56:34:12:EF:CD:AB",
				},
			},
			expectedResult: &sensor.SensorData{
				MacAddress:  "AB:89:67:45:23:01",
				Name:        "",
				BatLevel:    100,
				RSSI:        -70,
				Uptime:      1800,
				Temperature: 125.0,
				Humidity:    187.5,
				DewPoint:    utility.CalcDewPoint(125.0, 187.5),
				Scanned:     time.Now(), // Should not be compared directly in test
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseWS02Data(tt.payload, tt.rssi, tt.sensors)

			if result.MacAddress != tt.expectedResult.MacAddress {
				t.Errorf("MacAddress = %v, want %v", result.MacAddress, tt.expectedResult.MacAddress)
			}
			if result.Name != tt.expectedResult.Name {
				t.Errorf("Name = %v, want %v", result.Name, tt.expectedResult.Name)
			}
			if result.BatLevel != tt.expectedResult.BatLevel {
				t.Errorf("BatLevel = %v, want %v", result.BatLevel, tt.expectedResult.BatLevel)
			}
			if result.RSSI != tt.expectedResult.RSSI {
				t.Errorf("RSSI = %v, want %v", result.RSSI, tt.expectedResult.RSSI)
			}
			if result.Uptime != tt.expectedResult.Uptime {
				t.Errorf("Uptime = %v, want %v", result.Uptime, tt.expectedResult.Uptime)
			}
			if result.Temperature != tt.expectedResult.Temperature {
				t.Errorf("Temperature = %v, want %v", result.Temperature, tt.expectedResult.Temperature)
			}
			if result.Humidity != tt.expectedResult.Humidity {
				t.Errorf("Humidity = %v, want %v", result.Humidity, tt.expectedResult.Humidity)
			}
			if result.DewPoint != tt.expectedResult.DewPoint {
				t.Errorf("DewPoint = %v, want %v", result.DewPoint, tt.expectedResult.DewPoint)
			}
		})
	}
}
