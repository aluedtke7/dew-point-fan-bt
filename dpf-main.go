package main

import (
	"dpf-bt/bluetooth"
	"dpf-bt/display"
	"dpf-bt/gpio"
	"dpf-bt/sensor"
	"dpf-bt/utility"
	"github.com/d2r2/go-logger"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
	bt "tinygo.org/x/bluetooth"
)

const maxSensorData = 20

var (
	buildTime    = "---"
	lg           = logger.NewPackageLogger("main", logger.InfoLevel)
	fanConfig    = sensor.FanConfig{}
	influxConfig = sensor.InfluxDbConfig{}
	resultData   = sensor.ResultData{}
	sensors      = sensor.Sensors{}
	sensorStore  = sensor.SensorStore{
		Inside:  *sensor.NewSensorDataStore(maxSensorData),
		Outside: *sensor.NewSensorDataStore(maxSensorData),
	}
	disp            display.Display
	ioPins          gpio.Gpio
	lcdDelay        int
	lcdScrollSpeed  int
	lcdScreenChange int
	ipAddress       string
	remoteOverride  = 0
)

// The main function is the entry point of the application. It initializes configurations, hardware, and
// services and manages the app lifecycle.
func main() {
	defer func() {
		_ = logger.FinalizeLogger()
	}()
	pathOfBinary, err := os.Executable()
	if err != nil {
		lg.Errorf("Couldn't get path of executable: %s", err)
		return
	}
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(filepath.Dir(pathOfBinary))
	viper.OnConfigChange(func(e fsnotify.Event) {
		lg.Info("Config file changed:", e.Name)
		readConfig()
	})
	viper.WatchConfig()
	readConfig()
	lg.Infof("Build timestamp: %s", buildTime)

	adapter := bt.DefaultAdapter
	err = adapter.Enable()
	if err != nil {
		lg.Panic("failed to enable BLE adapter")
	}

	disp, err = display.New(false, lcdScrollSpeed, lcdDelay)
	if err != nil {
		lg.Errorf("Couldn't initialize display: %s", err)
	} else {
		disp.Backlight(true)
		ipAddress = utility.LogNetworkInterfacesAndGetIpAdr()
		display.StartScreen(disp, buildTime, ipAddress)
	}
	ioPins, err = gpio.New()
	if err != nil {
		lg.Errorf("Couldn't initialize GPIO: %s", err)
	}

	var ctrlChan = make(chan os.Signal, 1)
	signal.Notify(ctrlChan, os.Interrupt, syscall.SIGTERM)
	// this goroutine is waiting for being stopped
	go func() {
		<-ctrlChan
		disp.Backlight(false)
		lg.Info("Ctrl+C received... Exiting")
		os.Exit(1)
	}()

	go showScreens()
	go startWebserver()
	if influxConfig.Enabled {
		go sendToInfluxDb()
	}

	err = adapter.Scan(onScan)
	if err != nil {
		lg.Panic("failed to register scan callback")
	}
}

func onScan(_ *bt.Adapter, scanResult bt.ScanResult) {
	if scanResult.LocalName() == "ThermoBeacon" {
		bluetooth.ProcessAdvertisement(scanResult, &sensors, &sensorStore)
	}
}

func computeResults(inside sensor.SensorData, outside sensor.SensorData, resultData *sensor.ResultData) {
	if remoteOverride > 0 {
		// manual override via REST api
		if remoteOverride == 1 {
			resultData.ShouldBeOn = true
			resultData.Reason = sensor.ReasonSoftOverrideOn
		} else {
			resultData.ShouldBeOn = false
			resultData.Reason = sensor.ReasonSoftOverrideOff
		}
		return
	}
	if inside.Scanned.IsZero() || outside.Scanned.IsZero() {
		resultData.ShouldBeOn = false
		resultData.Reason = sensor.ReasonNoData
		return
	}
	last5Minute := time.Now().Add(-5 * time.Minute)
	if inside.Scanned.Before(last5Minute) || outside.Scanned.Before(last5Minute) {
		resultData.ShouldBeOn = false
		resultData.Reason = sensor.ReasonNoEnoughData
		return
	}
	if inside.Temperature < fanConfig.MinTempInside {
		resultData.ShouldBeOn = false
		resultData.Reason = sensor.ReasonInsideTempTooLow
		return
	}
	if outside.Temperature < fanConfig.MinTempOutside {
		resultData.ShouldBeOn = false
		resultData.Reason = sensor.ReasonOutsideTempTooLow
		return
	}
	if inside.Humidity < fanConfig.MinHumidityInside {
		resultData.ShouldBeOn = false
		resultData.Reason = sensor.ReasonInsideHumidityTooLow
		return
	}
	deltaDp := inside.DewPoint - outside.DewPoint
	if deltaDp < fanConfig.MinDiff {
		resultData.ShouldBeOn = false
		resultData.Reason = sensor.ReasonDewPointUnderHyst
		return
	}
	if deltaDp >= fanConfig.MinDiff+fanConfig.Hysteresis {
		resultData.ShouldBeOn = true
		resultData.Reason = sensor.ReasonDewPointOverHyst
		return
	}
	if deltaDp >= fanConfig.MinDiff && deltaDp < fanConfig.MinDiff+fanConfig.Hysteresis {
		// we don't change the fan state since we don't know if the dew point is rising or falling
		resultData.Reason = sensor.ReasonDewPointInBetween
		return
	}
	resultData.ShouldBeOn = false
	resultData.Reason = sensor.ReasonUnknown
}
