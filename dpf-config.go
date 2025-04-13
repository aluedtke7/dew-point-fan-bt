package main

import "github.com/spf13/viper"

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
	if len(sensors.InsideData.MacAddress) != 17 || len(sensors.OutsideData.MacAddress) != 17 {
		lg.Fatal("Invalid MAC address! Must be 17 characters long.")
	}
	lg.Infof("Inside sensor:  MAC %s - Temp cal = %.2f - Humidity cal = %.2f",
		sensors.InsideData.MacAddress, sensors.InsideCalibration.Temperature, sensors.InsideCalibration.Humidity)
	lg.Infof("Outside sensor: MAC %s - Temp cal = %.2f - Humidity cal = %.2f",
		sensors.OutsideData.MacAddress, sensors.OutsideCalibration.Temperature, sensors.OutsideCalibration.Humidity)

	lcdDelay = viper.GetInt("lcd.delay")
	lcdScrollSpeed = viper.GetInt("lcd.scrollSpeed")
	lcdScreenChange = viper.GetInt("lcd.screenChange")
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

	influxConfig.Org = viper.GetString("influx.org")
	influxConfig.Bucket = viper.GetString("influx.bucket")
	influxConfig.Token = viper.GetString("influx.token")
	influxConfig.Url = viper.GetString("influx.url")
}
