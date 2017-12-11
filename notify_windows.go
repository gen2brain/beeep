// +build windows,!linux,!darwin,!js

package beeep

import (
	"syscall"
	"unsafe"
)

const (
	mbOk                  = 0x00000000
	mbServiceNotification = 0x00200000
)

var (
	user32     = syscall.NewLazyDLL("user32.dll")
	messageBox = user32.NewProc("MessageBoxW")
)

// Notify sends desktop notification.
func Notify(title, message string) error {
	messageBox.Call(0,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(title))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(message))),
		uintptr(mbOk|mbServiceNotification))

	return nil
}
