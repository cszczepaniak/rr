package workouts

import (
	"context"
	"log/slog"
)

type Store interface {
	createWorkout(ctx context.Context, w Workout) error
	saveWorkout(ctx context.Context, w Workout) error
	getWorkout(ctx context.Context, id string) (Workout, error)
}

type Service struct {
	store Store
}

func New(s Store) *Service {
	return &Service{
		store: s,
	}
}

func (s *Service) Start(ctx context.Context) (string, error) {
	workout := newDefaultWorkout()

	err := s.store.createWorkout(ctx, workout)
	if err != nil {
		return "", err
	}

	return workout.ID, nil
}

func (s *Service) GetCurrentStage(ctx context.Context, id string) (Stage, error) {
	w, err := s.store.getWorkout(ctx, id)
	if err != nil {
		return nil, err
	}

	if w.IsComplete() {
		return End{}, nil
	}

	return w.Stages[w.CurrentStage], nil
}

func (s *Service) Advance(ctx context.Context, id string) (Stage, error) {
	w, err := s.store.getWorkout(ctx, id)
	if err != nil {
		return nil, err
	}

	w = w.Advance()
	err = s.store.saveWorkout(ctx, w)
	if err != nil {
		return nil, err
	}

	if w.IsComplete() {
		slog.Info("advance", "status", "ending workout")
		return End{}, nil
	}

	return w.Stages[w.CurrentStage], nil
}
