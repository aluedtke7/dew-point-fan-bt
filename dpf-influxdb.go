package main

import "time"

func sendToInfluxDb() {
	func() {
		// Create a ticker to trigger events every minute
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				lg.Infof("Sending average values to InfluxDB: %d, %d",
					sensorStore.Inside.Size(), sensorStore.Outside.Size())
				lg.Infof("Inside  (T/H): %5.1fC - %5.1f%%", sensorStore.Inside.AverageTemperature(),
					sensorStore.Inside.AverageHumidity())
				lg.Infof("Outside (T/H): %5.1fC - %5.1f%%", sensorStore.Outside.AverageTemperature(),
					sensorStore.Outside.AverageHumidity())
			}
		}
	}()
}
