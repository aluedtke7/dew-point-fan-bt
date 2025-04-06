package main

import (
	"dpf-bt/bluetooth"
	"dpf-bt/display"
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

var (
	fanConfig       sensor.FanConfig = sensor.FanConfig{}
	fanData         sensor.FanData   = sensor.FanData{}
	sensors         sensor.Sensors   = sensor.Sensors{}
	disp            display.Display
	lcdDelay        int
	lcdScrollSpeed  int
	lcdScreenChange int
	ipAddress       string
	lg              = logger.NewPackageLogger("main", logger.InfoLevel)
)

func main() {
	defer func() {
		_ = logger.FinalizeLogger()
	}()
	viper.SetConfigName("config")
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

	// disp, err = display.New(false, lcdScrollSpeed, lcdDelay)
	disp, err = display.NewTerminalLcd(false, lcdScrollSpeed, lcdDelay)
	if err != nil {
		lg.Errorf("Couldn't initialize display: %s", err)
	} else {
		disp.Backlight(true)
		ipAddress = utility.LogNetworkInterfacesAndGetIpAdr()
		display.StartScreen(disp, ipAddress)
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

	go func() {
		// Create a ticker to trigger events every 'lcdScreenChange' seconds
		ticker := time.NewTicker(time.Duration(lcdScreenChange) * time.Second)
		defer ticker.Stop()
		step := 0
		// Loop to handle toggling and communication through channels
		for {
			select {
			case <-ticker.C:
				switch step {
				case 0, 2, 4, 6:
					display.MainScreen(disp, sensors.InsideData, sensors.OutsideData)
				case 1, 5:
					display.FanInfoScreen(disp, fanData, sensors.InsideData, sensors.OutsideData)
				case 3, 7:
					display.InfoScreen(disp, sensors.InsideData, sensors.OutsideData)
				case 8:
					display.StartScreen(disp, ipAddress)
				}
				step += 1
				if step > 8 {
					step = 0
				}
			}
		}
	}()

	err = adapter.Scan(onScan)
	if err != nil {
		lg.Panic("failed to register scan callback")
	}
}

func readConfig() {
	err := viper.ReadInConfig()
	if err != nil {
		lg.Fatalf("Fatal error reading config file: %s \n", err)
	}
	sensors.InsideData.MacAddress = viper.GetString("inside.mac")
	sensors.InsideCalibration.Temperature = viper.GetFloat64("inside.temperature-calibration")
	sensors.InsideCalibration.Humidity = viper.GetFloat64("inside.humidity-calibration")
	sensors.OutsideData.MacAddress = viper.GetString("outside.mac")
	sensors.OutsideCalibration.Temperature = viper.GetFloat64("outside.temperature-calibration")
	sensors.OutsideCalibration.Humidity = viper.GetFloat64("outside.humidity-calibration")
	lcdDelay = viper.GetInt("lcd.delay")
	lcdScrollSpeed = viper.GetInt("lcd.scrollSpeed")
	lcdScreenChange = viper.GetInt("lcd.screenChange")
	lg.Infof("Inside sensor:  MAC %s - Temp cal = %.2f - Humidity cal = %.2f",
		sensors.InsideData.MacAddress, sensors.InsideCalibration.Temperature, sensors.InsideCalibration.Humidity)
	lg.Infof("Outside sensor: MAC %s - Temp cal = %.2f - Humidity cal = %.2f",
		sensors.OutsideData.MacAddress, sensors.OutsideCalibration.Temperature, sensors.OutsideCalibration.Humidity)
	if len(sensors.InsideData.MacAddress) != 17 || len(sensors.OutsideData.MacAddress) != 17 {
		lg.Fatal("Invalid MAC address! Must be 17 characters long.")
	}
	// if lcdDelay < 1 || lcdDelay > 10 {
	// 	lg.Fatal("Invalid LCD delay! Must be between 1 and 10 seconds.")
	// }
	if lcdScrollSpeed < 100 || lcdScrollSpeed > 10000 {
		lg.Fatal("Invalid LCD scroll speed! Must be between 100 and 10.000 ms.")
	}
	if lcdScreenChange < 3 || lcdScreenChange > 10 {
		lg.Fatal("Invalid LCD screen change interval! Must be between 3 and 10 seconds.")
	}
	fanConfig.MinDiff = viper.GetFloat64("fan.minDiff")
	fanConfig.Hysteresis = viper.GetFloat64("fan.hysteresis")
	fanConfig.MinHumidityInside = viper.GetFloat64("fan.minHumidityInside")
	fanConfig.MinTempInside = viper.GetFloat64("fan.minTempInside")
	fanConfig.MinTempOutside = viper.GetFloat64("fan.minTempOutside")
}

func onScan(_ *bt.Adapter, scanResult bt.ScanResult) {
	if scanResult.LocalName() == "ThermoBeacon" {
		bluetooth.ProcessAdvertisement(scanResult, &sensors)
	}
}
