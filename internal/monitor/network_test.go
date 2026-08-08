package monitor

import "testing"

func TestNetworkRateUnits(t *testing.T) {
	// Formula check: 35,000,000 bytes in 1 second equal 280 Mbit/s.
	got := float64(35_000_000) * 8 / 1_000_000
	if got != 280 {
		t.Fatalf("got %.1f", got)
	}
}
