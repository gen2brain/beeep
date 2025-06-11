//go:build windows && !linux && !freebsd && !netbsd && !openbsd && !darwin && !js

package beeep

import (
	"fmt"
	"os"
	"time"
)

// Alert displays a desktop notification and plays a default system sound.
func Alert(title, message string, icon any) error {
	var img string
	switch i := icon.(type) {
	case string:
		img = i
	case []byte:
		var err error
		img, err = bytesToFilename(i)
		if err != nil {
			return err
		}
		defer os.Remove(img)
	default:
		return fmt.Errorf("unsupported argument: %T", icon)
	}

	if isWindows10 {
		if err := toastNotify(title, message, img, true); err != nil {
			return err
		}

		time.Sleep(time.Millisecond * 100)
	} else {
		if err := balloonNotify(title, message, img); err != nil {
			return err
		}
	}

	messageBeep()

	return nil
}
