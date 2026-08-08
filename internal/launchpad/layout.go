// Package launchpad implements the legacy Novation Launchpad Mini MIDI layout.
package launchpad

import "fmt"

// LED colour values use the legacy Launchpad red/green velocity encoding.
const (
	Off   byte = 0x00
	Red   byte = 0x0F
	Green byte = 0x3C
	Amber byte = 0x3F
)

// GridNote returns the note for a matrix pad. Row 0 is the top row, col 0 the left column.
func GridNote(row, col int) byte {
	if row < 0 || row > 7 || col < 0 || col > 7 {
		panic("grid coordinate out of range")
	}
	return byte((8-row)*10 + col + 1)
}

// Color calculates a visible legacy colour. red and green range from 0 to 3.
func Color(red, green byte) byte {
	if red > 3 || green > 3 {
		panic("colour component out of range")
	}
	if red == 0 && green == 0 {
		return Off
	}
	return 0x0C + red + 0x10*green // normal copy-mode flags + red/green intensities
}

// Message creates a channel-1 note-on LED message.
func Message(note, color byte) []byte { return []byte{0x90, note, color} }

func ValidatePort(name string) error {
	if name == "" {
		return fmt.Errorf("empty MIDI port")
	}
	return nil
}
