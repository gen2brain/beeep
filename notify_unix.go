//go:build (linux || freebsd || netbsd || openbsd || illumos) && !nodbus

package beeep

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"

	"github.com/esiqveland/notify"
	"github.com/godbus/dbus/v5"
)

// Notify sends desktop notification.
//
// On Linux it tries to send notification via D-Bus, and it will fall back to `notify-send` binary.
func Notify(title, message, icon string) error {
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

		args := []string{title, message, "-a", AppID, "-i", icon, "-t", strconv.Itoa(int(timeout.Milliseconds())), "-u"}
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

	dbus1 := func() error {
		conn, err := dbus.SessionBus()
		if err != nil {
			return err
		}
		defer conn.Close()

		n := notify.Notification{
			AppName:       AppID,
			AppIcon:       icon,
			Summary:       title,
			Body:          message,
			ExpireTimeout: timeout,
		}

		if urgent {
			soundHint := notify.HintSoundWithName("bell")
			n.Hints = map[string]dbus.Variant{
				soundHint.ID: soundHint.Variant,
			}
			n.SetUrgency(notify.UrgencyCritical)
		} else {
			n.SetUrgency(notify.UrgencyNormal)
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
