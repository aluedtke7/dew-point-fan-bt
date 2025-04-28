package main

import (
	"dpf-bt/sensor"
	"reflect"
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
		expectedResult sensor.ResultData
		lastResult     sensor.ResultData
	}{
		{
			name:           "RemoteOverrideOn",
			remoteOverride: 1,
			expectedResult: sensor.ResultData{
				ShouldBeOn: true,
				Reason:     sensor.ReasonSoftOverrideOn,
			},
			lastResult: sensor.ResultData{
				ShouldBeOn: false,
				Reason:     sensor.ReasonNone,
			},
		},
		{
			name:           "RemoteOverrideOff",
			remoteOverride: 2,
			expectedResult: sensor.ResultData{
				ShouldBeOn: false,
				Reason:     sensor.ReasonSoftOverrideOff,
			},
			lastResult: sensor.ResultData{
				ShouldBeOn: true,
				Reason:     sensor.ReasonNone,
			},
		},
		{
			name: "NoDataInside",
			inside: sensor.SensorData{
				Scanned: time.Time{},
			},
			outside: sensor.SensorData{
				Scanned: time.Now(),
			},
			expectedResult: sensor.ResultData{
				ShouldBeOn: false,
				Reason:     sensor.ReasonNoData,
			},
			lastResult: sensor.ResultData{
				ShouldBeOn: true,
				Reason:     sensor.ReasonNone,
			},
		},
		{
			name: "NoDataOutside",
			inside: sensor.SensorData{
				Scanned: time.Now(),
			},
			outside: sensor.SensorData{
				Scanned: time.Time{},
			},
			expectedResult: sensor.ResultData{
				ShouldBeOn: false,
				Reason:     sensor.ReasonNoData,
			},
			lastResult: sensor.ResultData{
				ShouldBeOn: true,
				Reason:     sensor.ReasonNone,
			},
		},
		{
			name: "DataTooOld",
			inside: sensor.SensorData{
				Scanned: time.Now().Add(-10 * time.Minute),
			},
			outside: sensor.SensorData{
				Scanned: time.Now(),
			},
			expectedResult: sensor.ResultData{
				ShouldBeOn: false,
				Reason:     sensor.ReasonNoEnoughData,
			},
			lastResult: sensor.ResultData{
				ShouldBeOn: true,
				Reason:     sensor.ReasonNone,
			},
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
			expectedResult: sensor.ResultData{
				ShouldBeOn: false,
				Reason:     sensor.ReasonInsideTempTooLow,
			},
			lastResult: sensor.ResultData{
				ShouldBeOn: true,
				Reason:     sensor.ReasonNone,
			},
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
			expectedResult: sensor.ResultData{
				ShouldBeOn: false,
				Reason:     sensor.ReasonOutsideTempTooLow,
			},
			lastResult: sensor.ResultData{
				ShouldBeOn: true,
				Reason:     sensor.ReasonNone,
			},
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
			expectedResult: sensor.ResultData{
				ShouldBeOn: false,
				Reason:     sensor.ReasonInsideHumidityTooLow,
			},
			lastResult: sensor.ResultData{
				ShouldBeOn: true,
				Reason:     sensor.ReasonNone,
			},
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
			expectedResult: sensor.ResultData{
				ShouldBeOn: false,
				Reason:     sensor.ReasonDewPointUnderHyst,
			},
			lastResult: sensor.ResultData{
				ShouldBeOn: true,
				Reason:     sensor.ReasonNone,
			},
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
			expectedResult: sensor.ResultData{
				ShouldBeOn: true,
				Reason:     sensor.ReasonDewPointOverHyst,
			},
			lastResult: sensor.ResultData{
				ShouldBeOn: false,
				Reason:     sensor.ReasonNone,
			},
		},
		{
			name: "DewPointInBetweenFromLow",
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
			expectedResult: sensor.ResultData{
				ShouldBeOn: false,
				Reason:     sensor.ReasonDewPointInBetween,
			},
			lastResult: sensor.ResultData{
				ShouldBeOn: false,
				Reason:     sensor.ReasonNone,
			},
		},
		{
			name: "DewPointInBetweenFromHigh",
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
			expectedResult: sensor.ResultData{
				ShouldBeOn: true,
				Reason:     sensor.ReasonDewPointInBetween,
			},
			lastResult: sensor.ResultData{
				ShouldBeOn: true,
				Reason:     sensor.ReasonNone,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set global variables used in computeResults
			remoteOverride = tt.remoteOverride
			fanConfig = tt.fanConfig

			computeResults(tt.inside, tt.outside, &tt.lastResult)

			if !reflect.DeepEqual(tt.lastResult, tt.expectedResult) {
				t.Errorf("expected ResultData = %+v, got %+v", tt.expectedResult, tt.lastResult)
			}
		})
	}
}
