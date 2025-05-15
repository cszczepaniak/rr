package workouts

import (
	"io/fs"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFSStore(t *testing.T) {
	dir := t.TempDir()

	s, err := NewFSStore(dir)
	require.NoError(t, err)

	w := newDefaultWorkout()

	require.NoError(t, s.createWorkout(t.Context(), w))

	gotW, err := s.getWorkout(t.Context(), w.ID)
	require.NoError(t, err)

	assert.Equal(t, w, gotW)

	w.CurrentStage += 10
	require.NoError(t, s.createWorkout(t.Context(), w))

	gotW, err = s.getWorkout(t.Context(), w.ID)
	require.NoError(t, err)

	assert.Equal(t, w, gotW)
}

func TestFSStore_PathEscape(t *testing.T) {
	dir := t.TempDir()

	s, err := NewFSStore(dir)
	require.NoError(t, err)

	_, err = s.getWorkout(t.Context(), "../bad_actor")
	var expErr *fs.PathError
	require.ErrorAs(t, err, &expErr)
}
