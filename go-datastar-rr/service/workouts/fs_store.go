package workouts

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
)

type fileSystemStore struct {
	dir string
}

func NewFSStore(dir string) fileSystemStore {
	return fileSystemStore{
		dir: dir,
	}
}

func (s fileSystemStore) createWorkout(ctx context.Context, w Workout) error {
	return s.persistWorkout(w, os.O_CREATE|os.O_WRONLY)
}

func (s fileSystemStore) saveWorkout(ctx context.Context, w Workout) error {
	return s.persistWorkout(w, os.O_TRUNC|os.O_WRONLY)
}

func (s fileSystemStore) getWorkout(ctx context.Context, id string) (Workout, error) {
	f, err := os.Open(filepath.Join(s.dir, id+".json"))
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
	f, err := os.OpenFile(filepath.Join(s.dir, w.ID+".json"), flags, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(w)
}
