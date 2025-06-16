//go:build (linux || freebsd || netbsd || openbsd || illumos) && nodbus

package beeep

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

// Notify sends desktop notification.
// The icon can be string with a path to png file or png []byte data. Stock icon names can also be used where supported.
//
// On Linux it tries to send notification via D-Bus, and it will fall back to `notify-send` binary.
func Notify(title, message string, icon any) error {
	return notify1(title, message, icon, false)
}

func notify1(title, message string, icon any, urgent bool) error {
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

	if _, err := os.Stat(img); err == nil {
		img = pathAbs(img)
	}

	cmd1 := func() error {
		cmd, err := exec.LookPath("notify-send")
		if err != nil {
			return err
		}

		args := []string{title, message, "-a", AppName, "-i", img, "-t", strconv.Itoa(int(timeout.Milliseconds())), "-u"}
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

		args := []string{"--title", title, "--passivepopup", message, strconv.Itoa(int(timeout.Seconds())), "--icon", img}
		c := exec.Command(cmd, args...)

		return c.Run()
	}

	err1 := cmd1()
	if err1 != nil {
		err2 := cmd2()
		if err2 != nil {
			return fmt.Errorf("beeep: notify-send: %w; kdialog: %w", err1, err2)
		}
	}

	return nil
}
