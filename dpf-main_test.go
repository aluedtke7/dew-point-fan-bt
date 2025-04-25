package main

import (
	"dpf-bt/sensor"
	"testing"
	"time"
)

func TestComputeResults(t *testing.T) {
	tests := []struct {
		name           string
		inside         sensor.SensorData
		outside        sensor.SensorData
		remoteOverride int
		fanConfig      sensor.FanConfig
		expectedOn     bool
		expectedReason sensor.Reason
	}{
		{
			name:           "RemoteOverrideOn",
			remoteOverride: 1,
			expectedOn:     true,
			expectedReason: sensor.ReasonSoftOverrideOn,
		},
		{
			name:           "RemoteOverrideOff",
			remoteOverride: 2,
			expectedOn:     false,
			expectedReason: sensor.ReasonSoftOverrideOff,
		},
		{
			name: "NoDataInside",
			inside: sensor.SensorData{
				Scanned: time.Time{},
			},
			outside: sensor.SensorData{
				Scanned: time.Now(),
			},
			expectedOn:     false,
			expectedReason: sensor.ReasonNoData,
		},
		{
			name: "NoDataOutside",
			inside: sensor.SensorData{
				Scanned: time.Now(),
			},
			outside: sensor.SensorData{
				Scanned: time.Time{},
			},
			expectedOn:     false,
			expectedReason: sensor.ReasonNoData,
		},
		{
			name: "DataTooOld",
			inside: sensor.SensorData{
				Scanned: time.Now().Add(-10 * time.Minute),
			},
			outside: sensor.SensorData{
				Scanned: time.Now(),
			},
			expectedOn:     false,
			expectedReason: sensor.ReasonNoEnoughData,
		},
		{
			name: "LowInsideTemperature",
			inside: sensor.SensorData{
				Temperature: 18.0,
				Scanned:     time.Now(),
			},
			outside: sensor.SensorData{
				Scanned: time.Now(),
			},
			fanConfig: sensor.FanConfig{
				MinTempInside: 20.0,
			},
			expectedOn:     false,
			expectedReason: sensor.ReasonInsideTempTooLow,
		},
		{
			name: "LowOutsideTemperature",
			inside: sensor.SensorData{
				Temperature: 21.0,
				Scanned:     time.Now(),
			},
			outside: sensor.SensorData{
				Temperature: 7.0,
				Scanned:     time.Now(),
			},
			fanConfig: sensor.FanConfig{
				MinTempOutside: 10.0,
			},
			expectedOn:     false,
			expectedReason: sensor.ReasonOutsideTempTooLow,
		},
		{
			name: "LowInsideHumidity",
			inside: sensor.SensorData{
				Humidity: 45.0,
				Scanned:  time.Now(),
			},
			outside: sensor.SensorData{
				Scanned: time.Now(),
			},
			fanConfig: sensor.FanConfig{
				MinHumidityInside: 50.0,
			},
			expectedOn:     false,
			expectedReason: sensor.ReasonInsideHumidityTooLow,
		},
		{
			name: "DewPointBelowHysteresis",
			inside: sensor.SensorData{
				DewPoint: 10.0,
				Scanned:  time.Now(),
			},
			outside: sensor.SensorData{
				DewPoint: 7.0,
				Scanned:  time.Now(),
			},
			fanConfig: sensor.FanConfig{
				MinDiff: 4.0,
			},
			expectedOn:     false,
			expectedReason: sensor.ReasonDewPointUnderHyst,
		},
		{
			name: "DewPointAboveHysteresis",
			inside: sensor.SensorData{
				DewPoint: 15.0,
				Scanned:  time.Now(),
			},
			outside: sensor.SensorData{
				DewPoint: 9.0,
				Scanned:  time.Now(),
			},
			fanConfig: sensor.FanConfig{
				MinDiff:    4.0,
				Hysteresis: 2.0,
			},
			expectedOn:     true,
			expectedReason: sensor.ReasonDewPointOverHyst,
		},
		{
			name: "DewPointInBetween",
			inside: sensor.SensorData{
				DewPoint: 13.0,
				Scanned:  time.Now(),
			},
			outside: sensor.SensorData{
				DewPoint: 9.0,
				Scanned:  time.Now(),
			},
			fanConfig: sensor.FanConfig{
				MinDiff:    4.0,
				Hysteresis: 2.0,
			},
			expectedOn:     false,
			expectedReason: sensor.ReasonDewPointInBetween,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set global variables used in computeResults
			remoteOverride = tt.remoteOverride
			fanConfig = tt.fanConfig

			var resultData sensor.ResultData
			computeResults(tt.inside, tt.outside, &resultData)

			if resultData.ShouldBeOn != tt.expectedOn {
				t.Errorf("expected ShouldBeOn = %v, got %v", tt.expectedOn, resultData.ShouldBeOn)
			}
			if resultData.Reason != tt.expectedReason {
				t.Errorf("expected Reason = %v, got %v", tt.expectedReason, resultData.Reason)
			}
		})
	}
}
