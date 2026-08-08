package monitor

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
)

// Snapshot contains current utilization for the four physical CPU cores.
type Snapshot struct {
	Cores               []float64
	Memory              MemorySnapshot
	SwapUsed, SwapTotal uint64
	Network             NetworkRate
	CPUTemperature      float64
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
	cores, err := PhysicalCores(threads)
	if err != nil {
		return Snapshot{}, err
	}
	temperature, err := ReadCPUTemperature()
	if err != nil {
		return Snapshot{}, fmt.Errorf("CPU temperature: %w", err)
	}
	snapshot := Snapshot{Cores: cores, Memory: memory, SwapUsed: swapTotal - swapFree, SwapTotal: swapTotal, CPUTemperature: temperature}
	if network != nil {
		rate, err := network.Sample()
		if err != nil {
			return Snapshot{}, fmt.Errorf("network: %w", err)
		}
		snapshot.Network = rate
	}
	return snapshot, nil
}

// PhysicalCores averages the two logical threads of each physical core on Shinobi:
// CPU 0+4, 1+5, 2+6, and 3+7 according to the verified lscpu topology.
func PhysicalCores(threads []float64) ([]float64, error) {
	if len(threads) < 8 {
		return nil, fmt.Errorf("expected 8 logical CPU threads, got %d", len(threads))
	}
	return []float64{(threads[0] + threads[4]) / 2, (threads[1] + threads[5]) / 2, (threads[2] + threads[6]) / 2, (threads[3] + threads[7]) / 2}, nil
}

func readSwap() (total, free uint64, err error) {
	values, err := readMemInfo("/proc/meminfo")
	if err != nil {
		return 0, 0, err
	}
	return values["SwapTotal"], values["SwapFree"], nil
}
