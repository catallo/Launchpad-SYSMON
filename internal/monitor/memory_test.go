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

func TestFineMemoryCells(t *testing.T) {
	cells := FineMemoryCells(MemorySnapshot{Total: 100, User: 48, System: 4, Cache: 36, Free: 12}, 8)
	if len(cells) != 8 {
		t.Fatalf("got %d cells", len(cells))
	}
	if cells[0].Class != MemoryUser || cells[0].Intensity != 3 {
		t.Fatalf("first=%+v", cells[0])
	}
}
