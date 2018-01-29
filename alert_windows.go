// +build windows,!linux,!darwin,!js

package beeep

import (
	"path/filepath"

	"gopkg.in/toast.v1"
)

// Alert displays a desktop notification and plays a default system sound.
func Alert(appid, title, message, appIcon string) error {
	var err error
	iconPath := ""
	if appIcon != "" {
		iconPath, err = filepath.Abs(appIcon)
		if err != nil {
			return err
		}
	}
	note := toastNotification(appid, title, message, iconPath)
	note.Audio = toast.Default
	return note.Push()
}
