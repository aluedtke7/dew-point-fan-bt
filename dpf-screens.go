package main

import (
	"dpf-bt/display"
	"time"
)

func showScreens() {
	func() {
		// Create a ticker to trigger events every 'lcdScreenChange' seconds
		ticker := time.NewTicker(time.Duration(lcdScreenChange) * time.Second)
		defer ticker.Stop()
		step := 0
		// Loop to handle toggling and communication through channels
		for {
			computeResults(sensors.InsideData, sensors.OutsideData, &resultData)
			ioPins.SetFan(resultData.ShouldBeOn)
			resultData.IsOn = ioPins.ReadFanSense()
			select {
			case <-ticker.C:
				switch step {
				case 0, 3, 6:
					display.MainScreen(disp, sensors.InsideData, sensors.OutsideData)
				case 1, 4, 7:
					display.ResultScreen(disp, resultData, sensors.InsideData, sensors.OutsideData, fanConfig)
				case 2, 5:
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
}
