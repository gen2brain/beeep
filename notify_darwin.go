//go:build darwin && !linux && !freebsd && !netbsd && !openbsd && !windows && !js

package beeep

import (
	"fmt"
	"image/png"
	"os"
	"os/exec"

	"github.com/jackmordaunt/icns/v3"
)

// Notify sends desktop notification.
// The icon can be string with a path to png file or png []byte data. Stock icon names can also be used where supported.
//
// On macOS, this will first try `terminal-notifier` and will fall back to AppleScript with `osascript`.
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
		cmd, err := exec.LookPath("terminal-notifier")
		if err != nil {
			return err
		}

		var img string

		if isBytes {
			tmp1, err := bytesToFilename(icon.([]byte))
			if err != nil {
				return err
			}
			defer os.Remove(tmp1)

			tmp2, err := pngToIcns(tmp1)
			if err != nil {
				return err
			}
			defer os.Remove(tmp2)

			img = tmp2
		} else {
			tmp, err := pngToIcns(pathAbs(icon.(string)))
			if err != nil {
				return err
			}
			defer os.Remove(tmp)

			img = tmp
		}

		var args []string
		if urgent {
			args = []string{"-title", title, "-message", message, "-group", AppName, "-appIcon", img, "-sound", "default"}
		} else {
			args = []string{"-title", title, "-message", message, "-group", AppName, "-appIcon", img}
		}
		c := exec.Command(cmd, args...)

		return c.Run()
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
