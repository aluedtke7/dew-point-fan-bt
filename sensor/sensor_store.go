package sensor

// SensorDataList struct manages a list of up to 20 SensorData entries.
type SensorDataList struct {
	data          []SensorData
	maxSensorData int
}

// SensorStore struct manages collections of sensor data for inside and outside environments.
type SensorStore struct {
	Inside  SensorDataList
	Outside SensorDataList
}

// NewSensorDataStore initializes a SensorDataList with a specified maximum capacity for storing sensor data.
// The minimum allowable maxData value is 5; lower values default to 5.
// Returns a pointer to the initialized SensorDataList.
func NewSensorDataStore(maxData int) *SensorDataList {
	if maxData < 5 {
		// Ensure a valid maximum value; fallback to 5 if invalid input is provided.
		maxData = 5
	}
	return &SensorDataList{
		data:          []SensorData{},
		maxSensorData: maxData,
	}
}

// AddSensorData adds a new SensorData to the store. It removes the oldest entry if the limit is exceeded.
func (store *SensorDataList) AddSensorData(sensor SensorData) {
	if len(store.data) >= store.maxSensorData {
		// If the list exceeds the limit, remove the oldest entry
		store.data = store.data[1:]
	}
	// Add the new sensor data to the list
	store.data = append(store.data, sensor)
}

// AverageTemperature calculates the average temperature from all SensorData entries in the store.
// Returns 0 if there are no SensorData entries.
func (store *SensorDataList) AverageTemperature() float64 {
	if len(store.data) == 0 {
		return 0
	}

	var totalTemperature float64
	for _, sensor := range store.data {
		totalTemperature += sensor.Temperature
	}
	return totalTemperature / float64(len(store.data))
}

// AverageHumidity calculates the average humidity of all SensorData entries in the store.
// Returns 0 if there are no SensorData entries.
func (store *SensorDataList) AverageHumidity() float64 {
	if len(store.data) == 0 {
		return 0
	}

	var totalHumidity float64
	for _, sensor := range store.data {
		totalHumidity += sensor.Humidity
	}
	return totalHumidity / float64(len(store.data))
}

// AverageDewPoint calculates the average dew point of all SensorData in the store.
// Returns 0 if there are no SensorData entries.
func (store *SensorDataList) AverageDewPoint() float64 {
	if len(store.data) == 0 {
		return 0
	}

	var totalDewPoint float64
	for _, sensor := range store.data {
		totalDewPoint += sensor.DewPoint
	}
	return totalDewPoint / float64(len(store.data))
}

func (store *SensorDataList) Size() int {
	return len(store.data)
}
