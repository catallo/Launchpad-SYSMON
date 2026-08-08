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
