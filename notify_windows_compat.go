//go:build windows && !go1.21 && !linux && !freebsd && !netbsd && !openbsd && !darwin && !js

package beeep

import (
	"fmt"
	"image/png"
	"os"
	"time"

	"github.com/sergeymakinen/go-ico"
	"github.com/tadvi/systray"
)

var isWindows10 = false

// Notify sends desktop notification.
// The icon can be string with a path to png file or png []byte data. Stock icon names can also be used where supported.
//
// On Windows 10/11 it will use Windows Runtime COM API and will fall back to PowerShell. Windows 7 will use win32 API.
func Notify(title, message string, icon any) error {
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

	if err := balloonNotify(title, message, img); err != nil {
		return err
	}

	return nil
}

func balloonNotify(title, message, icon string) error {
	tray, err := systray.New()
	if err != nil {
		return err
	}

	tmp, err := pngToIco(pathAbs(icon))
	if err != nil {
		return err
	}
	defer os.Remove(tmp)

	err = tray.ShowCustom(tmp, title)
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

	err = tray.ShowMessage(title, message, true)
	if err != nil {
		return err
	}

	return nil
}

func toastNotify(title, message, icon string, urgent bool) error {
	return nil
}

func pngToIco(icon string) (string, error) {
	var out string

	f, err := os.Open(icon)
	if err != nil {
		return out, err
	}
	defer f.Close()

	img, err := png.Decode(f)
	if err != nil {
		return out, err
	}

	tmp, err := os.CreateTemp(os.TempDir(), "beeep")
	if err != nil {
		return out, err
	}
	defer tmp.Close()

	out = tmp.Name()

	err = ico.Encode(tmp, img)
	if err != nil {
		return out, err
	}

	return out, nil
}
