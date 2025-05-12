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

func (s Stages) MarshalJSON() ([]byte, error) {
	var marshalable []marshalledStage
	for _, stage := range s {
		marshalled, err := marshalStage(stage)
		if err != nil {
			return nil, err
		}

		marshalable = append(marshalable, marshalledStage{
			Type:  reflect.TypeOf(stage).String(),
			Value: marshalled,
		})
	}

	return json.Marshal(marshalable)
}

func (s *Stages) UnmarshalJSON(bs []byte) error {
	var data []marshalledStage
	err := json.Unmarshal(bs, &data)
	if err != nil {
		return err
	}

	stages := make(Stages, 0, len(data))
	for _, mStage := range data {
		type unmarshalledStage struct {
			Category string        `json:"category"`
			Name     string        `json:"name"`
			Reps     int           `json:"reps,omitzero"`
			Duration time.Duration `json:"duration,omitzero"`
		}

		var uStage unmarshalledStage
		err := json.Unmarshal(mStage.Value, &uStage)
		if err != nil {
			return err
		}

		switch mStage.Type {
		case reflect.TypeOf(Reps{}).String():
			stages = append(stages, newReps(uStage.Category, uStage.Name, uStage.Reps))
		case reflect.TypeOf(Hold{}).String():
			stages = append(stages, newHold(uStage.Category, uStage.Name, uStage.Duration))
		case reflect.TypeOf(Rest{}).String():
			stages = append(stages, Rest{
				stageBase: stageBase{
					category: uStage.Category,
					name:     uStage.Name,
				},
				Duration: uStage.Duration,
			})
		}
	}

	*s = stages

	return nil
}
