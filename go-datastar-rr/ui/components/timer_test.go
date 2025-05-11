package components

import (
	"testing"
	"time"
)

func TestTimerProps(t *testing.T) {
	props := TimerProps{
		Duration:  time.Second,
		CountIn:   2 * time.Second,
		AutoStart: false,
	}

	t.Log(props.FormatSignalsJSON())
}
