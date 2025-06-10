//go:build windows && !linux && !freebsd && !netbsd && !openbsd && !darwin && !js

package beeep

import "time"

// Alert displays a desktop notification and plays a default system sound.
func Alert(title, message, icon string) error {
	if isWindows10 {
		if err := toastNotify(title, message, icon, true); err != nil {
			return err
		}

		time.Sleep(time.Millisecond * 10)
	} else {
		if err := balloonNotify(title, message, icon, true); err != nil {
			return err
		}
	}

	return nil
}
