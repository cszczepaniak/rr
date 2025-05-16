package workouts

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/alecthomas/assert/v2"
)

func TestJSON(t *testing.T) {
	workout := Workout{
		ID: "foo",
		Stages: []Stage{
			newReps("category1", "name1", 12),
			newHold("category2", "name2", time.Second),
			newRest(
				newHold("category3", "name3", time.Second),
				time.Minute,
			),
		},
	}

	bs, err := json.Marshal(workout)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(bs))

	var unmarshalledWorkout Workout
	err = json.Unmarshal(bs, &unmarshalledWorkout)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, workout, unmarshalledWorkout)
}
