//go:build windows && !go1.21 && !linux && !freebsd && !netbsd && !openbsd && !darwin && !js

package beeep

import (
	"time"

	"github.com/tadvi/systray"
)

var isWindows10 = false

// Notify sends desktop notification.
func Notify(title, message, icon string) error {
	if err := balloonNotify(title, message, icon, false); err != nil {
		return err
	}

	return nil
}

func balloonNotify(title, message, icon string, urgent bool) error {
	tray, err := systray.New()
	if err != nil {
		return err
	}

	err = tray.ShowCustom(pathAbs(icon), title)
	if err != nil {
		return err
	}

	go func() {
		go func() {
			_ = tray.Run()
		}()
		time.Sleep(timeout)
		_ = tray.Stop()
	}()

	err = tray.ShowMessage(title, message, false)
	if err != nil {
		return err
	}

	if urgent {
		err = Beep(DefaultFreq, DefaultDuration)
		if err != nil {
			return err
		}
	}

	return nil
}

func toastNotify(title, message, icon string, urgent bool) error {
	return nil
}
