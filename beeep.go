// Package beeep provides a cross-platform library for sending desktop notifications and beeps.
package beeep

import (
	"errors"
	"runtime"
)

var (
	// ErrUnsupported is returned when operating system is not supported.
	ErrUnsupported = errors.New("beeep: unsupported operating system: " + runtime.GOOS)
)
