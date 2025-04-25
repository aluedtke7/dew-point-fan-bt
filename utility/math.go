package utility

import "math"

// CalcDewPoint calculates the dew point temperature (°C) based on the given temperature (°C) and relative humidity (%).
// It uses the Magnus formula with coefficients adjusted for temperature above or below 0°C.
// Returns the calculated dew point temperature as a float64.
func CalcDewPoint(temperature, humidity float64) float64 {
	var a, b float64

	if temperature >= 0 {
		a = 7.5
		b = 237.3
	} else if temperature < 0 {
		a = 7.6
		b = 240.7
	}

	// saturation vapor pressure in hPa
	sdd := 6.1078 * math.Pow(10, (a*temperature)/(b+temperature))

	// vapor pressure in hPa
	dd := sdd * (humidity / 100)

	// v parameter
	v := math.Log10(dd / 6.1078)

	// dew point temperature (°C)
	return RoundDouble((b*v)/(a-v), 1)
}

// RoundDouble rounds a float64 value to the specified number of decimal places based on the given precision.
func RoundDouble(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}
