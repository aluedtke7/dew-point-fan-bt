package sensor

import (
	"testing"
)

func TestNewSensorDataStore(t *testing.T) {
	tests := []struct {
		name       string
		maxData    int
		expectSize int
	}{
		{name: "default minimum limit", maxData: 3, expectSize: 5},
		{name: "exact minimum limit", maxData: 5, expectSize: 5},
		{name: "above minimum limit", maxData: 10, expectSize: 10},
		{name: "zero max data", maxData: 0, expectSize: 5},
		{name: "negative max data", maxData: -1, expectSize: 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := NewSensorDataStore(tt.maxData)
			if store.maxSensorData != tt.expectSize {
				t.Errorf("expected maxSensorData %d, got %d", tt.expectSize, store.maxSensorData)
			}
		})
	}
}

func TestAverageTemperature(t *testing.T) {
	tests := []struct {
		name     string
		data     []SensorData
		expected float64
	}{
		{
			name:     "no data",
			data:     []SensorData{},
			expected: 0,
		},
		{
			name: "single entry",
			data: []SensorData{
				{MacAddress: "MAC1", Temperature: 20.5},
			},
			expected: 20.5,
		},
		{
			name: "multiple entries",
			data: []SensorData{
				{MacAddress: "MAC1", Temperature: 20.5},
				{MacAddress: "MAC2", Temperature: 22.0},
				{MacAddress: "MAC3", Temperature: 25.5},
			},
			expected: 22.7, // Rounded to one decimal place
		},
		{
			name: "all temperatures zero",
			data: []SensorData{
				{MacAddress: "MAC1", Temperature: 0},
				{MacAddress: "MAC2", Temperature: 0},
			},
			expected: 0,
		},
		{
			name: "negative temperatures",
			data: []SensorData{
				{MacAddress: "MAC1", Temperature: -5},
				{MacAddress: "MAC2", Temperature: -10},
			},
			expected: -7.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &SensorDataList{
				data: tt.data,
			}
			result := store.AverageTemperature()
			if result != tt.expected {
				t.Errorf("expected %.1f, got %.1f", tt.expected, result)
			}
		})
	}
}

func TestAddSensorData(t *testing.T) {
	tests := []struct {
		name         string
		initialData  []SensorData
		newData      SensorData
		maxData      int
		expectedSize int
		expectedData []SensorData
	}{
		{
			name:         "add to empty list",
			initialData:  []SensorData{},
			newData:      SensorData{MacAddress: "MAC1", Temperature: 25.5},
			maxData:      5,
			expectedSize: 1,
			expectedData: []SensorData{{MacAddress: "MAC1", Temperature: 25.5}},
		},
		{
			name: "append when not full",
			initialData: []SensorData{
				{MacAddress: "MAC1", Temperature: 20.5},
			},
			newData:      SensorData{MacAddress: "MAC2", Temperature: 22.5},
			maxData:      5,
			expectedSize: 2,
			expectedData: []SensorData{
				{MacAddress: "MAC1", Temperature: 20.5},
				{MacAddress: "MAC2", Temperature: 22.5},
			},
		},
		{
			name: "overwrite oldest when full",
			initialData: []SensorData{
				{MacAddress: "MAC1", Temperature: 20.5},
				{MacAddress: "MAC2", Temperature: 22.5},
				{MacAddress: "MAC3", Temperature: 23.0},
				{MacAddress: "MAC4", Temperature: 24.1},
				{MacAddress: "MAC5", Temperature: 25.2},
			},
			newData:      SensorData{MacAddress: "MAC6", Temperature: 26.5},
			maxData:      5,
			expectedSize: 5,
			expectedData: []SensorData{
				{MacAddress: "MAC2", Temperature: 22.5},
				{MacAddress: "MAC3", Temperature: 23.0},
				{MacAddress: "MAC4", Temperature: 24.1},
				{MacAddress: "MAC5", Temperature: 25.2},
				{MacAddress: "MAC6", Temperature: 26.5},
			},
		},
		{
			name: "add to list with maxData 1",
			initialData: []SensorData{
				{MacAddress: "MAC1", Temperature: 25.0},
			},
			newData:      SensorData{MacAddress: "MAC2", Temperature: 30.0},
			maxData:      1,
			expectedSize: 1,
			expectedData: []SensorData{
				{MacAddress: "MAC2", Temperature: 30.0},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &SensorDataList{
				data:          tt.initialData,
				maxSensorData: tt.maxData,
			}
			store.AddSensorData(tt.newData)
			if store.Size() != tt.expectedSize {
				t.Errorf("expected size %d, got %d", tt.expectedSize, store.Size())
			}
			for i, expected := range tt.expectedData {
				if store.data[i] != expected {
					t.Errorf("expected data at index %d to be %+v, got %+v", i, expected, store.data[i])
				}
			}
		})
	}
}

func TestAverageHumidity(t *testing.T) {
	tests := []struct {
		name     string
		data     []SensorData
		expected float64
	}{
		{
			name:     "no data",
			data:     []SensorData{},
			expected: 0,
		},
		{
			name: "single entry",
			data: []SensorData{
				{MacAddress: "MAC1", Humidity: 55.5},
			},
			expected: 55.5,
		},
		{
			name: "multiple entries",
			data: []SensorData{
				{MacAddress: "MAC1", Humidity: 50.0},
				{MacAddress: "MAC2", Humidity: 60.0},
				{MacAddress: "MAC3", Humidity: 65.0},
			},
			expected: 58.3, // Rounded to one decimal place
		},
		{
			name: "all humidities zero",
			data: []SensorData{
				{MacAddress: "MAC1", Humidity: 0},
				{MacAddress: "MAC2", Humidity: 0},
			},
			expected: 0,
		},
		{
			name: "negative humidities",
			data: []SensorData{
				{MacAddress: "MAC1", Humidity: -10.0},
				{MacAddress: "MAC2", Humidity: -20.0},
			},
			expected: -15.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &SensorDataList{
				data: tt.data,
			}
			result := store.AverageHumidity()
			if result != tt.expected {
				t.Errorf("expected %.1f, got %.1f", tt.expected, result)
			}
		})
	}
}

func TestAverageDewPoint(t *testing.T) {
	tests := []struct {
		name     string
		data     []SensorData
		expected float64
	}{
		{
			name:     "no data",
			data:     []SensorData{},
			expected: 0,
		},
		{
			name: "single entry",
			data: []SensorData{
				{MacAddress: "MAC1", DewPoint: 10.5},
			},
			expected: 10.5,
		},
		{
			name: "multiple entries",
			data: []SensorData{
				{MacAddress: "MAC1", DewPoint: 12.5},
				{MacAddress: "MAC2", DewPoint: 15.0},
				{MacAddress: "MAC3", DewPoint: 14.0},
			},
			expected: 13.8, // Rounded to one decimal place
		},
		{
			name: "all dew points zero",
			data: []SensorData{
				{MacAddress: "MAC1", DewPoint: 0},
				{MacAddress: "MAC2", DewPoint: 0},
			},
			expected: 0,
		},
		{
			name: "negative dew points",
			data: []SensorData{
				{MacAddress: "MAC1", DewPoint: -5.0},
				{MacAddress: "MAC2", DewPoint: -7.5},
			},
			expected: -6.3, // Rounded to one decimal place
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &SensorDataList{
				data: tt.data,
			}
			result := store.AverageDewPoint()
			if result != tt.expected {
				t.Errorf("expected %.1f, got %.1f", tt.expected, result)
			}
		})
	}
}
