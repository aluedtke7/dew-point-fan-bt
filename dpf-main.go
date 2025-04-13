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
	"syscall"
	"time"
	bt "tinygo.org/x/bluetooth"
)

const maxSensorData = 20

var (
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
)

// main is the entry point of the application. It initializes configurations, hardware, and services, and manages app lifecycle.
func main() {
	defer func() {
		_ = logger.FinalizeLogger()
	}()
	viper.SetConfigName("config-private")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	viper.OnConfigChange(func(e fsnotify.Event) {
		lg.Info("Config file changed:", e.Name)
		readConfig()
	})
	viper.WatchConfig()
	readConfig()

	adapter := bt.DefaultAdapter
	err := adapter.Enable()
	if err != nil {
		lg.Panic("failed to enable BLE adapter")
	}

	disp, err = display.New(false, lcdScrollSpeed, lcdDelay)
	if err != nil {
		lg.Errorf("Couldn't initialize display: %s", err)
	} else {
		disp.Backlight(true)
		ipAddress = utility.LogNetworkInterfacesAndGetIpAdr()
		display.StartScreen(disp, ipAddress)
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
	// go sendToInfluxDb()

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
	if inside.Scanned.IsZero() || outside.Scanned.IsZero() {
		resultData.ShouldBeOn = false
		resultData.Outcome = sensor.ReasonNoData
		return
	}
	last5Minute := time.Now().Add(-5 * time.Minute)
	if inside.Scanned.Before(last5Minute) || outside.Scanned.Before(last5Minute) {
		resultData.ShouldBeOn = false
		resultData.Outcome = sensor.ReasonNoEnoughData
		return
	}
	if inside.Temperature < fanConfig.MinTempInside {
		resultData.ShouldBeOn = false
		resultData.Outcome = sensor.ReasonInsideTempTooLow
		return
	}
	if outside.Temperature < fanConfig.MinTempOutside {
		resultData.ShouldBeOn = false
		resultData.Outcome = sensor.ReasonOutsideTempTooLow
		return
	}
	if inside.Humidity < fanConfig.MinHumidityInside {
		resultData.ShouldBeOn = false
		resultData.Outcome = sensor.ReasonInsideHumidityTooLow
		return
	}
	deltaDp := inside.DewPoint - outside.DewPoint
	if deltaDp < fanConfig.MinDiff {
		resultData.ShouldBeOn = false
		resultData.Outcome = sensor.ReasonDewPointUnderHyst
		return
	}
	if deltaDp > fanConfig.MinDiff+fanConfig.Hysteresis {
		resultData.ShouldBeOn = true
		resultData.Outcome = sensor.ReasonDewPointOverHyst
		return
	}
	resultData.ShouldBeOn = false
	resultData.Outcome = sensor.ReasonDewPointUnderHyst
}
