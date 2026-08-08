// launchpad-control displays Shinobi CPU-thread status on a Novation Launchpad Mini.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"launchpad-sysmon/internal/launchpad"
	"launchpad-sysmon/internal/monitor"
)

// This is the same raw ALSA MIDI device used by the proven launchpad-text app.
const defaultDevice = "/dev/snd/midiC6D0"

type output interface{ Write([]byte) (int, error) }

func main() {
	device := flag.String("device", defaultDevice, "raw ALSA MIDI device")
	interval := flag.Duration("interval", time.Second, "refresh period")
	demo := flag.Bool("demo", false, "print metrics, do not send MIDI")
	clear := flag.Bool("clear", false, "turn all LEDs off and exit")
	flag.Parse()
	log.SetFlags(log.Ldate | log.Ltime)

	if *demo {
		runDemo(*interval)
		return
	}
	out, err := os.OpenFile(*device, os.O_WRONLY, 0)
	if err != nil {
		log.Fatalf("open MIDI device %q: %v", *device, err)
	}
	defer out.Close()
	if *clear {
		clearAll(out)
		return
	}
	clearAll(out)
	defer clearAll(out)
	network, err := monitor.NewNetworkSampler("enp0s31f6")
	if err != nil {
		log.Fatalf("network sampler: %v", err)
	}
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(signals)
	log.Printf("Monitoring via raw MIDI device %q; Ctrl+C clears LEDs.", *device)
	for {
		s, err := monitor.ReadWithNetwork(network)
		if err != nil {
			log.Printf("read metrics: %v", err)
		} else if err := render(out, s); err != nil {
			log.Fatal(err)
		}
		select {
		case <-signals:
			return
		case <-time.After(*interval):
		}
	}
}

func render(out output, s monitor.Snapshot) error {
	if len(s.Threads) != 8 {
		return fmt.Errorf("expected 8 logical CPU-thread values, got %d", len(s.Threads))
	}
	for thread, value := range s.Threads {
		column, firstRow := launchpad.ThreadBlock(thread)
		active, color := monitor.Bar(value, 4), monitor.Color(value)
		for offset := 0; offset < 4; offset++ {
			led := launchpad.Off
			if offset >= 4-active {
				led = color
			}
			if _, err := out.Write(launchpad.Message(launchpad.GridNote(firstRow+offset, column), led)); err != nil {
				return err
			}
		}
	}
	for row, class := range monitor.MemoryCells(s.Memory, 8) {
		if _, err := out.Write(launchpad.Message(launchpad.GridNote(7-row, 4), memoryColor(class))); err != nil {
			return err
		}
	}
	for row := 0; row < 8; row++ {
		led := launchpad.Off
		if s.SwapTotal > 0 && row >= 8-monitor.Bar(float64(s.SwapUsed)*100/float64(s.SwapTotal), 8) {
			led = swapColor(s.SwapUsed, s.SwapTotal)
		}
		if _, err := out.Write(launchpad.Message(launchpad.GridNote(row, 5), led)); err != nil {
			return err
		}
	}
	if err := renderNetworkBar(out, 6, s.Network.DownloadMbit, 280); err != nil {
		return err
	}
	if err := renderNetworkBar(out, 7, s.Network.UploadMbit, 50); err != nil {
		return err
	}
	if err := renderTemperatureBar(out, s.CPUTemperature, 85); err != nil {
		return err
	}
	return nil
}

func renderNetworkBar(out output, column int, rate, fullScale float64) error {
	active := monitor.Bar(rate/fullScale*100, 8)
	color := monitor.Color(rate / fullScale * 100)
	for row := 0; row < 8; row++ {
		led := launchpad.Off
		if row >= 8-active {
			led = color
		}
		if _, err := out.Write(launchpad.Message(launchpad.GridNote(row, column), led)); err != nil {
			return err
		}
	}
	return nil
}

func renderTemperatureBar(out output, temperature, fullScale float64) error {
	active := monitor.Bar(temperature/fullScale*100, 8)
	color := monitor.Color(temperature / fullScale * 100)
	for row := 0; row < 8; row++ {
		led := launchpad.Off
		if row >= 8-active {
			led = color
		}
		if _, err := out.Write(launchpad.Message(launchpad.SideButtonNote(row), led)); err != nil {
			return err
		}
	}
	return nil
}

func swapColor(used, total uint64) byte {
	if total == 0 || used == 0 {
		return launchpad.Off
	}
	percent := float64(used) * 100 / float64(total)
	if percent < 50 {
		return launchpad.Green
	}
	if percent < 80 {
		return launchpad.Amber
	}
	return launchpad.Red
}

func memoryColor(class monitor.MemoryClass) byte {
	switch class {
	case monitor.MemoryUser:
		return launchpad.Amber
	case monitor.MemorySystem:
		return launchpad.Red
	case monitor.MemoryCache:
		return launchpad.Green
	default:
		return launchpad.Off
	}
}

func clearAll(out output) {
	// Matrix plus any outer button left on by the earlier incorrect implementation.
	for row := 0; row < 8; row++ {
		for col := 0; col < 9; col++ {
			_, _ = out.Write(launchpad.Message(byte(row*16+col), launchpad.Off))
		}
	}
	for controller := byte(104); controller <= 111; controller++ {
		_, _ = out.Write([]byte{0xB0, controller, launchpad.Off})
	}
	log.Print("Launchpad LEDs cleared.")
}

func runDemo(interval time.Duration) {
	for {
		s, err := monitor.Read()
		if err != nil {
			log.Print(err)
		} else {
			fmt.Printf("CPU threads: %s\nRAM: user %.1f GiB | system %.1f GiB | cache %.1f GiB | free %.1f GiB\nSwap: %.1f / %.1f GiB | Network ↓ %.1f Mbit/s | ↑ %.1f Mbit/s | CPU %.1f °C\n", formatThreads(s.Threads), kibToGiB(s.Memory.User), kibToGiB(s.Memory.System), kibToGiB(s.Memory.Cache), kibToGiB(s.Memory.Free), kibToGiB(s.SwapUsed), kibToGiB(s.SwapTotal), s.Network.DownloadMbit, s.Network.UploadMbit, s.CPUTemperature)
		}
		if interval <= 0 {
			return
		}
		time.Sleep(interval)
	}
}
func formatThreads(values []float64) string {
	parts := make([]string, len(values))
	for i, value := range values {
		parts[i] = fmt.Sprintf("T%d %5.1f%%", i+1, value)
	}
	return strings.Join(parts, " | ")
}

func kibToGiB(value uint64) float64 { return float64(value) / 1024 / 1024 }
