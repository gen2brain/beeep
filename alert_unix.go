//go:build linux || freebsd || netbsd || openbsd || illumos

package beeep

// Alert displays a desktop notification and plays a beep.
func Alert(title, message string, icon any) error {
	if err := notify1(title, message, icon, true); err != nil {
		return err
	}

	err := Beep(DefaultFreq, DefaultDuration)
	if err != nil {
		return err
	}

	return nil
}
