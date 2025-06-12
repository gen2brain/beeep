//go:build (linux || freebsd || netbsd || openbsd || illumos) && !nodbus

package beeep

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"

	"github.com/esiqveland/notify"
	"github.com/godbus/dbus/v5"
)

// Notify sends desktop notification.
// The icon can be string with a path to png file or png []byte data. Stock icon names can also be used where supported.
//
// On Linux it tries to send notification via D-Bus, and it will fall back to `notify-send`.
//
// On macOS, this will first try `terminal-notifier` and will fall back to AppleScript with `osascript`.
//
// On Windows 10/11 it will use Windows Runtime COM API and will fall back to PowerShell. Windows 7 will use win32 API.
//
// On the Web it uses the Notification API, in Firefox it just works, in Chrome you must call it from some "user gesture"
// like `onclick`, and you must use TLS certificate, it doesn't work with plain http.
func Notify(title, message string, icon any) error {
	return notify1(title, message, icon, false)
}

func notify1(title, message string, ico any, urgent bool) error {
	var isString, isBytes bool
	switch ico.(type) {
	case string:
		isString = true
	case []byte:
		isBytes = true
	default:
		return fmt.Errorf("unsupported argument: %T", ico)
	}

	var icon string
	if isString {
		icon = ico.(string)
		if _, err := os.Stat(icon); err == nil {
			icon = pathAbs(icon)
		}
	}

	cmd1 := func() error {
		cmd, err := exec.LookPath("notify-send")
		if err != nil {
			return err
		}

		if isBytes {
			tmp, err := bytesToFilename(ico.([]byte))
			if err != nil {
				return err
			}
			defer os.Remove(tmp)

			icon = tmp
		}

		args := []string{title, message, "-a", AppName, "-i", icon, "-t", strconv.Itoa(int(timeout.Milliseconds())), "-u"}
		if urgent {
			args = append(args, "critical")
		} else {
			args = append(args, "normal")
		}
		c := exec.Command(cmd, args...)

		return c.Run()
	}

	cmd2 := func() error {
		cmd, err := exec.LookPath("kdialog")
		if err != nil {
			return err
		}

		if isBytes {
			tmp, err := bytesToFilename(ico.([]byte))
			if err != nil {
				return err
			}
			defer os.Remove(tmp)

			icon = tmp
		}

		args := []string{"--title", title, "--passivepopup", message, strconv.Itoa(int(timeout.Seconds())), "--icon", icon}
		c := exec.Command(cmd, args...)

		return c.Run()
	}

	dbus1 := func() error {
		conn, err := dbus.SessionBus()
		if err != nil {
			return err
		}
		defer conn.Close()

		n := notify.Notification{
			AppName:       AppName,
			AppIcon:       icon,
			Summary:       title,
			Body:          message,
			ExpireTimeout: timeout,
		}

		n.Hints = map[string]dbus.Variant{}

		if urgent {
			soundHint := notify.HintSoundWithName("bell")
			n.Hints[soundHint.ID] = soundHint.Variant
			n.SetUrgency(notify.UrgencyCritical)
		} else {
			n.SetUrgency(notify.UrgencyNormal)
		}

		if isBytes {
			rgba, err := bytesToRGBA(ico.([]byte))
			if err != nil {
				return err
			}

			imageHint := notify.HintImageDataRGBA(rgba)
			n.Hints[imageHint.ID] = imageHint.Variant
		}

		notifier, err := notify.New(conn, notify.WithLogger(log.New(io.Discard, "", log.Flags())))
		if err != nil {
			return err
		}
		defer notifier.Close()

		_, err = notifier.SendNotification(n)
		if err != nil {
			return err
		}

		return nil
	}

	err := dbus1()
	if err != nil {
		err1 := cmd1()
		if err1 != nil {
			err2 := cmd2()
			if err2 != nil {
				return fmt.Errorf("beeep: dbus: %w; notify-send: %w; kdialog: %w", err, err1, err2)
			}
		}
	}

	return nil
}

func bytesToRGBA(data []byte) (*image.RGBA, error) {
	i, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	if img, ok := i.(*image.RGBA); ok {
		return img, nil
	}

	b := i.Bounds()
	img := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(img, img.Bounds(), i, b.Min, draw.Src)

	return img, nil
}
