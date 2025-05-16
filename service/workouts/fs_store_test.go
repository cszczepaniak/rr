package workouts

import (
	"errors"
	"io/fs"
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestFSStore(t *testing.T) {
	dir := t.TempDir()

	s, err := NewFSStore(dir)
	assert.NoError(t, err)

	w := newDefaultWorkout()

	assert.NoError(t, s.createWorkout(t.Context(), w))

	gotW, err := s.getWorkout(t.Context(), w.ID)
	assert.NoError(t, err)

	assert.Equal(t, w, gotW)

	w.CurrentStage += 10
	assert.NoError(t, s.createWorkout(t.Context(), w))

	gotW, err = s.getWorkout(t.Context(), w.ID)
	assert.NoError(t, err)

	assert.Equal(t, w, gotW)
}

func TestFSStore_PathEscape(t *testing.T) {
	dir := t.TempDir()

	s, err := NewFSStore(dir)
	assert.NoError(t, err)

	_, err = s.getWorkout(t.Context(), "../bad_actor")
	var expErr *fs.PathError
	assert.True(t, errors.As(err, &expErr))
}
