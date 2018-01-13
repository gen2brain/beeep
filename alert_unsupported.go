// +build !linux,!windows,!darwin,!js

package beeep

import (
	"errors"
	"runtime"
)

// Alert displays a desktop notification and plays a beep.
func Alert(title, message, appIcon string) error {
	return errors.New("beeep: unsupported operating system: %v", runtime.GOOS)
}
