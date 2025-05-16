package workouts

import (
	"cmp"
	"encoding/json"
	"fmt"
	"time"

	"github.com/a-h/templ"
)

type timerProps struct {
	Duration  time.Duration
	CountIn   time.Duration
	AutoStart bool
}

type timerSignals struct {
	CountInTicks int       `json:"countInTicks"`
	TimerTicks   int       `json:"timerTicks"`
	CountingIn   bool      `json:"countingIn"`
	Ticks        int       `json:"ticks"`
	Started      bool      `json:"started"`
	AutoStart    bool      `json:"autoStart"`
	Done         bool      `json:"done"`
	Interval     *struct{} `json:"interval"`
}

func (props timerProps) formatSignalsJSON() string {
	signalsJSON := timerSignals{
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

func (props timerProps) formatOnDoneCallback() templ.JSExpression {
	if props.CountIn > 0 {
		return templ.JSExpression(fmt.Sprintf("() => resetTo(%d)", durToTicks(props.Duration)))
	}

	return templ.JSExpression("")
}

func (props timerProps) FormatStartTimerCall() templ.ComponentScript {
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
