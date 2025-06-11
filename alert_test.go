package beeep

import (
	"testing"
)

func TestAlert(t *testing.T) {
	err := Alert("Alert title", "Message body", "testdata/warning.png")
	if err != nil {
		t.Error(err)
	}
}
