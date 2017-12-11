package beeep

import (
	"testing"
)

func TestNotify(t *testing.T) {
	err := Notify("Title", "Message body")
	if err != nil {
		t.Error(err)
	}
}
