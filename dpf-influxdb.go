package main

import (
	"context"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"time"
)

func sendToInfluxDb() {
	client := influxdb2.NewClient(influxConfig.Url, influxConfig.Token)
	writeAPI := client.WriteAPIBlocking(influxConfig.Org, influxConfig.Bucket)
	tags := map[string]string{
		// "manual_override": strconv.FormatBool(fanStatus),
		// "remote_override": strconv.Itoa(remoteOverride),
		// "venting":         strconv.FormatBool(fanShouldBeOn),
	}

	func() {
		// Create a ticker to trigger events every minute
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if sensorStore.Inside.Size() < 10 || sensorStore.Outside.Size() < 10 {
					lg.Warnf("NOT sending to InfluxDB due to insufficient data (Inside/Outside): %d, %d",
						sensorStore.Inside.Size(), sensorStore.Outside.Size())
				} else {
					lg.Infof("Sending average values to InfluxDB (Inside/Outside): %d, %d",
						sensorStore.Inside.Size(), sensorStore.Outside.Size())
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
					point := write.NewPoint("dp", tags, fields, time.Now())
					if err := writeAPI.WritePoint(context.Background(), point); err != nil {
						lg.Error(err)
					}

					lg.Infof("Inside  (T/H): %5.1fC - %5.1f%%", sensorStore.Inside.AverageTemperature(),
						sensorStore.Inside.AverageHumidity())
					lg.Infof("Outside (T/H): %5.1fC - %5.1f%%", sensorStore.Outside.AverageTemperature(),
						sensorStore.Outside.AverageHumidity())
				}
			}
		}
	}()
}
