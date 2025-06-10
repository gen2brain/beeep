//go:build windows && !linux && !freebsd && !netbsd && !openbsd && !darwin && !js

package beeep

import (
	"time"

	"git.sr.ht/~jackmordaunt/go-toast"
	"github.com/tadvi/systray"
	"golang.org/x/sys/windows/registry"
)

var isWindows10 bool

func init() {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion`, registry.QUERY_VALUE)
	if err != nil {
		return
	}
	defer k.Close()

	maj, _, err := k.GetIntegerValue("CurrentMajorVersionNumber")
	if err != nil {
		return
	}

	isWindows10 = maj == 10
}

// Notify sends desktop notification.
func Notify(title, message, icon string) error {
	if isWindows10 {
		if err := toastNotify(title, message, icon, false); err != nil {
			return err
		}
	} else {
		if err := balloonNotify(title, message, icon); err != nil {
			return err
		}
	}

	time.Sleep(time.Millisecond * 10)

	return nil
}

func balloonNotify(title, message, icon string) error {
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

	return tray.ShowMessage(title, message, false)
}

func toastNotify(title, message, icon string, urgent bool) error {
	n := toast.Notification{
		AppID: AppName,
		Title: title,
		Body:  message,
		Icon:  pathAbs(icon),
	}

	n.Duration = toast.Short
	if timeout.Seconds() > 10 {
		n.Duration = toast.Long
	}

	if urgent {
		n.Audio = toast.Default
	} else {
		n.Audio = toast.Silent
	}

	return n.Push()
}
