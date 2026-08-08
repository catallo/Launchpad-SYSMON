package monitor

import "testing"

func TestPhysicalCores(t *testing.T) {
	got, err := PhysicalCores([]float64{10, 20, 30, 40, 50, 60, 70, 80})
	if err != nil {
		t.Fatal(err)
	}
	want := []float64{30, 40, 50, 60}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("core %d: got %.1f want %.1f", i, got[i], want[i])
		}
	}
}
