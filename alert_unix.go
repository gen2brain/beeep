//go:build linux || freebsd || netbsd || openbsd || illumos

package beeep

// Alert displays a desktop notification and plays a beep.
func Alert(title, message, appIcon string) error {
	if err := notify1(title, message, appIcon, true); err != nil {
		return err
	}

	return Beep(DefaultFreq, DefaultDuration)
}
