package monitor

import "testing"

func TestBar(t *testing.T) {
	if got := Bar(0, 8); got != 0 {
		t.Fatalf("0%% = %d", got)
	}
	if got := Bar(50, 8); got != 4 {
		t.Fatalf("50%% = %d", got)
	}
	if got := Bar(100, 8); got != 8 {
		t.Fatalf("100%% = %d", got)
	}
	if got := Bar(200, 8); got != 8 {
		t.Fatalf("clamped = %d", got)
	}
}

func TestColor(t *testing.T) {
	if got := Color(0); got != Off {
		t.Fatalf("off color = %d", got)
	}
	if got := Color(80); got != Amber {
		t.Fatalf("mid color = %d", got)
	}
	if got := Color(90); got != Red {
		t.Fatalf("high color = %d", got)
	}
}

func TestFineBar(t *testing.T) {
	cases := []struct {
		percent   float64
		lit       int
		intensity byte
	}{
		{0, 0, 0}, {1, 1, 1}, {8.33, 1, 1}, {8.34, 1, 2}, {16.66, 1, 2},
		{16.67, 1, 3}, {25, 1, 3}, {25.01, 2, 1}, {100, 4, 3},
	}
	for _, tc := range cases {
		lit, intensity := FineBar(tc.percent, 4)
		if lit != tc.lit || intensity != tc.intensity {
			t.Fatalf("%.2f%%: got (%d,%d), want (%d,%d)", tc.percent, lit, intensity, tc.lit, tc.intensity)
		}
	}
}
