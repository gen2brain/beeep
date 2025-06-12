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

func notify1(title, message, icon string, urgent bool) error {
	if _, err := os.Stat(icon); err == nil {
		icon = pathAbs(icon)
	}

	cmd1 := func() error {
		cmd, err := exec.LookPath("notify-send")
		if err != nil {
			return err
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

		args := []string{"--title", title, "--passivepopup", message, strconv.Itoa(int(timeout.Seconds())), "--icon", icon}
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
