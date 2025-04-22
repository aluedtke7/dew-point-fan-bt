package main

import "github.com/spf13/viper"

// readConfig initializes application configuration values from the configuration file using the Viper library.
// It reads and validates the sensor, LCD, fan, and InfluxDB configurations, ensuring all values are correctly set.
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
	if lcdDelay < 1 || lcdDelay > 60 {
		lg.Fatal("Invalid LCD delay! Must be between 1 and 60 seconds.")
	}
	lcdScrollSpeed = viper.GetInt("lcd.scrollSpeed")
	if lcdScrollSpeed < 100 || lcdScrollSpeed > 10000 {
		lg.Fatal("Invalid LCD scroll speed! Must be between 100 and 10.000 ms.")
	}
	lcdScreenChange = viper.GetInt("lcd.screenChange")
	if lcdScreenChange < 3 || lcdScreenChange > 10 {
		lg.Fatal("Invalid LCD screen change interval! Must be between 3 and 10 seconds.")
	}

	fanConfig.MinDiff = viper.GetFloat64("fan.minDiff")
	if fanConfig.MinDiff < 1 || fanConfig.MinDiff > 10 {
		lg.Fatal("Invalid minimal difference! Must be between 1 and 10°C.")
	}
	fanConfig.Hysteresis = viper.GetFloat64("fan.hysteresis")
	if fanConfig.Hysteresis < 0.1 || fanConfig.Hysteresis > 5 {
		lg.Fatal("Invalid hysteresis! Must be between 0.1 and 5°C.")
	}
	fanConfig.MinHumidityInside = viper.GetFloat64("fan.minHumidityInside")
	if fanConfig.MinHumidityInside < 30 || fanConfig.MinHumidityInside > 70 {
		lg.Fatal("Invalid minimal inside humidity! Must be between 30 and 70%.")
	}
	fanConfig.MinTempInside = viper.GetFloat64("fan.minTempInside")
	if fanConfig.MinTempInside < 10 || fanConfig.MinTempInside > 40 {
		lg.Fatal("Invalid minimal inside temperature! Must be between 10 and 40°C.")
	}
	fanConfig.MinTempOutside = viper.GetFloat64("fan.minTempOutside")
	if fanConfig.MinTempOutside < -20 || fanConfig.MinTempOutside > 20 {
		lg.Fatal("Invalid minimal outside temperature! Must be between -20 and 20°C.")
	}

	influxConfig.Enabled = viper.GetBool("influx.enabled")
	influxConfig.Org = viper.GetString("influx.org")
	influxConfig.Bucket = viper.GetString("influx.bucket")
	influxConfig.Token = viper.GetString("influx.token")
	influxConfig.Url = viper.GetString("influx.url")
}
