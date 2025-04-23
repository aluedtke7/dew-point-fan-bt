package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// sensorData represents the data collected from a sensor, including its name, temperature, humidity, and dew point.
type sensorData struct {
	Name        string  `json:"name"`
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
	DewPoint    float64 `json:"dew_point"`
}

// info represents the main structure for current system data, including sensor readings and fan control states.
type info struct {
	Update         string       `json:"update"`
	Sensors        []sensorData `json:"sensors"`
	Reason         int          `json:"reason"`
	Venting        bool         `json:"venting"`
	Override       bool         `json:"override"`
	RemoteOverride int          `json:"remote_override"`
	DiffMin        float64      `json:"diff_min"`
	Hysteresis     float64      `json:"hysteresis"`
}

// remoteControl represents the structure for managing remote override control for the fan system.
// Override defines the override state, where the fan is forced to a specific state through external input.
type remoteControl struct {
	Override int `json:"override"`
}

var remoteOverride = 0

// startWebserver initializes and starts a web server to display sensor data and control fan settings interactively.
func startWebserver() {
	go func() {
		shouldBeOn := "OFF"
		if resultData.ShouldBeOn {
			shouldBeOn = "ON"
		}
		isOn := "OFF"
		if resultData.IsOn {
			isOn = "ON"
		}
		// plain text browser page
		webHandler := func(w http.ResponseWriter, req *http.Request) {
			_, _ = fmt.Fprintf(w, "Dew Point Fan\n"+
				"-----------------------------------------------------\n"+
				"Inside   DP: %6.1f, Temp: %5.1f°C, Humidity: %5.1f%%\n"+
				"Outside  DP: %6.1f, Temp: %5.1f°C, Humidity: %5.1f%%\n"+
				"Diff     DP: %6.1f\n"+
				"Fan should be %s                         Fan is %s",
				sensorStore.Inside.AverageDewPoint(),
				sensorStore.Inside.AverageTemperature(),
				sensorStore.Inside.AverageHumidity(),
				sensorStore.Outside.AverageDewPoint(),
				sensorStore.Outside.AverageTemperature(),
				sensorStore.Outside.AverageHumidity(),
				sensorStore.Inside.AverageDewPoint()-sensorStore.Outside.AverageDewPoint(),
				shouldBeOn, isOn,
			)
		}
		http.HandleFunc("/", webHandler)

		// data in JSON format
		infoHandler := func(w http.ResponseWriter, req *http.Request) {
			if req.Method == "GET" {
				inf := new(info)
				inf.Update = time.Now().Format(time.DateTime)
				inf.Sensors = []sensorData{
					{"Inside",
						sensorStore.Inside.AverageTemperature(),
						sensorStore.Inside.AverageHumidity(),
						sensorStore.Inside.AverageDewPoint(),
					},
					{"Outside",
						sensorStore.Outside.AverageTemperature(),
						sensorStore.Outside.AverageHumidity(),
						sensorStore.Outside.AverageDewPoint(),
					},
				}
				inf.Reason = int(resultData.Outcome)
				inf.Venting = resultData.ShouldBeOn
				inf.Override = resultData.ShouldBeOn != resultData.IsOn
				inf.RemoteOverride = remoteOverride
				inf.DiffMin = fanConfig.MinDiff
				inf.Hysteresis = fanConfig.Hysteresis
				j, _ := json.MarshalIndent(inf, "", "  ")
				_, _ = w.Write(j)
			}
		}
		http.HandleFunc("/info", infoHandler)

		// POST handler for changing/overriding fan state
		overrideHandler := func(w http.ResponseWriter, req *http.Request) {
			if req.Method == "POST" {
				lg.Info("POST API called")
				decoder := json.NewDecoder(req.Body)
				remote := &remoteControl{}
				err := decoder.Decode(remote)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				remoteOverride = remote.Override
				j, _ := json.MarshalIndent(remote, "", "  ")
				_, _ = w.Write(j)
			}
		}
		http.HandleFunc("/override", overrideHandler)
		log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
	}()
}
