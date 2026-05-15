//go:build darwin && !linux && !freebsd && !netbsd && !openbsd && !windows && !js

package beeep

/*
#include <libproc.h>
*/
import "C"

import (
	"fmt"
	"image/png"
	"os"
	"os/exec"
	"time"
	"unsafe"

	"github.com/jackmordaunt/icns/v3"
)

// Notify sends desktop notification.
// The icon can be string with a path to png file or png []byte data. Stock icon names can also be used where supported.
//
// On macOS, this will first try `alerter` and will fall back to AppleScript with `osascript`.
func Notify(title, message string, icon any) error {
	return notify1(title, message, icon, false)
}

func notify1(title, message string, icon any, urgent bool) error {
	var isBytes bool
	switch icon.(type) {
	case string:
	case []byte:
		isBytes = true
	default:
		return fmt.Errorf("unsupported argument: %T", icon)
	}

	cmd1 := func() error {
		cmd, err := exec.LookPath("alerter")
		if err != nil {
			return err
		}

		var tmpFiles []string
		cleanup := func() {
			for _, f := range tmpFiles {
				os.Remove(f)
			}
		}

		var img string

		if isBytes {
			tmp1, err := bytesToFilename(icon.([]byte))
			if err != nil {
				return err
			}
			tmpFiles = append(tmpFiles, tmp1)

			tmp2, err := pngToIcns(tmp1)
			if err != nil {
				cleanup()
				return err
			}
			tmpFiles = append(tmpFiles, tmp2)

			img = tmp2
		} else {
			tmp, err := pngToIcns(pathAbs(icon.(string)))
			if err != nil {
				return err
			}
			tmpFiles = append(tmpFiles, tmp)
			img = tmp
		}

		var args []string
		if urgent {
			args = []string{"--title", title, "--message", message, "--group", AppName, "--app-icon", img, "--sound", "default"}
		} else {
			args = []string{"--title", title, "--message", message, "--group", AppName, "--app-icon", img}
		}

		c := exec.Command(cmd, args...)

		if err := c.Start(); err != nil {
			cleanup()
			return err
		}

		// Block until alerter has finished setup and entered its run loop,
		// ensuring the icon temp file is not removed before alerter reads it.
		waitUntilIdle(c.Process.Pid, 5*time.Second)

		// Schedule cleanup for when alerter eventually exits.
		go func() {
			c.Wait()
			cleanup()
		}()

		return nil
	}

	cmd2 := func() error {
		osa, err := exec.LookPath("osascript")
		if err != nil {
			return err
		}

		var script string
		if urgent {
			script = fmt.Sprintf("display notification %q with title %q sound name \"default\"", message, title)
		} else {
			script = fmt.Sprintf("display notification %q with title %q", message, title)
		}

		cmd := exec.Command(osa, "-e", script)

		return cmd.Run()
	}

	err1 := cmd1()
	if err1 != nil {
		err2 := cmd2()
		if err2 != nil {
			return fmt.Errorf("beeep: terminal-notifier: %w; osascript: %w", err1, err2)
		}
	}

	return nil
}

func pngToIcns(icon string) (string, error) {
	var out string

	f, err := os.Open(icon)
	if err != nil {
		return out, err
	}
	defer f.Close()

	img, err := png.Decode(f)
	if err != nil {
		return out, err
	}

	tmp, err := os.CreateTemp(os.TempDir(), "beeep*.icns")
	if err != nil {
		return out, err
	}
	defer tmp.Close()

	out = tmp.Name()

	err = icns.Encode(tmp, img)
	if err != nil {
		return out, err
	}

	return out, nil
}

func pidTaskInfo(pid int) (uint64, error) {
	var info C.struct_proc_taskinfo

	ret := C.proc_pidinfo(C.int(pid), C.PROC_PIDTASKINFO, 0, unsafe.Pointer(&info), C.int(unsafe.Sizeof(info)))
	if ret <= 0 {
		return 0, fmt.Errorf("proc_pidinfo failed")
	}

	return uint64(info.pti_total_user) + uint64(info.pti_total_system), nil
}

func waitUntilIdle(pid int, timeout time.Duration) {
	deadline := time.Now().Add(timeout)
	var prev uint64
	seenActivity := false
	stable := 0

	for time.Now().Before(deadline) {
		time.Sleep(10 * time.Millisecond)

		cur, err := pidTaskInfo(pid)
		if err != nil {
			return
		}

		if cur > 0 {
			seenActivity = true
		}

		if seenActivity {
			if cur == prev {
				stable++
				if stable >= 3 {
					return
				}
			} else {
				stable = 0
			}
		}

		prev = cur
	}
}
