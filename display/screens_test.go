package display

import (
	"testing"
)

func TestFormatUpDays(t *testing.T) {
	tests := []struct {
		name     string
		seconds  uint32
		expected string
	}{
		{
			name:     "exact_one_day",
			seconds:  24 * 60 * 60,
			expected: "1d",
		},
		{
			name:     "less_than_one_day",
			seconds:  60 * 60,
			expected: "0d",
		},
		{
			name:     "multiple_days",
			seconds:  3 * 24 * 60 * 60,
			expected: "3d",
		},
		{
			name:     "large_number_of_days",
			seconds:  1000 * 24 * 60 * 60,
			expected: "1000d",
		},
		{
			name:     "zero_seconds",
			seconds:  0,
			expected: "0d",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatUpDays(tt.seconds)
			if got != tt.expected {
				t.Errorf("formatUpDays(%v) = %v, want %v", tt.seconds, got, tt.expected)
			}
		})
	}
}
