// +build !linux,!windows,!darwin,!js

package beeep

import (
	"errors"
	"runtime"
)

// Notify sends desktop notification.
func Notify(title, message string) error {
	return errors.New("beeep: unsupported operating system: %v", runtime.GOOS)
}
