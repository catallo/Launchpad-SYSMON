# Launchpad-SYSMON

A small Go system monitor for a Novation Launchpad Mini on Linux.

The 8×8 pad grid shows:

- columns 1–4: physical CPU-core usage
- column 5: RAM split into application memory, kernel/shared memory, file cache, and free memory
- column 6: swap usage
- column 7: download rate, scaled to 280 Mbit/s
- column 8: upload rate, scaled to 50 Mbit/s

Green, amber, and red indicate low, medium, and high utilization. The monitor uses the legacy two-colour MIDI protocol and writes directly to an ALSA raw-MIDI device.

## Build

```bash
go test ./...
go build -o bin/launchpad-sysmon ./cmd/launchpad-control
```

## Run

```bash
./bin/launchpad-sysmon --interval=500ms
```

The default device is `/dev/snd/midiC6D0`; adjust it with `--device` if your Launchpad uses another ALSA MIDI device.

## Service

The included user service starts the monitor with a 500 ms refresh interval:

```bash
systemctl --user enable --now launchpad-sysmon.service
```
