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
		s, err := monitor.ReadWithNetwork(network, *interval)
		if err != nil {
			log.Printf("read metrics: %v", err)
		} else if err := render(out, s); err != nil {
			log.Fatal(err)
		}
		// CPU sampling itself occupies the configured update interval.
		select {
		case <-signals:
			return
		default:
		}
	}
}

func render(out output, s monitor.Snapshot) error {
	if len(s.Threads) != 8 {
		return fmt.Errorf("expected 8 logical CPU-thread values, got %d", len(s.Threads))
	}
	for thread, value := range s.Threads {
		column, firstRow := launchpad.ThreadBlock(thread)
		full, partial := monitor.FineBar(value, 4)
		for offset := 0; offset < 4; offset++ {
			led := launchpad.Off
			if offset >= 4-full {
				led = cpuColor(value, 3)
			}
			if partial > 0 && offset == 3-full {
				led = cpuColor(value, partial)
			}
			if _, err := out.Write(launchpad.Message(launchpad.GridNote(firstRow+offset, column), led)); err != nil {
				return err
			}
		}
	}
	for row, cell := range monitor.FineMemoryCells(s.Memory, 8) {
		if _, err := out.Write(launchpad.Message(launchpad.GridNote(7-row, 4), memoryColor(cell.Class, cell.Intensity))); err != nil {
			return err
		}
	}
	swapPercent := 0.0
	if s.SwapTotal > 0 {
		swapPercent = float64(s.SwapUsed) * 100 / float64(s.SwapTotal)
	}
	if err := renderFineBar(out, func(row int) byte { return launchpad.GridNote(row, 5) }, swapPercent, swapColor); err != nil {
		return err
	}
	if err := renderFineBar(out, func(row int) byte { return launchpad.GridNote(row, 6) }, s.Network.DownloadMbit/280*100, cpuColor); err != nil {
		return err
	}
	if err := renderFineBar(out, func(row int) byte { return launchpad.GridNote(row, 7) }, s.Network.UploadMbit/50*100, cpuColor); err != nil {
		return err
	}
	if err := renderFineBar(out, launchpad.SideButtonNote, s.CPUTemperature/90*100, cpuColor); err != nil {
		return err
	}
	return nil
}

func cpuColor(percent float64, intensity byte) byte {
	// Normal legacy flags (0x0C); red and green each have three intensity levels.
	var red, green byte
	switch {
	case percent < 75:
		green = intensity
	case percent < 90:
		red, green = intensity, intensity
	default:
		red = intensity
	}
	return launchpad.Color(red, green)
}

func renderFineBar(out output, note func(int) byte, percent float64, color func(float64, byte) byte) error {
	full, partial := monitor.FineBar(percent, 8)
	for row := 0; row < 8; row++ {
		led := launchpad.Off
		if row >= 8-full {
			led = color(percent, 3)
		}
		if partial > 0 && row == 7-full {
			led = color(percent, partial)
		}
		if _, err := out.Write(launchpad.Message(note(row), led)); err != nil {
			return err
		}
	}
	return nil
}

func swapColor(percent float64, intensity byte) byte {
	switch {
	case percent <= 0:
		return launchpad.Off
	case percent < 50:
		return launchpad.Color(0, intensity)
	case percent < 80:
		return launchpad.Color(intensity, intensity)
	default:
		return launchpad.Color(intensity, 0)
	}
}

func memoryColor(class monitor.MemoryClass, intensity byte) byte {
	switch class {
	case monitor.MemoryUser:
		return launchpad.Color(intensity, intensity)
	case monitor.MemorySystem:
		return launchpad.Color(intensity, 0)
	case monitor.MemoryCache:
		return launchpad.Color(0, intensity)
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
