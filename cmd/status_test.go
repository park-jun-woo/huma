package cmd

import (
	"math"
	"testing"
)

func TestPct(t *testing.T) {
	tests := []struct {
		name  string
		n     int
		total int
		want  float64
	}{
		{"zero total", 5, 0, 0},
		{"zero n", 0, 10, 0},
		{"half", 5, 10, 50},
		{"full", 10, 10, 100},
		{"one third", 1, 3, 33.333333},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := pct(tt.n, tt.total)
			if math.Abs(got-tt.want) > 0.001 {
				t.Errorf("pct(%d, %d) = %f, want %f", tt.n, tt.total, got, tt.want)
			}
		})
	}
}
