//go:build windows && !linux && !freebsd && !netbsd && !openbsd && !darwin && !js
// +build windows,!linux,!freebsd,!netbsd,!openbsd,!darwin,!js

package beeep

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"syscall"
	"time"

	toast "github.com/go-toast/toast"
	"github.com/tadvi/systray"
	"golang.org/x/sys/windows/registry"
)

var isWindows10 bool
var applicationID string

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

	if isWindows10 {
		applicationID = appID()
	}
}

// Notify sends desktop notification with default timeout
func Notify(title, message, appIcon string) error {
	return NotifyEx(title, message, appIcon, 3*time.Second)
}

// NotifyEx sends desktop notification with timeout
func NotifyEx(title, message, appIcon string, timeout time.Duration) error {
	if isWindows10 {
		return toastNotify(title, message, appIcon)
	}

	err := baloonNotify(title, message, appIcon, false)
	if err != nil {
		e := msgNotify(title, message, timeout)
		if e != nil {
			return errors.New("beeep: " + err.Error() + "; " + e.Error())
		}
	}

	return nil

}

func msgNotify(title, message string, timeout time.Duration) error {
	msg, err := exec.LookPath("msg")
	if err != nil {
		return err
	}
	if timeout <= 0 {
		timeout = time.Second * 3
	}

	cmd := exec.Command(msg, "*", fmt.Sprintf("/TIME:%d", timeout/time.Second), title+"\n\n"+message)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd.Run()
}

func baloonNotify(title, message, appIcon string, bigIcon bool) error {
	tray, err := systray.New()
	if err != nil {
		return err
	}

	err = tray.ShowCustom(pathAbs(appIcon), title)
	if err != nil {
		return err
	}

	go func() {
		go func() {
			_ = tray.Run()
		}()
		time.Sleep(3 * time.Second)
		_ = tray.Stop()
	}()

	return tray.ShowMessage(title, message, bigIcon)
}

func toastNotify(title, message, appIcon string) error {
	notification := toastNotification(title, message, pathAbs(appIcon))
	return notification.Push()
}

func toastNotification(title, message, appIcon string) toast.Notification {
	return toast.Notification{
		AppID:   applicationID,
		Title:   title,
		Message: message,
		Icon:    appIcon,
	}
}

func appID() string {
	defID := "{1AC14E77-02E7-4E5D-B744-2EB1AE5198B7}\\WindowsPowerShell\\v1.0\\powershell.exe"
	cmd := exec.Command("powershell", "-NoProfile", "Get-StartApps")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.Output()
	if err != nil {
		return defID
	}

	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.Contains(line, "powershell.exe") {
			sp := strings.Split(line, " ")
			if len(sp) > 0 {
				return sp[len(sp)-1]
			}
		}
	}

	return defID
}
