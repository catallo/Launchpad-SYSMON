package monitor

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
)

// Snapshot contains current utilization for eight logical CPU threads.
type Snapshot struct {
	Threads             []float64
	Memory              MemorySnapshot
	SwapUsed, SwapTotal uint64
	Network             NetworkRate
	CPUTemperature      float64
}

// Read samples utilization for all logical CPU threads. Shinobi has eight threads.
func Read() (Snapshot, error) {
	return ReadWithNetwork(nil, 500*time.Millisecond)
}

// ReadWithNetwork samples CPU and Ethernet over sampleInterval.
func ReadWithNetwork(network *NetworkSampler, sampleInterval time.Duration) (Snapshot, error) {
	threads, err := cpu.Percent(sampleInterval, true)
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
	logicalThreads, err := LogicalThreads(threads)
	if err != nil {
		return Snapshot{}, err
	}
	temperature, err := ReadCPUTemperature()
	if err != nil {
		return Snapshot{}, fmt.Errorf("CPU temperature: %w", err)
	}
	snapshot := Snapshot{Threads: logicalThreads, Memory: memory, SwapUsed: swapTotal - swapFree, SwapTotal: swapTotal, CPUTemperature: temperature}
	if network != nil {
		rate, err := network.Sample()
		if err != nil {
			return Snapshot{}, fmt.Errorf("network: %w", err)
		}
		snapshot.Network = rate
	}
	return snapshot, nil
}

// LogicalThreads returns the eight logical CPU threads exposed by Shinobi.
func LogicalThreads(threads []float64) ([]float64, error) {
	if len(threads) < 8 {
		return nil, fmt.Errorf("expected 8 logical CPU threads, got %d", len(threads))
	}
	return append([]float64(nil), threads[:8]...), nil
}

func readSwap() (total, free uint64, err error) {
	values, err := readMemInfo("/proc/meminfo")
	if err != nil {
		return 0, 0, err
	}
	return values["SwapTotal"], values["SwapFree"], nil
}
