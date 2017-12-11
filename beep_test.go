package beeep

import (
	"testing"
)

func TestBeep(t *testing.T) {
	err := Beep(DefaultFreq, DefaultDuration)
	if err != nil {
		t.Error(err)
	}
}
