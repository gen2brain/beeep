// +build windows,!linux,!darwin,!js

package beeep

import (
	"os/exec"
)

// Notify sends desktop notification.
func Notify(title, message string) error {
	msg, err := exec.LookPath("msg")
	if err != nil {
		return err
	}

	cmd := exec.Command(msg, "*", "/TIME:3", title+"\n\n"+message)
	return cmd.Start()
}
