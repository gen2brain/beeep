package beeep

import (
	"testing"
)

func TestNotify(t *testing.T) {
	err := Notify("Notify title", "Message body", "assets/icon128.png")
	if err != nil {
		t.Error(err)
	}
}
