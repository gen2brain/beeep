// Package beeep provides a cross-platform library for sending desktop notifications and beeps.
package beeep

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

var (
	// ErrUnsupported is returned when an operating system is not supported.
	ErrUnsupported = fmt.Errorf("beeep: unsupported operating system: %s", runtime.GOOS)
)

// AppName is the name of app.
// This should be the application's formal name, rather than some sort of ID.
var AppName = "DefaultAppName"

// timeout is notification duration (where applicable).
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

func bytesToFilename(data []byte) (string, error) {
	var out string

	tmp, err := os.CreateTemp(os.TempDir(), "beeep*.png")
	if err != nil {
		return out, err
	}
	defer tmp.Close()

	_, err = tmp.Write(data)
	if err != nil {
		return out, err
	}

	out = tmp.Name()

	return out, nil
}
