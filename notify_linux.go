// +build linux,!windows,!darwin,!js

package beeep

import (
	"errors"
	"os/exec"

	"github.com/godbus/dbus"
)

// Notify sends desktop notification.
//
// On Linux it tries to send notification via D-Bus and it will fallback to `notify-send` binary.
//
// FIXME: On Windows this just executes `msg` command, available since Vista.
//
// On macOS this executes AppleScript with `osascript` binary.
//
// On Web, in Firefox it just works, in Chrome you must call it from some "user gesture" like `onclick`,
// and you must use TLS certificate, it doesn't work with plain http.
func Notify(title, message string) error {
	cmd := func() error {
		send, err := exec.LookPath("sw-notify-send")
		if err != nil {
			send, err = exec.LookPath("notify-send")
			if err != nil {
				return err
			}
		}

		c := exec.Command(send, title, message)
		return c.Start()
	}

	conn, err := dbus.SessionBus()
	if err != nil {
		return cmd()
	}

	defer conn.Close()

	obj := conn.Object("org.freedesktop.Notifications", dbus.ObjectPath("/org/freedesktop/Notifications"))

	call := obj.Call("org.freedesktop.Notifications", 0, "", uint32(0), "", title, message, []string{}, map[string]dbus.Variant{}, int32(-1))
	if call.Err != nil {
		e := cmd()
		if e != nil {
			return errors.New("beeep: " + call.Err.Error() + "; " + e.Error())
		}
	}

	return nil
}
