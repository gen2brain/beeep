//go:build windows && !linux && !freebsd && !netbsd && !openbsd && !darwin && !js

package beeep

import (
	"syscall"
)

const (
	// DefaultFreq - frequency, in Hz, middle A
	DefaultFreq = 587.0
	// DefaultDuration - duration in milliseconds
	DefaultDuration = 500
)

// Beep beeps the PC speaker (https://en.wikipedia.org/wiki/PC_speaker).
//
// On Windows it uses Beep function via syscall.
func Beep(freq float64, duration int) error {
	if freq == 0 {
		freq = DefaultFreq
	} else if freq > 32767 {
		freq = 32767
	} else if freq < 37 {
		freq = DefaultFreq
	}

	if duration == 0 {
		duration = DefaultDuration
	}

	kernel32, _ := syscall.LoadLibrary("kernel32.dll")
	beep32, _ := syscall.GetProcAddress(kernel32, "Beep")

	defer syscall.FreeLibrary(kernel32)

	_, _, e := syscall.SyscallN(beep32, uintptr(int(freq)), uintptr(duration))
	if e != 0 {
		return e
	}

	return nil
}

func messageBeep() {
	user32, _ := syscall.LoadLibrary("user32.dll")
	beep32, _ := syscall.GetProcAddress(user32, "MessageBeep")

	defer syscall.FreeLibrary(user32)

	_, _, _ = syscall.SyscallN(beep32, uintptr(0x00000000))
}
