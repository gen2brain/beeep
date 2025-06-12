//go:build linux || freebsd || netbsd || openbsd || illumos

package beeep

import (
	"fmt"
	"os"
	"syscall"
	"time"
	"unsafe"
)

// Constants
const (
	// DefaultFreq - frequency, in Hz, middle A
	DefaultFreq = 440.0
	// DefaultDuration - duration in milliseconds
	DefaultDuration = 200
)

const (
	// linux/input-event-codes.h
	evSnd   = 0x12 // Event type
	sndTone = 0x02 // Sound
)

// inputEvent represents linux/input.h event structure.
type inputEvent struct {
	Time  syscall.Timeval // time in seconds since the epoch at which event occurred
	Type  uint16          // event type
	Code  uint16          // event code related to the event type
	Value int32           // event value related to the event type
}

// Beep beeps the PC speaker (https://en.wikipedia.org/wiki/PC_speaker).
//
// On Linux it needs permission to access `/dev/input/by-path/platform-pcspkr-event-spkr` file for writing,
// and `pcspkr` module must be loaded. User must be in the correct group, usually `input`.
//
// If it cannot open device files, it will fall back to sending Bell character (https://en.wikipedia.org/wiki/Bell_character).
// For bell character in X11 terminals you can enable a bell with `xset b on`. For console check `setterm` and `--blength` or `--bfreq` options.
//
// On macOS, it will first try to use `osascript` and will fall back to sending bell character.
// Enable `Audible bell` in Terminal --> Preferences --> Settings --> Advanced.
//
// On Windows it uses Beep function via syscall.
//
// On the Web it plays hard-coded beep sound.
func Beep(freq float64, duration int) error {
	if freq == 0 {
		freq = DefaultFreq
	} else if freq > 20000 {
		freq = 20000
	} else if freq < 0 {
		freq = DefaultFreq
	}

	if duration == 0 {
		duration = DefaultDuration
	}

	f, err := os.OpenFile("/dev/input/by-path/platform-pcspkr-event-spkr", os.O_WRONLY, 0644)
	if err != nil {
		// Output the only beep we can
		_, err = os.Stdout.Write([]byte{7})
		if err != nil {
			return fmt.Errorf("beeep: error writing to stdout: %w", err)
		}

		return nil
	}

	defer f.Close()

	ev := inputEvent{}
	ev.Type = evSnd
	ev.Code = sndTone
	ev.Value = int32(freq)

	d := *(*[unsafe.Sizeof(ev)]byte)(unsafe.Pointer(&ev))

	// Start beep
	_, err = f.Write(d[:])
	if err != nil {
		return fmt.Errorf("beeep: error writing to pcspkr: %w", err)
	}

	time.Sleep(time.Duration(duration) * time.Millisecond)

	ev.Value = 0
	d = *(*[unsafe.Sizeof(ev)]byte)(unsafe.Pointer(&ev))

	// Stop beep
	_, err = f.Write(d[:])
	if err != nil {
		return fmt.Errorf("beeep: error writing to pcspkr: %w", err)
	}

	return nil
}
