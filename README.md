# Launchpad Control

Lokaler System-Monitor für das angeschlossene Novation Launchpad Mini auf **Shinobi**.

## Ziel

Die 8×8-Matrix zeigt Ressourcenbalken (CPU, RAM, Datenträger, GPU, VRAM, Temperaturen und Netzwerk). Die obere und rechte Tastenreihe werden später für Ansichten und Details verwendet.

## Sicherheitsprinzip

Das Programm steuert ausschließlich das Launchpad. Es startet, stoppt oder verändert keine Systemdienste. Hardwarezugriff erfolgt ausschließlich als Benutzer `sco`.

## Status

Grundgerüst und Tests vorhanden. Der MIDI-Treiber und das Layout folgen als nächster Schritt.

## Entwicklung

```bash
cd /home/sco/Projects/Launchpad-Control
go test ./...
```
