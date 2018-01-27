// +build !linux,!windows,!darwin,!js

package beeep

// Notify sends desktop notification.
func Notify(title, message string) error {
	return ErrUnsupported
}
