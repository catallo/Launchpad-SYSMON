package monitor

import "testing"

func TestMemoryCellsKeepsClassesVisible(t *testing.T) {
	s := MemorySnapshot{Total: 100, User: 48, System: 4, Cache: 36, Free: 12}
	got := MemoryCells(s, 8)
	want := []MemoryClass{MemoryUser, MemoryUser, MemoryUser, MemorySystem, MemoryCache, MemoryCache, MemoryCache, MemoryFree}
	if len(got) != len(want) {
		t.Fatalf("cells=%d", len(got))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("cell %d: got %d want %d", i, got[i], want[i])
		}
	}
}

func TestMemoryCellsEmpty(t *testing.T) {
	if got := MemoryCells(MemorySnapshot{}, 8); got != nil {
		t.Fatalf("got %#v", got)
	}
}
