package workouts

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"
)

func marshalStage(s Stage) ([]byte, error) {
	switch s := s.(type) {
	case Hold:
		return json.Marshal(struct {
			Category string        `json:"category"`
			Name     string        `json:"name"`
			Duration time.Duration `json:"duration"`
		}{
			Category: s.category,
			Name:     s.name,
			Duration: s.Duration,
		})
	case Reps:
		return json.Marshal(struct {
			Category string `json:"category"`
			Name     string `json:"name"`
			Reps     int    `json:"reps"`
		}{
			Category: s.category,
			Name:     s.name,
			Reps:     s.Reps,
		})
	case Rest:
		return json.Marshal(struct {
			Category string        `json:"category"`
			Name     string        `json:"name"`
			Duration time.Duration `json:"duration"`
		}{
			Category: s.category,
			Name:     s.name,
			Duration: s.Duration,
		})
	}

	return nil, fmt.Errorf("marshalling not supported for %T", s)
}

type marshalledStage struct {
	Type  string          `json:"type"`
	Value json.RawMessage `json:"value"`
}

func (w Workout) MarshalJSON() ([]byte, error) {
	marshalable := struct {
		ID     string            `json:"id"`
		Stages []marshalledStage `json:"stages"`
	}{
		ID: w.ID,
	}
	for _, s := range w.Stages {
		marshalled, err := marshalStage(s)
		if err != nil {
			return nil, err
		}

		marshalable.Stages = append(marshalable.Stages, marshalledStage{
			Type:  reflect.TypeOf(s).String(),
			Value: marshalled,
		})
	}

	return json.Marshal(marshalable)
}

func (w *Workout) UnmarshalJSON(bs []byte) error {
	var data struct {
		ID     string            `json:"id"`
		Stages []json.RawMessage `json:"stages"`
	}
	err := json.Unmarshal(bs, &data)
	if err != nil {
		return err
	}

	w.ID = data.ID
	for _, mStage := range data.Stages {
		type unmarshalledStage struct {
			Type  string `json:"type"`
			Value struct {
				Category string        `json:"category"`
				Name     string        `json:"name"`
				Reps     int           `json:"reps,omitzero"`
				Duration time.Duration `json:"duration,omitzero"`
			} `json:"value"`
		}

		var uStage unmarshalledStage
		err := json.Unmarshal(mStage, &uStage)
		if err != nil {
			return err
		}

		switch uStage.Type {
		case reflect.TypeOf(Reps{}).String():
			w.Stages = append(w.Stages, newReps(uStage.Value.Category, uStage.Value.Name, uStage.Value.Reps))
		case reflect.TypeOf(Hold{}).String():
			w.Stages = append(w.Stages, newHold(uStage.Value.Category, uStage.Value.Name, uStage.Value.Duration))
		case reflect.TypeOf(Rest{}).String():
			w.Stages = append(w.Stages, Rest{
				stageBase: stageBase{
					category: uStage.Value.Category,
					name:     uStage.Value.Name,
				},
				Duration: uStage.Value.Duration,
			})
		}
	}

	return nil
}
