package workouts

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFSStore(t *testing.T) {
	dir := t.TempDir()

	s := NewFSStore(dir)

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
