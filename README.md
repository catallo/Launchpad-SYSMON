# Launchpad-SYSMON

A small Go system monitor for a Novation Launchpad Mini on Linux.

The 8×8 pad grid shows:

- columns 1–4: usage of eight logical CPU threads
- column 5: RAM split into application memory, kernel/shared memory, file cache, and free memory
- column 6: swap usage
- column 7: download rate, scaled to 280 Mbit/s
- column 8: upload rate, scaled to 50 Mbit/s

Green, amber, and red indicate low, medium, and high utilization. The monitor uses the legacy two-colour MIDI protocol and writes directly to an ALSA raw-MIDI device.

## Current hardware assumptions

This application is currently tailored to the Shinobi workstation, not a general-purpose hardware monitor. It requires Linux, at least eight logical CPU threads, a Linux `coretemp` CPU-temperature sensor, and the hard-coded network interface `enp0s31f6`. The first eight CPU threads are displayed. The default raw MIDI device is `/dev/snd/midiC6D0`; override that path with `--device`.

The published binary targets **Linux x86-64 (amd64)**. Other operating systems and CPU architectures are not advertised as supported. Automated tests and a CLI-help smoke test do not replace a real Launchpad/hardware test.

## Install a pre-built binary

Download `launchpad-sysmon-linux-amd64` and `SHA256SUMS` from the same [Release](https://github.com/catallo/Launchpad-SYSMON/releases), then verify and install them. Installation only proceeds if checksum verification succeeds:

```bash
sha256sum --check SHA256SUMS &&
  install -D -m 755 launchpad-sysmon-linux-amd64 "$HOME/.local/bin/launchpad-sysmon"
```

Run the installed binary by its full path, or add `$HOME/.local/bin` to your `PATH`.

## Build from source

Use the Go version declared in `go.mod` or a newer supported Go release:

```bash
go test ./...
go build -o bin/launchpad-sysmon ./cmd/launchpad-sysmon
```

The `bin/` and `dist/` directories are ignored by Git. Compiled programs are delivered as build artifacts and Release assets, not committed to source. The previously tracked binary remains available in Git history.

## Automated builds and releases

[GitHub Actions](https://github.com/catallo/Launchpad-SYSMON/actions/workflows/build.yml) runs unit tests with the race detector, `go vet`, a CGO-disabled Linux x86-64 build, a hardware-free `--help` smoke test, and SHA-256 checksum verification.

It runs on pushes and pull requests to `master`, on version tags matching `v*`, and manually through the Actions page. Successful builds provide a `launchpad-sysmon-binaries` artifact containing the binary and `SHA256SUMS`. Actions artifacts expire after seven days and require a GitHub login to download; ZIP artifacts do not preserve executable permissions.

To publish a release, first push the tested commit to `master`, then push a **new** semantic version tag such as `v0.1.0` in a separate push. The tag must point to a commit in `master` history. GitHub rebuilds and tests that tagged source, then publishes the binary, checksums, and GitHub-generated release notes. Tags with a prerelease suffix are marked as prereleases; build metadata alone does not mark a prerelease. Existing releases are never overwritten. Branch builds and manual runs do not publish releases.

No personal access token or other repository secret is needed; only the release job receives write access to repository contents.

## Run

For a local source build:

```bash
./bin/launchpad-sysmon --interval=500ms
```

For the installed release binary:

```bash
"$HOME/.local/bin/launchpad-sysmon" --interval=500ms --device=/dev/snd/midiC6D0
```

Normal operation writes to the MIDI device and changes the Launchpad LEDs. To inspect flags without opening that device:

```bash
"$HOME/.local/bin/launchpad-sysmon" --help
```

## Service

No systemd unit is currently shipped in this repository. If you already have a local `launchpad-sysmon.service`, keep its executable path and MIDI-device permissions in mind when installing a release. The build/release workflow does not install software, modify services, or access a connected Launchpad.
