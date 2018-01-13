// +build windows,!linux,!darwin,!js

package beeep

import (
	"bytes"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	toast "gopkg.in/toast.v1"
)

// Notify sends desktop notification.
func Notify(title, message, appIcon string) error {
	if isWindows10() {
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

func isWindows10() bool {
	ver := getWindowsVersionString()
	parts := strings.Split(ver, ".")
	if len(parts) < 1 {
		return false
	}
	i, err := strconv.Atoi(parts[0])
	if err != nil {
		return false
	}
	return i == 10
}

// Returns the Windows version string, such as "10.0.16299.125" on Windows 10
func getWindowsVersionString() string {
	cmd := exec.Command("cmd", "ver")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
	s := strings.ToLower(strings.Replace(out.String(), "\r\n", "", -1))
	p1 := strings.Index(s, "[version")
	p2 := strings.Index(s, "]")
	var ver string
	if p1 == -1 || p2 == -1 {
		ver = "unknown"
	} else {
		ver = s[p1+9 : p2]
	}
	return ver
}
