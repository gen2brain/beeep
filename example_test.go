package beeep_test

import (
	_ "embed"

	"github.com/gen2brain/beeep"
)

//go:embed testdata/info.png
var icon []byte

func ExampleNotify() {
	_ = beeep.Notify("Title", "MessageBody", icon) // icon is embedded file
}

func ExampleAlert() {
	beeep.AppName = "My App Name" // change the default app name

	_ = beeep.Alert("Title", "MessageBody", "testdata/warning.png")
}

func ExampleBeep() {
	_ = beeep.Beep(beeep.DefaultFreq, beeep.DefaultDuration)
}
