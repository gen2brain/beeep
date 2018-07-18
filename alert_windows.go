// +build windows,!linux,!freebsd,!netbsd,!openbsd,!darwin,!js

package beeep

import (
	"path/filepath"

	toast "gopkg.in/toast.v1"
)

// Alert displays a desktop notification and plays a default system sound.
func Alert(title, message, appIcon string) error {
	var err error
	iconPath := ""
	if appIcon != "" {
		iconPath, err = filepath.Abs(appIcon)
		if err != nil {
			return err
		}
	}
	note := toastNotification(title, message, iconPath)
	note.Audio = toast.Default
	return note.Push()
}
