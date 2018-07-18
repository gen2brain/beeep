// +build windows,!linux,!freebsd,!netbsd,!openbsd,!darwin,!js

package beeep

import (
	"os/exec"
	"path/filepath"

	"golang.org/x/sys/windows/registry"
	toast "gopkg.in/toast.v1"
)

var isWindows10 = false

func init() {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion`, registry.QUERY_VALUE)
	if err != nil {
		return
	}
	defer k.Close()

	maj, _, err := k.GetIntegerValue("CurrentMajorVersionNumber")
	if err != nil {
		return
	}

	isWindows10 = maj == 10
}

// Notify sends desktop notification.
func Notify(title, message, appIcon string) error {
	if isWindows10 {
		return toastNotify(title, message, appIcon)
	}
	return msgNotify(title, message)
}

func msgNotify(title, message string) error {
	msg, err := exec.LookPath("msg")
	if err != nil {
		return err
	}
	cmd := exec.Command(msg, "*", "/TIME:3", title+"\n\n"+message)
	return cmd.Start()
}

func toastNotify(title, message, appIcon string) error {
	var err error
	iconPath := ""
	if appIcon != "" {
		iconPath, err = filepath.Abs(appIcon)
		if err != nil {
			return err
		}
	}
	notification := toastNotification(title, message, iconPath)
	return notification.Push()
}

func toastNotification(title, message, appIcon string) toast.Notification {
	// NOTE: a real appID is required since Windows 10 Fall Creator's Update,
	// issue https://github.com/go-toast/toast/issues/9
	appID := "{1AC14E77-02E7-4E5D-B744-2EB1AE5198B7}\\WindowsPowerShell\\v1.0\\powershell.exe"
	return toast.Notification{
		AppID:   appID,
		Title:   title,
		Message: message,
		Icon:    appIcon,
	}
}
