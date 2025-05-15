package components

import (
	"cmp"
	"encoding/json"
	"fmt"
	"time"

	"github.com/a-h/templ"
)

type TimerProps struct {
	Duration  time.Duration
	CountIn   time.Duration
	AutoStart bool
}

type TimerSignals struct {
	CountInTicks int       `json:"countInTicks"`
	TimerTicks   int       `json:"timerTicks"`
	CountingIn   bool      `json:"countingIn"`
	Ticks        int       `json:"ticks"`
	Started      bool      `json:"started"`
	AutoStart    bool      `json:"autoStart"`
	Done         bool      `json:"done"`
	Interval     *struct{} `json:"interval"`
}

func (props TimerProps) FormatSignalsJSON() string {
	signalsJSON := TimerSignals{
		CountInTicks: durToTicks(props.CountIn),
		TimerTicks:   durToTicks(props.Duration),
		CountingIn:   props.CountIn > 0,
		Ticks:        durToTicks(cmp.Or(props.CountIn, props.Duration)),
		Started:      false,
		Done:         false,
		AutoStart:    props.AutoStart,
		Interval:     nil,
	}

	bs, _ := json.Marshal(signalsJSON)
	return string(bs)
}

func (props TimerProps) formatOnDoneCallback() templ.JSExpression {
	if props.CountIn > 0 {
		return templ.JSExpression(fmt.Sprintf("() => resetTo(%d)", durToTicks(props.Duration)))
	}

	return templ.JSExpression("")
}

func (props TimerProps) FormatStartTimerCall() templ.ComponentScript {
	if props.CountIn > 0 {
		return templ.JSFuncCall(
			"startTimer",
			props.formatOnDoneCallback(),
		)
	}

	return templ.JSFuncCall("startTimer")
}

const tick = 100 * time.Millisecond

func durToTicks(dur time.Duration) int {
	return int(dur / tick)
}
