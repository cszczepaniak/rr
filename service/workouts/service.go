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

	return hackRestOffset(w), nil
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

	return hackRestOffset(w), nil
}

func hackRestOffset(w Workout) Stage {
	stage := w.Stages[w.CurrentStage]
	rest, ok := stage.(Rest)
	if !ok {
		return stage
	}

	// HACK: if this is Rest, we're going to swap in the _next_ stage so we can render what's
	// coming next during the rest.
	nextIdx := w.CurrentStage + 1
	if nextIdx >= len(w.Stages) {
		// This would be surprising! Let's log but continue with the orignal Rest.
		slog.Error("advance", "status", "rest with no next stage")
		return stage
	}

	nextStage := w.Stages[nextIdx]
	return newRest(nextStage, rest.Duration)
}
