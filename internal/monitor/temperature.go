package monitor

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// ReadCPUTemperature returns the coretemp CPU package temperature in Celsius.
func ReadCPUTemperature() (float64, error) {
	entries, err := filepath.Glob("/sys/class/hwmon/hwmon*")
	if err != nil {
		return 0, err
	}
	for _, entry := range entries {
		name, err := os.ReadFile(filepath.Join(entry, "name"))
		if err != nil || strings.TrimSpace(string(name)) != "coretemp" {
			continue
		}
		value, err := os.ReadFile(filepath.Join(entry, "temp1_input"))
		if err != nil {
			return 0, fmt.Errorf("read %s: %w", entry, err)
		}
		milliC, err := strconv.ParseFloat(strings.TrimSpace(string(value)), 64)
		if err != nil {
			return 0, fmt.Errorf("parse %s: %w", entry, err)
		}
		return milliC / 1000, nil
	}
	return 0, fmt.Errorf("coretemp package sensor not found")
}
