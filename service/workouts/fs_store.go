package workouts

import (
	"context"
	"encoding/json"
	"os"
)

type fileSystemStore struct {
	root *os.Root
}

func NewFSStore(dir string) (fileSystemStore, error) {
	r, err := os.OpenRoot(dir)
	if err != nil {
		return fileSystemStore{}, err
	}

	return fileSystemStore{
		root: r,
	}, nil
}

func (s fileSystemStore) createWorkout(ctx context.Context, w Workout) error {
	return s.persistWorkout(w, os.O_CREATE|os.O_WRONLY)
}

func (s fileSystemStore) saveWorkout(ctx context.Context, w Workout) error {
	return s.persistWorkout(w, os.O_TRUNC|os.O_WRONLY)
}

func (s fileSystemStore) getWorkout(ctx context.Context, id string) (Workout, error) {
	f, err := s.root.Open(id + ".json")
	if err != nil {
		return Workout{}, err
	}
	defer f.Close()

	var workout Workout
	err = json.NewDecoder(f).Decode(&workout)
	if err != nil {
		return Workout{}, err
	}

	return workout, nil
}

func (s fileSystemStore) persistWorkout(w Workout, flags int) error {
	f, err := s.root.OpenFile(w.ID+".json", flags, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(w)
}
