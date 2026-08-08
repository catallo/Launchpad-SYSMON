// Package monitor maps normalized system metrics to Launchpad indicators.
package monitor

import "math"

// Legacy Launchpad Mini velocities: red and green intensity bit fields.
const (
	Off   byte = 0x00
	Green byte = 0x3C
	Amber byte = 0x3F
	Red   byte = 0x0F
)

// Bar returns the number of active pads for percentage on a bar of length pads.
func Bar(percent float64, pads int) int {
	if pads <= 0 || percent <= 0 {
		return 0
	}
	if percent >= 100 {
		return pads
	}
	return int(math.Ceil(percent / 100 * float64(pads)))
}

// Color represents normal, warning, and critical utilization.
func Color(percent float64) byte {
	switch {
	case percent <= 0:
		return Off
	case percent < 75:
		return Green
	case percent < 90:
		return Amber
	default:
		return Red
	}
}

// FineBar maps a percentage to a bar with three brightness levels per pad.
// Four pads therefore represent twelve distinct utilization steps.
func FineBar(percent float64, pads int) (lit int, intensity byte) {
	if pads <= 0 || percent <= 0 {
		return 0, 0
	}
	steps := pads * 3
	if percent >= 100 {
		return pads, 3
	}
	step := int(math.Ceil(percent / 100 * float64(steps)))
	return int(math.Ceil(float64(step) / 3)), byte((step-1)%3 + 1)
}
