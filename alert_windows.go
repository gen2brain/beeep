// +build windows,!linux,!freebsd,!netbsd,!openbsd,!darwin,!js

package beeep

import (
	"path/filepath"

	toast "gopkg.in/toast.v1"
)

// Alert displays a desktop notification and plays a default system sound.
func Alert(title, message, appIcon string) error {
	note := toastNotification(title, message, pathAbs(appIcon))
	note.Audio = toast.Default
	return note.Push()
}
