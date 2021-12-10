package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func requireFilesEqual(t *testing.T, expectedPath, actualPath string) {
	expected, err := ioutil.ReadFile(expectedPath)
	require.Nil(t, err)

	actual, err := ioutil.ReadFile(actualPath)
	require.Nil(t, err)

	require.Equal(t, expected, actual)
}

func TestCopy(t *testing.T) {
	t.Run("nonexistent source file", func(t *testing.T) {
		err := Copy("nonexistent", "/tmp/out", 0, 0)
		require.ErrorIs(t, err, os.ErrNotExist)
	})

	t.Run("offset can not be larger than source size", func(t *testing.T) {
		err := Copy("testdata/input.txt", "/tmp/out", 1000000, 0)
		require.ErrorIs(t, err, ErrOffsetExceedsFileSize)
	})

	t.Run("files with unknown size are unsupported", func(t *testing.T) {
		err := Copy("/dev/urandom", "/tmp/out", 0, 0)
		require.Equal(t, err, ErrUnsupportedFile)
	})

	t.Run("whole file copy", func(t *testing.T) {
		err := Copy("testdata/input.txt", "/tmp/out", 0, 0)

		require.Nil(t, err)
		requireFilesEqual(t, "testdata/out_offset0_limit0.txt", "/tmp/out")
	})

	t.Run("file copy (offset 0, limit 10)", func(t *testing.T) {
		err := Copy("testdata/input.txt", "/tmp/out", 0, 10)

		require.Nil(t, err)
		requireFilesEqual(t, "testdata/out_offset0_limit10.txt", "/tmp/out")
	})

	t.Run("file copy (offset 0, limit 1000)", func(t *testing.T) {
		err := Copy("testdata/input.txt", "/tmp/out", 0, 1000)

		require.Nil(t, err)
		requireFilesEqual(t, "testdata/out_offset0_limit1000.txt", "/tmp/out")
	})

	t.Run("file copy (offset 0, limit 10000)", func(t *testing.T) {
		err := Copy("testdata/input.txt", "/tmp/out", 0, 10000)

		require.Nil(t, err)
		requireFilesEqual(t, "testdata/out_offset0_limit10000.txt", "/tmp/out")
	})

	t.Run("file copy (offset 100, limit 1000)", func(t *testing.T) {
		err := Copy("testdata/input.txt", "/tmp/out", 100, 1000)

		require.Nil(t, err)
		requireFilesEqual(t, "testdata/out_offset100_limit1000.txt", "/tmp/out")
	})

	t.Run("file copy (offset 6000, limit 1000)", func(t *testing.T) {
		err := Copy("testdata/input.txt", "/tmp/out", 6000, 1000)

		require.Nil(t, err)
		requireFilesEqual(t, "testdata/out_offset6000_limit1000.txt", "/tmp/out")
	})
}
