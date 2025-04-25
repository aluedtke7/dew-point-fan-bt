package main

import (
	"context"
	influxdb "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"time"
)

const (
	tickInterval       = time.Minute
	minRequiredSamples = 10
	measurementName    = "dp"
)

// sendToInfluxDb sends aggregated sensor data to InfluxDB at scheduled intervals
func sendToInfluxDb() {
	client := influxdb.NewClient(influxConfig.Url, influxConfig.Token)
	writeAPI := client.WriteAPIBlocking(influxConfig.Org, influxConfig.Bucket)
	tags := make(map[string]string)

	ticker := time.NewTicker(tickInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if !hasEnoughData() {
				logInsufficientData()
				continue
			}
			logDataTransmissionStart()
			point := createDataPoint(tags)
			if err := writeAPI.WritePoint(context.Background(), point); err != nil {
				lg.Error(err)
				continue
			}
			logAverageValues()
		}
	}
}

// hasEnoughData checks if both inside and outside sensor data lists have at least the minimum required
// number of samples.
func hasEnoughData() bool {
	return sensorStore.Inside.Size() >= minRequiredSamples &&
		sensorStore.Outside.Size() >= minRequiredSamples
}

// logInsufficientData logs a warning when sensor data is not enough for sending to InfluxDB.
func logInsufficientData() {
	lg.Warnf("NOT sending to InfluxDB due to insufficient data (Inside/Outside): %d, %d",
		sensorStore.Inside.Size(), sensorStore.Outside.Size())
}

// logDataTransmissionStart logs the start of data transmission to InfluxDB, including the size of inside
// and outside data lists.
func logDataTransmissionStart() {
	lg.Infof("Sending average values to InfluxDB (Inside/Outside): %d, %d",
		sensorStore.Inside.Size(), sensorStore.Outside.Size())
}

// createDataPoint generates a data point with sensor readings and additional metadata for InfluxDB storage.
func createDataPoint(tags map[string]string) *write.Point {
	ventingValue := 0
	if resultData.IsOn {
		ventingValue = 1
	}

	fields := map[string]interface{}{
		"temp_i":     sensorStore.Inside.AverageTemperature(),
		"temp_o":     sensorStore.Outside.AverageTemperature(),
		"dewpoint_i": sensorStore.Inside.AverageDewPoint(),
		"dewpoint_o": sensorStore.Outside.AverageDewPoint(),
		"hum_i":      sensorStore.Inside.AverageHumidity(),
		"hum_o":      sensorStore.Outside.AverageHumidity(),
		"retry_i":    0,
		"retry_o":    0,
		"vent_val":   ventingValue,
	}
	return write.NewPoint(measurementName, tags, fields, time.Now())
}

// logAverageValues logs the average temperature and humidity for both inside and outside sensor data stores.
func logAverageValues() {
	lg.Infof("Inside  (T/H): %5.1fC - %5.1f%%",
		sensorStore.Inside.AverageTemperature(),
		sensorStore.Inside.AverageHumidity())
	lg.Infof("Outside (T/H): %5.1fC - %5.1f%%",
		sensorStore.Outside.AverageTemperature(),
		sensorStore.Outside.AverageHumidity())
}
