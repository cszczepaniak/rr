package workouts

import (
	"context"
	"errors"
	"sync"
)

type memoryStore struct {
	mu             sync.Mutex
	inProgressByID map[string]Workout
}

func NewMemoryStore() *memoryStore {
	return &memoryStore{
		inProgressByID: make(map[string]Workout),
	}
}

func (s *memoryStore) createWorkout(ctx context.Context, w Workout) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.inProgressByID[w.ID]; ok {
		return errors.New("workout already existed")
	}

	return s.saveWorkout(ctx, w)
}

func (s *memoryStore) saveWorkout(ctx context.Context, w Workout) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.inProgressByID[w.ID] = w
	return nil
}

func (s *memoryStore) getWorkout(ctx context.Context, id string) (Workout, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	w, ok := s.inProgressByID[id]
	if !ok {
		return Workout{}, errors.New("workout not found")
	}

	return w, nil
}
