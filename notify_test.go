package beeep

import (
	_ "embed"
	"testing"
)

//go:embed testdata/info.png
var data []byte

func TestNotify(t *testing.T) {
	err := Notify("Notify title", "Message body", data)
	if err != nil {
		t.Error(err)
	}
}
