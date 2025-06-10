//go:build js

package beeep

// Alert displays a desktop notification and plays a beep.
func Alert(title, message, icon string) error {
	if err := Notify(title, message, icon); err != nil {
		return err
	}

	return Beep(DefaultFreq, DefaultDuration)
}
