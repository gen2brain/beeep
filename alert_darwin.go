//go:build darwin && !linux && !freebsd && !netbsd && !openbsd && !windows && !js

package beeep

// Alert displays a desktop notification and plays a default system sound.
func Alert(title, message string, icon any) error {
	return notify1(title, message, icon, true)
}
