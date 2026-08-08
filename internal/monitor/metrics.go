package monitor

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/mem"
)

// Snapshot contains percentages rendered on the monitor grid.
type Snapshot struct{ CPU, RAM, RootDisk, GPU, VRAM float64 }

func Read() (Snapshot, error) {
	cpuPercent, err := cpu.Percent(0, false)
	if err != nil || len(cpuPercent) == 0 {
		return Snapshot{}, fmt.Errorf("CPU: %w", err)
	}
	memory, err := mem.VirtualMemory()
	if err != nil {
		return Snapshot{}, fmt.Errorf("RAM: %w", err)
	}
	root, err := disk.Usage("/")
	if err != nil {
		return Snapshot{}, fmt.Errorf("root disk: %w", err)
	}
	gpu, vram := nvidia()
	return Snapshot{CPU: cpuPercent[0], RAM: memory.UsedPercent, RootDisk: root.UsedPercent, GPU: gpu, VRAM: vram}, nil
}

func nvidia() (float64, float64) {
	output, err := exec.Command("nvidia-smi", "--query-gpu=utilization.gpu,memory.used,memory.total", "--format=csv,noheader,nounits").Output()
	if err != nil {
		return 0, 0
	}
	fields := strings.Split(strings.TrimSpace(string(output)), ",")
	if len(fields) != 3 {
		return 0, 0
	}
	gpu, e1 := strconv.ParseFloat(strings.TrimSpace(fields[0]), 64)
	used, e2 := strconv.ParseFloat(strings.TrimSpace(fields[1]), 64)
	total, e3 := strconv.ParseFloat(strings.TrimSpace(fields[2]), 64)
	if e1 != nil || e2 != nil || e3 != nil || total <= 0 {
		return 0, 0
	}
	return gpu, used / total * 100
}
