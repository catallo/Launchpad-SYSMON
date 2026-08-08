// launchpad-control displays Shinobi system status on a Novation Launchpad Mini.
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

	"launchpad-control/internal/launchpad"
	"launchpad-control/internal/monitor"

	"gitlab.com/gomidi/midi/v2/drivers"
	"gitlab.com/gomidi/midi/v2/drivers/rtmididrv"
)

const defaultPort = "Launchpad Mini"

func main() {
	portName := flag.String("port", defaultPort, "part of the MIDI output port name")
	interval := flag.Duration("interval", time.Second, "refresh period")
	demo := flag.Bool("demo", false, "print metrics, do not send MIDI")
	clear := flag.Bool("clear", false, "turn all LEDs off and exit")
	flag.Parse()
	log.SetFlags(log.Ldate | log.Ltime)

	if *demo {
		runDemo(*interval)
		return
	}
	driver, err := rtmididrv.New()
	if err != nil {
		log.Fatal(err)
	}
	defer driver.Close()
	out, err := output(driver, *portName)
	if err != nil {
		log.Fatal(err)
	}
	if err := out.Open(); err != nil {
		log.Fatalf("open MIDI port %q: %v", out.String(), err)
	}
	defer out.Close()
	if *clear {
		clearAll(out)
		return
	}
	defer clearAll(out)
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(signals)
	log.Printf("Monitoring via MIDI port %q; Ctrl+C clears LEDs.", out.String())
	for {
		s, err := monitor.Read()
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

func output(driver *rtmididrv.Driver, wanted string) (drivers.Out, error) {
	ports, err := driver.Outs()
	if err != nil {
		return nil, err
	}
	for _, port := range ports {
		if strings.Contains(strings.ToLower(port.String()), strings.ToLower(wanted)) {
			return port, nil
		}
	}
	var names []string
	for _, p := range ports {
		names = append(names, p.String())
	}
	return nil, fmt.Errorf("MIDI output containing %q not found; available: %s", wanted, strings.Join(names, "; "))
}

func render(out drivers.Out, s monitor.Snapshot) error {
	values := []float64{s.CPU, s.RAM, s.RootDisk, s.GPU, s.VRAM}
	for col, value := range values {
		active := monitor.Bar(value, 8)
		color := monitor.Color(value)
		for row := 0; row < 8; row++ {
			c := launchpad.Off
			if row >= 8-active {
				c = color
			}
			if err := out.Send(launchpad.Message(launchpad.GridNote(row, col), c)); err != nil {
				return err
			}
		}
	}
	return nil
}

func clearAll(out drivers.Out) {
	for row := 0; row < 8; row++ {
		for col := 0; col < 8; col++ {
			_ = out.Send(launchpad.Message(launchpad.GridNote(row, col), launchpad.Off))
		}
	}
	log.Print("Launchpad LEDs cleared.")
}
func runDemo(interval time.Duration) {
	for {
		s, err := monitor.Read()
		if err != nil {
			log.Print(err)
		} else {
			fmt.Printf("CPU %5.1f%% | RAM %5.1f%% | Disk %5.1f%% | GPU %5.1f%% | VRAM %5.1f%%\n", s.CPU, s.RAM, s.RootDisk, s.GPU, s.VRAM)
		}
		if interval <= 0 {
			return
		}
		time.Sleep(interval)
	}
}
