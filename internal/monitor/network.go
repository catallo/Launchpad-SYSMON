package monitor

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// NetworkRate is the current Ethernet throughput in megabits per second.
type NetworkRate struct{ DownloadMbit, UploadMbit float64 }

type NetworkSampler struct {
	iface          string
	received, sent uint64
	sampled        time.Time
}

func NewNetworkSampler(iface string) (*NetworkSampler, error) {
	rx, tx, err := readNetworkBytes("/proc/net/dev", iface)
	if err != nil {
		return nil, err
	}
	return &NetworkSampler{iface: iface, received: rx, sent: tx, sampled: time.Now()}, nil
}

func (s *NetworkSampler) Sample() (NetworkRate, error) {
	rx, tx, err := readNetworkBytes("/proc/net/dev", s.iface)
	if err != nil {
		return NetworkRate{}, err
	}
	now := time.Now()
	seconds := now.Sub(s.sampled).Seconds()
	if seconds <= 0 {
		return NetworkRate{}, fmt.Errorf("invalid sample interval")
	}
	rate := NetworkRate{DownloadMbit: float64(rx-s.received) * 8 / 1_000_000 / seconds, UploadMbit: float64(tx-s.sent) * 8 / 1_000_000 / seconds}
	s.received, s.sent, s.sampled = rx, tx, now
	return rate, nil
}

func readNetworkBytes(path, iface string) (uint64, uint64, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		parts := strings.Split(line, ":")
		if len(parts) != 2 || strings.TrimSpace(parts[0]) != iface {
			continue
		}
		fields := strings.Fields(parts[1])
		if len(fields) < 9 {
			return 0, 0, fmt.Errorf("bad counters for %s", iface)
		}
		rx, e1 := strconv.ParseUint(fields[0], 10, 64)
		tx, e2 := strconv.ParseUint(fields[8], 10, 64)
		if e1 != nil || e2 != nil {
			return 0, 0, fmt.Errorf("parse counters for %s", iface)
		}
		return rx, tx, nil
	}
	if err := scanner.Err(); err != nil {
		return 0, 0, err
	}
	return 0, 0, fmt.Errorf("interface %s not found", iface)
}
