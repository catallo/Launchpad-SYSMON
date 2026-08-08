package monitor

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
)

// Snapshot contains the current percentage for each logical CPU thread.
type Snapshot struct {
	Threads []float64
}

// Read samples utilization for all logical CPU threads. Shinobi has eight threads.
func Read() (Snapshot, error) {
	threads, err := cpu.Percent(250*time.Millisecond, true)
	if err != nil {
		return Snapshot{}, fmt.Errorf("CPU threads: %w", err)
	}
	if len(threads) < 8 {
		return Snapshot{}, fmt.Errorf("expected at least 8 logical CPU threads, got %d", len(threads))
	}
	return Snapshot{Threads: threads[:8]}, nil
}
