package monitor

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
)

// Snapshot contains the current percentage for each logical CPU thread.
type Snapshot struct {
	Threads             []float64
	Memory              MemorySnapshot
	SwapUsed, SwapTotal uint64
	Network             NetworkRate
}

// Read samples utilization for all logical CPU threads. Shinobi has eight threads.
func Read() (Snapshot, error) {
	return ReadWithNetwork(nil)
}

// ReadWithNetwork includes a sampled Ethernet rate when a sampler is supplied.
func ReadWithNetwork(network *NetworkSampler) (Snapshot, error) {
	threads, err := cpu.Percent(250*time.Millisecond, true)
	if err != nil {
		return Snapshot{}, fmt.Errorf("CPU threads: %w", err)
	}
	if len(threads) < 8 {
		return Snapshot{}, fmt.Errorf("expected at least 8 logical CPU threads, got %d", len(threads))
	}
	memory, err := ReadMemory()
	if err != nil {
		return Snapshot{}, fmt.Errorf("memory: %w", err)
	}
	swapTotal, swapFree, err := readSwap()
	if err != nil {
		return Snapshot{}, fmt.Errorf("swap: %w", err)
	}
	snapshot := Snapshot{Threads: threads[:8], Memory: memory, SwapUsed: swapTotal - swapFree, SwapTotal: swapTotal}
	if network != nil {
		rate, err := network.Sample()
		if err != nil {
			return Snapshot{}, fmt.Errorf("network: %w", err)
		}
		snapshot.Network = rate
	}
	return snapshot, nil
}

func readSwap() (total, free uint64, err error) {
	values, err := readMemInfo("/proc/meminfo")
	if err != nil {
		return 0, 0, err
	}
	return values["SwapTotal"], values["SwapFree"], nil
}
