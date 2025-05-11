package workouts

import (
	"context"
	"errors"
	"log/slog"
	"sync"
)

type inProgressWorkout struct {
	stageIndex int
	workout    Workout
}

type Service struct {
	mu             sync.Mutex
	inProgressByID map[string]inProgressWorkout
}

func New() *Service {
	return &Service{
		inProgressByID: make(map[string]inProgressWorkout),
	}
}

func (s *Service) Start(ctx context.Context) (Workout, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	workout := newDefaultWorkout()

	s.inProgressByID[workout.ID] = inProgressWorkout{
		workout: workout,
	}
	return s.inProgressByID[workout.ID].workout, nil
}

func (s *Service) Advance(ctx context.Context, id string) (Stage, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	w, ok := s.inProgressByID[id]
	if !ok {
		return nil, errors.New("workout not found")
	}

	slog.Info("stageindex", "n", w.stageIndex)
	w.stageIndex++
	slog.Info("stageindex", "n", w.stageIndex)
	if w.stageIndex >= len(w.workout.Stages) {
		return End{}, nil
	}

	s.inProgressByID[id] = w

	return w.workout.Stages[w.stageIndex], nil
}
