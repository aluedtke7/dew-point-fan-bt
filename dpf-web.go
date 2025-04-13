package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type sensorData struct {
	Name        string  `json:"name"`
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
	DewPoint    float64 `json:"dew_point"`
}

type info struct {
	Update         string       `json:"update"`
	Sensors        []sensorData `json:"sensors"`
	Venting        bool         `json:"venting"`
	Override       bool         `json:"override"`
	RemoteOverride int          `json:"remote_override"`
	DiffMin        float64      `json:"diff_min"`
	Hysteresis     float64      `json:"hysteresis"`
}

type remoteControl struct {
	Override int `json:"override"`
}

var remoteOverride = 0

func startWebserver() {
	// a little http server to show current values
	go func() {
		shouldBeOn := "OFF"
		if resultData.ShouldBeOn {
			shouldBeOn = "ON"
		}
		isOn := "OFF"
		if resultData.IsOn {
			isOn = "ON"
		}
		// browser page plain text
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

		// POST handler for changing fanIsOn
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
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()
}
