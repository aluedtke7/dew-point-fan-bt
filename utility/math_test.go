package utility

import (
	"math"
	"testing"
)

func TestCalcDewPoint(t *testing.T) {
	tests := []struct {
		name        string
		temperature float64
		humidity    float64
		expected    float64
	}{
		{
			name:        "Positive temperature and high humidity",
			temperature: 25.0,
			humidity:    80.0,
			expected:    21.3,
		},
		{
			name:        "Negative temperature and high humidity",
			temperature: -5.0,
			humidity:    90.0,
			expected:    -6.4,
		},
		{
			name:        "Zero temperature and high humidity",
			temperature: 0.0,
			humidity:    85.0,
			expected:    -2.2,
		},
		{
			name:        "Positive temperature and low humidity",
			temperature: 30.0,
			humidity:    30.0,
			expected:    10.5,
		},
		{
			name:        "Negative temperature and low humidity",
			temperature: -10.0,
			humidity:    20.0,
			expected:    -28.7,
		},
		{
			name:        "Extreme positive temperature",
			temperature: 50.0,
			humidity:    50.0,
			expected:    36.7,
		},
		{
			name:        "Extreme negative temperature",
			temperature: -40.0,
			humidity:    50.0,
			expected:    -46.4,
		},
		{
			name:        "Zero temperature and zero humidity",
			temperature: 0.0,
			humidity:    0.0,
			expected:    -math.Inf(1),
		},
		{
			name:        "Negative temperature and zero humidity",
			temperature: -10.0,
			humidity:    0.0,
			expected:    -math.Inf(1),
		},
		{
			name:        "Positive temperature and 100 humidity",
			temperature: 20.0,
			humidity:    100.0,
			expected:    20.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			const precision = 6 // Precision for comparison
			result := CalcDewPoint(tt.temperature, tt.humidity)
			if math.Abs(result-tt.expected) > math.Pow(10, -precision) {
				t.Errorf("CalcDewPoint(%v, %v) = %v, expected %v",
					tt.temperature, tt.humidity, result, tt.expected)
			}
		})
	}
}

func TestRoundDouble(t *testing.T) {
	tests := []struct {
		name      string
		val       float64
		precision uint
		expected  float64
	}{
		{
			name:      "Round to zero precision",
			val:       123.456,
			precision: 0,
			expected:  123,
		},
		{
			name:      "Round to one decimal place",
			val:       123.456,
			precision: 1,
			expected:  123.5,
		},
		{
			name:      "Round to two decimal places",
			val:       123.456,
			precision: 2,
			expected:  123.46,
		},
		{
			name:      "No rounding needed",
			val:       100.0,
			precision: 2,
			expected:  100.0,
		},
		{
			name:      "Round negative number",
			val:       -123.456,
			precision: 2,
			expected:  -123.46,
		},
		{
			name:      "Round to zero decimal places (negative number)",
			val:       -123.756,
			precision: 0,
			expected:  -124,
		},
		{
			name:      "Small number rounding",
			val:       0.0049,
			precision: 2,
			expected:  0.0,
		},
		{
			name:      "Small number rounding up",
			val:       0.0051,
			precision: 2,
			expected:  0.01,
		},
		{
			name:      "Large number rounding",
			val:       98765432.123456,
			precision: 3,
			expected:  98765432.123,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RoundDouble(tt.val, tt.precision)
			if result != tt.expected {
				t.Errorf("RoundDouble(%v, %v) = %v, expected %v",
					tt.val, tt.precision, result, tt.expected)
			}
		})
	}
}
