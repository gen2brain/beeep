//go:build (linux && nodbus) || (freebsd && nodbus) || (netbsd && nodbus) || (openbsd && nodbus) || illumos
// +build linux,nodbus freebsd,nodbus netbsd,nodbus openbsd,nodbus illumos

package beeep

import (
	"errors"
	"os/exec"
	"strconv"
	"time"
)

// Notify sends desktop notification with default timeout
func Notify(title, message, appIcon string) error {
	return NotifyEx(title, message, appIcon, -1)
}

// NotifyEx sends desktop notification with timeout
func NotifyEx(title, message, appIcon string, timeout time.Duration) error {
	appIcon = pathAbs(appIcon)
	if timeout > 0 {
		timeout = timeout / time.Millisecond
	}
	timeOut := strconv.FormatInt(int64(timeout), 10)

	cmd := func() error {
		send, err := exec.LookPath("sw-notify-send")
		if err != nil {
			send, err = exec.LookPath("notify-send")
			if err != nil {
				return err
			}
		}

		c := exec.Command(send, title, message, "-i", appIcon, "-t", timeOut)
		return c.Run()
	}

	knotify := func() error {
		send, err := exec.LookPath("kdialog")
		if err != nil {
			return err
		}
		c := exec.Command(send, "--title", title, "--passivepopup", message, strconv.Itoa(timeout/1e3), "--icon", appIcon)
		return c.Run()
	}

	err := cmd()
	if err != nil {
		e := knotify()
		if e != nil {
			return errors.New("beeep: " + err.Error() + "; " + e.Error())
		}
	}

	return nil
}
