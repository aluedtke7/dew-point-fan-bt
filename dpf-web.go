package main

import (
	"dpf-bt/sensor"
	"encoding/json"
	"fmt"
	"github.com/d2r2/go-logger"
	"net/http"
	"strings"
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

const (
	webServerPort = ":8080"
	webServerHost = "0.0.0.0"
)

type webServer struct {
	sensorStore    *sensor.SensorStore
	resultData     *sensor.ResultData
	fanConfig      *sensor.FanConfig
	remoteOverride *int
}

var lgWeb = logger.NewPackageLogger("web", logger.InfoLevel)

// startWebserver initializes and starts a web server to display sensor data and control fan settings interactively.
func startWebserver() {
	srv := &webServer{
		sensorStore:    &sensorStore,
		resultData:     &resultData,
		fanConfig:      &fanConfig,
		remoteOverride: &remoteOverride,
	}

	go func() {
		http.HandleFunc("/", srv.handleMainPage)
		http.HandleFunc("/info", srv.handleInfo)
		http.HandleFunc("/override", srv.handleOverride)

		lgWeb.Fatal(http.ListenAndServe(webServerHost+webServerPort, nil))
	}()
}

func (s *webServer) handleMainPage(w http.ResponseWriter, _ *http.Request) {
	var b strings.Builder

	shouldBeOn := s.getFanStateText(s.resultData.ShouldBeOn)
	isOn := s.getFanStateText(s.resultData.IsOn)
	dewPointDiff := s.sensorStore.Inside.AverageDewPoint() - s.sensorStore.Outside.AverageDewPoint()

	b.WriteString("Dew Point Fan\n-----------------------------------------------------\n")
	_, _ = fmt.Fprintf(&b, "Inside   DP: %6.1f, Temp: %5.1f°C, Humidity: %5.1f%%\n",
		s.sensorStore.Inside.AverageDewPoint(),
		s.sensorStore.Inside.AverageTemperature(),
		s.sensorStore.Inside.AverageHumidity())
	_, _ = fmt.Fprintf(&b, "Outside  DP: %6.1f, Temp: %5.1f°C, Humidity: %5.1f%%\n",
		s.sensorStore.Outside.AverageDewPoint(),
		s.sensorStore.Outside.AverageTemperature(),
		s.sensorStore.Outside.AverageHumidity())
	_, _ = fmt.Fprintf(&b, "Diff     DP: %6.1f\n", dewPointDiff)
	_, _ = fmt.Fprintf(&b, "Fan should be %s                         Fan is %s", shouldBeOn, isOn)

	_, _ = fmt.Fprint(w, b.String())
}

func (s *webServer) handleInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	inf := &info{
		Update:         time.Now().Format(time.DateTime),
		Sensors:        s.getSensorData(),
		Reason:         int(s.resultData.Reason),
		Venting:        s.resultData.ShouldBeOn,
		Override:       s.resultData.ShouldBeOn != s.resultData.IsOn,
		RemoteOverride: *s.remoteOverride,
		DiffMin:        s.fanConfig.MinDiff,
		Hysteresis:     s.fanConfig.Hysteresis,
	}

	if err := s.writeJSON(w, inf); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *webServer) handleOverride(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	lgWeb.Info("POST API called")

	var remote remoteControl
	if err := json.NewDecoder(r.Body).Decode(&remote); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	lgWeb.Infof("POST API called with override: %d", remote.Override)
	*s.remoteOverride = remote.Override

	if err := s.writeJSON(w, remote); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *webServer) getSensorData() []sensorData {
	return []sensorData{
		{
			Name:        "Inside",
			Temperature: s.sensorStore.Inside.AverageTemperature(),
			Humidity:    s.sensorStore.Inside.AverageHumidity(),
			DewPoint:    s.sensorStore.Inside.AverageDewPoint(),
		},
		{
			Name:        "Outside",
			Temperature: s.sensorStore.Outside.AverageTemperature(),
			Humidity:    s.sensorStore.Outside.AverageHumidity(),
			DewPoint:    s.sensorStore.Outside.AverageDewPoint(),
		},
	}
}

func (s *webServer) getFanStateText(state bool) string {
	if state {
		return "ON"
	}
	return "OFF"
}

func (s *webServer) writeJSON(w http.ResponseWriter, v interface{}) error {
	j, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	_, err = w.Write(j)
	return err
}
