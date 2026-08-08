package monitor

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

// MemoryClass is a non-overlapping physical-RAM category.
type MemoryClass uint8

const (
	MemoryFree   MemoryClass = iota
	MemoryUser               // anonymous memory of applications
	MemorySystem             // kernel memory and shared memory
	MemoryCache              // filesystem cache, buffers, reclaimable slab
)

// MemorySnapshot partitions physical memory in KiB. The classes add up to Total.
type MemorySnapshot struct {
	Total, User, System, Cache, Free uint64
}

// ReadMemory derives a useful, non-overlapping htop-style memory breakdown from /proc/meminfo.
func ReadMemory() (MemorySnapshot, error) {
	values, err := readMemInfo("/proc/meminfo")
	if err != nil {
		return MemorySnapshot{}, err
	}
	total := values["MemTotal"]
	if total == 0 {
		return MemorySnapshot{}, fmt.Errorf("MemTotal missing")
	}
	shmem := values["Shmem"]
	anon := subtract(values["Active(anon)"]+values["Inactive(anon)"], shmem)
	kernel := values["SUnreclaim"] + values["KernelStack"] + values["PageTables"] + values["Unevictable"] + shmem
	cache := values["Buffers"] + subtract(values["Cached"], shmem) + values["SReclaimable"]
	used := anon + kernel + cache
	if used > total {
		return MemorySnapshot{}, fmt.Errorf("memory categories exceed MemTotal")
	}
	return MemorySnapshot{Total: total, User: anon, System: kernel, Cache: cache, Free: total - used}, nil
}

func subtract(value, remove uint64) uint64 {
	if remove > value {
		return 0
	}
	return value - remove
}

func readMemInfo(path string) (map[string]uint64, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	result := make(map[string]uint64)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 2 {
			continue
		}
		key := strings.TrimSuffix(fields[0], ":")
		value, err := strconv.ParseUint(fields[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", key, err)
		}
		result[key] = value
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

// MemoryCells returns bottom-to-top class values for a whole 8-cell RAM bar.
// Tiny nonzero categories are kept visible where possible, so kernel usage is not hidden by rounding.
func MemoryCells(s MemorySnapshot, cells int) []MemoryClass {
	if cells <= 0 || s.Total == 0 {
		return nil
	}
	classes := []MemoryClass{MemoryUser, MemorySystem, MemoryCache, MemoryFree}
	amounts := []uint64{s.User, s.System, s.Cache, s.Free}
	counts := make([]int, len(classes))
	remaining := cells
	for i, amount := range amounts {
		if amount > 0 && remaining > 0 {
			counts[i] = 1
			remaining--
		}
	}
	for remaining > 0 {
		best := -1
		var bestNeed float64
		for i, amount := range amounts {
			if amount == 0 {
				continue
			}
			need := float64(amount)/float64(s.Total)*float64(cells) - float64(counts[i])
			if best == -1 || need > bestNeed {
				best, bestNeed = i, need
			}
		}
		counts[best]++
		remaining--
	}
	result := make([]MemoryClass, 0, cells)
	for i, class := range classes {
		for n := 0; n < counts[i]; n++ {
			result = append(result, class)
		}
	}
	return result
}

// MemoryCell is one three-level LED slot in the RAM bar.
type MemoryCell struct {
	Class     MemoryClass
	Intensity byte
}

// FineMemoryCells returns exactly pads RAM cells. Each category's final cell
// reflects its fractional share at one of three brightness levels.
func FineMemoryCells(s MemorySnapshot, pads int) []MemoryCell {
	classes := MemoryCells(s, pads)
	if len(classes) == 0 {
		return nil
	}
	amounts := map[MemoryClass]uint64{
		MemoryUser: s.User, MemorySystem: s.System, MemoryCache: s.Cache, MemoryFree: s.Free,
	}
	last := make(map[MemoryClass]int)
	for i, class := range classes {
		last[class] = i
	}
	result := make([]MemoryCell, len(classes))
	for i, class := range classes {
		result[i] = MemoryCell{Class: class, Intensity: 3}
	}
	for class, index := range last {
		if class == MemoryFree {
			result[index].Intensity = 0
			continue
		}
		exact := float64(amounts[class]) / float64(s.Total) * float64(pads)
		fraction := exact - math.Floor(exact)
		if fraction > 0 {
			result[index].Intensity = byte(math.Ceil(fraction * 3))
		}
	}
	return result
}
