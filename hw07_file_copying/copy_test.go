package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	// Place your code here.
	t.Run("/dev/urandom is unsupported", func(t *testing.T) {
		err := Copy("/dev/urandom", "/tmp/urandom", 0, 0)
		require.Equal(t, err, ErrUnsupportedFile)
	})

	t.Run("nonexistent source file", func(t *testing.T) {
		err := Copy("nonexistent", "/tmp/f", 0, 0)
		require.ErrorIs(t, err, os.ErrNotExist)
	})
}
