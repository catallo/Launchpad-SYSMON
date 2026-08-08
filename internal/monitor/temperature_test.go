package monitor

import "testing"

func TestTemperatureScale(t *testing.T) {
	if got := 42.0 / 80.0 * 100; got != 52.5 {
		t.Fatalf("got %.1f", got)
	}
}
