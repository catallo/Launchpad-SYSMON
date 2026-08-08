package monitor

import "testing"

func TestLogicalThreads(t *testing.T) {
	input := []float64{10, 20, 30, 40, 50, 60, 70, 80}
	got, err := LogicalThreads(input)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 8 {
		t.Fatalf("got %d threads", len(got))
	}
	for i := range input {
		if got[i] != input[i] {
			t.Fatalf("thread %d: got %.1f want %.1f", i, got[i], input[i])
		}
	}
}
