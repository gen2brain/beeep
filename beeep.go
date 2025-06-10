// Package beeep provides a cross-platform library for sending desktop notifications and beeps.
package beeep

import (
	"fmt"
	"path/filepath"
	"runtime"
	"time"
)

var (
	// ErrUnsupported is returned when an operating system is not supported.
	ErrUnsupported = fmt.Errorf("beeep: unsupported operating system: %s", runtime.GOOS)
)

// AppID .
var AppID string

var timeout = time.Second * 5

func pathAbs(path string) string {
	var err error
	var abs string

	if path != "" {
		abs, err = filepath.Abs(path)
		if err != nil {
			abs = path
		}
	}

	return abs
}
