package workouts

import (
	"iter"
	"slices"
	"time"

	"github.com/google/uuid"
)

type Stage interface {
	Category() string
	Name() string
}

type stageBase struct {
	category string
	name     string
}

func (b stageBase) Category() string {
	return b.category
}

func (b stageBase) Name() string {
	return b.name
}

type Reps struct {
	stageBase
	Reps int `json:"reps"`
}

func newReps(c, n string, reps int) Reps {
	return Reps{
		stageBase: stageBase{
			category: c,
			name:     n,
		},
		Reps: reps,
	}
}

type Hold struct {
	stageBase
	Duration time.Duration `json:"dur"`
}

func newHold(c, n string, dur time.Duration) Hold {
	return Hold{
		stageBase: stageBase{
			category: c,
			name:     n,
		},
		Duration: dur,
	}
}

type Rest struct {
	stageBase
	Duration time.Duration `json:"dur"`
}

func newRest(from Stage, dur time.Duration) Rest {
	return Rest{
		stageBase: stageBase{
			category: from.Category(),
			name:     from.Name(),
		},
		Duration: dur,
	}
}

func newWithRest(s Stage, restDur time.Duration) iter.Seq[Stage] {
	return func(yield func(Stage) bool) {
		if !yield(s) || !yield(newRest(s, restDur)) {
			return
		}
	}
}

type End struct {
	stageBase
}

type Workout struct {
	ID           string  `json:"id"`
	Stages       []Stage `json:"stages"`
	CurrentStage int     `json:"current_stage"`
}

func newDefaultWorkout() Workout {
	return Workout{
		ID: uuid.NewString(),
		Stages: slices.Collect(concat(
			iterOf(newReps("Warmup", "Shoulder Rolls", 10)),
			iterOf(newReps("Warmup", "Scapular Shrugs", 10)),
			iterOf(newReps("Warmup", "Cat/Camel", 10)),
			iterOf(newReps("Warmup", "Band Pulldowns", 10)),
			iterOf(newReps("Warmup", "Band Dislocates", 10)),
			iterOf(newReps("Warmup", "Wrist Mobility", 1)),
			iterOf(newHold("Warmup", "Hamstring Stretch (1 of 2)", 30*time.Second)),
			iterOf(newHold("Warmup", "Hamstring Stretch (2 of 2)", 30*time.Second)),
			repeat(3,
				newWithRest(
					newHold("Skill Work", "Parallel Bar Support", 30*time.Second),
					90*time.Second,
				),
				newWithRest(
					newHold("Skill Work", "Handstand (Wall Start)", 30*time.Second),
					90*time.Second,
				),
			),
			repeat(3,
				newWithRest(
					newReps("Strength Work (Set 1)", "RTO Pushup", 8),
					90*time.Second,
				),
				newWithRest(
					newReps("Strength Work (Set 1)", "Tuck Ice Cream Makers", 8),
					90*time.Second,
				),
			),
			repeat(3,
				newWithRest(
					newHold("Strength Work (Set 2)", "L-Sit (Foot Supported)", 30*time.Second),
					90*time.Second,
				),
				newWithRest(
					newReps("Strength Work (Set 2)", "One-foot Step Ups", 8),
					90*time.Second,
				),
			),
			repeat(3,
				newWithRest(
					newReps("Strength Work (Set 3)", "Pullups", 8),
					90*time.Second,
				),
				newWithRest(
					newReps("Strength Work (Set 3)", "Ring Dips", 8),
					90*time.Second,
				),
			),
			repeat(3,
				newWithRest(
					newHold("Bodyline Drills", "Plank (Elbows)", 30*time.Second),
					60*time.Second,
				),
				newWithRest(
					newHold("Bodyline Drills", "Side Plank (Left Elbow)", 30*time.Second),
					60*time.Second,
				),
				newWithRest(
					newHold("Bodyline Drills", "Side Plank (Right Elbow)", 30*time.Second),
					60*time.Second,
				),
				newWithRest(
					newHold("Bodyline Drills", "Hollow Hold", 30*time.Second),
					60*time.Second,
				),
				newWithRest(
					newHold("Bodyline Drills", "Superman", 30*time.Second),
					60*time.Second,
				),
			),
		)),
	}
}

func concat(seqs ...iter.Seq[Stage]) iter.Seq[Stage] {
	return func(yield func(Stage) bool) {
		for _, seq := range seqs {
			for s := range seq {
				if !yield(s) {
					return
				}
			}
		}
	}
}

func iterOf(stages ...Stage) iter.Seq[Stage] {
	return slices.Values(stages)
}

func repeat(n int, seqs ...iter.Seq[Stage]) iter.Seq[Stage] {
	flattened := concat(seqs...)
	return func(yield func(Stage) bool) {
		for range n {
			for s := range flattened {
				if !yield(s) {
					return
				}
			}
		}
	}
}
