package workouts

import (
	"cmp"
	"encoding/json"
	"time"
)

type timerProps struct {
	autoAdvance bool
	workoutID   string

	duration  time.Duration
	countIn   time.Duration
	autoStart bool
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
		CountInTicks: durToTicks(props.countIn),
		TimerTicks:   durToTicks(props.duration),
		CountingIn:   props.countIn > 0,
		Ticks:        durToTicks(cmp.Or(props.countIn, props.duration)),
		Started:      false,
		Done:         false,
		AutoStart:    props.autoStart,
		Interval:     nil,
	}

	bs, _ := json.Marshal(signalsJSON)
	return string(bs)
}

const tick = 100 * time.Millisecond

func durToTicks(dur time.Duration) int {
	return int(dur / tick)
}
