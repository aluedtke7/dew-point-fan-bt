package display

import (
	"dpf-bt/sensor"
	"fmt"
	"math"
	"strings"
	"time"
)

// printLine formats and prints text to a specific line on the display, with optional scrolling.
func printLine(disp Display, line int, text string, scroll bool) {
	if scroll {
		disp.PrintLine(line, strings.TrimSpace(text), scroll)
	} else {
		disp.PrintLine(line, strings.TrimRight(text, " "), scroll)
	}
}

// formatUpDays converts a given time in seconds to a string representing the number of full days.
func formatUpDays(seconds uint32) string {
	days := seconds / (24 * 3600)
	return fmt.Sprintf("%dd", days)
}

// StartScreen initializes the display with a startup message and the provided IP address.
func StartScreen(display Display, buildTime string, ip string) {
	printLine(display, 0, "DewPointFan BT v1", false)
	printLine(display, 1, buildTime, false)
	printLine(display, 2, "", false)
	printLine(display, 3, "IP: "+ip, false)
}

// MainScreen displays sensor data for inside and outside environments, including temperature,
// humidity, and dew point.
func MainScreen(display Display, sensorInside sensor.SensorData, sensorOutside sensor.SensorData) {
	printLine(display, 0, "DPF   Inside Outside", false)
	printLine(display, 1, fmt.Sprintf("Temp: %5.1fC  %5.1fC", sensorInside.Temperature,
		sensorOutside.Temperature), false)
	printLine(display, 2, fmt.Sprintf("Hum:  %5.1f%%  %5.1f%%", sensorInside.Humidity,
		sensorOutside.Humidity), false)
	printLine(display, 3, fmt.Sprintf("DP:   %5.1fC  %5.1fC", sensorInside.DewPoint,
		sensorOutside.DewPoint), false)
}

// InfoScreen displays information about sensor data on a display, including RSSI, battery levels,
// and uptime details.
func InfoScreen(display Display, sensorInside sensor.SensorData, sensorOutside sensor.SensorData) {
	printLine(display, 0, "DPF   Inside Outside", false)
	printLine(display, 1, fmt.Sprintf("RSSI:%7d %7d", sensorInside.RSSI,
		sensorOutside.RSSI), false)
	printLine(display, 2, fmt.Sprintf("Bat: %7d %7d", sensorInside.BatLevel,
		sensorOutside.BatLevel), false)
	printLine(display, 3, fmt.Sprintf("Up:  %7s %7s", formatUpDays(sensorInside.Uptime),
		formatUpDays(sensorOutside.Uptime)), false)
}

// ResultScreen displays fan status, its operation reason, dew point differences, and sensors last-seen durations.
func ResultScreen(display Display, result sensor.ResultData, sensorInside sensor.SensorData,
	sensorOutside sensor.SensorData, fanConfig sensor.FanConfig) {
	isOn := "OFF"
	shouldBeOn := "OFF"
	if result.IsOn {
		isOn = "ON"
	}
	if result.ShouldBeOn {
		shouldBeOn = "ON"
	}
	now := time.Now()
	insideLastSeen := int32(math.Min(float64(now.Sub(sensorInside.Scanned).Seconds()), 9999))
	outsideLastSeen := int32(math.Min(float64(now.Sub(sensorOutside.Scanned).Seconds()), 9999))
	printLine(display, 0, fmt.Sprintf("Fan is %s (%s)", isOn, shouldBeOn), false)
	printLine(display, 1, fmt.Sprintf(" %18s ", sensor.ReasonName[result.Reason]), false)
	printLine(display, 2, fmt.Sprintf("Dp diff:%5.1fC (%3.1f)",
		sensorInside.DewPoint-sensorOutside.DewPoint, fanConfig.MinDiff+fanConfig.Hysteresis), false)
	printLine(display, 3, fmt.Sprintf("In/Out:  %4ds %4ds", insideLastSeen, outsideLastSeen), false)
}
