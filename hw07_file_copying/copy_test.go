package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const tempFile = "/tmp/out"

func requireEqualToTmp(t *testing.T, expectedPath string) {
	t.Helper()

	expected, err := ioutil.ReadFile(expectedPath)
	require.Nil(t, err)

	actual, err := ioutil.ReadFile(tempFile)
	require.Nil(t, err)

	require.Equal(t, expected, actual)
}

func TestCopy(t *testing.T) {
	t.Run("nonexistent source file", func(t *testing.T) {
		err := Copy("nonexistent", tempFile, 0, 0)
		require.ErrorIs(t, err, os.ErrNotExist)
	})

	t.Run("offset can not be larger than source size", func(t *testing.T) {
		err := Copy("testdata/input.txt", tempFile, 1000000, 0)
		require.ErrorIs(t, err, ErrOffsetExceedsFileSize)
	})

	t.Run("files with unknown size are unsupported", func(t *testing.T) {
		err := Copy("/dev/urandom", tempFile, 0, 0)
		require.Equal(t, err, ErrUnsupportedFile)
	})

	t.Run("whole file copy", func(t *testing.T) {
		err := Copy("testdata/input.txt", tempFile, 0, 0)

		require.Nil(t, err)
		requireEqualToTmp(t, "testdata/out_offset0_limit0.txt")
	})

	t.Run("file copy (offset 0, limit 10)", func(t *testing.T) {
		err := Copy("testdata/input.txt", tempFile, 0, 10)

		require.Nil(t, err)
		requireEqualToTmp(t, "testdata/out_offset0_limit10.txt")
	})

	t.Run("file copy (offset 0, limit 1000)", func(t *testing.T) {
		err := Copy("testdata/input.txt", tempFile, 0, 1000)

		require.Nil(t, err)
		requireEqualToTmp(t, "testdata/out_offset0_limit1000.txt")
	})

	t.Run("file copy (offset 0, limit 10000)", func(t *testing.T) {
		err := Copy("testdata/input.txt", tempFile, 0, 10000)

		require.Nil(t, err)
		requireEqualToTmp(t, "testdata/out_offset0_limit10000.txt")
	})

	t.Run("file copy (offset 100, limit 1000)", func(t *testing.T) {
		err := Copy("testdata/input.txt", tempFile, 100, 1000)

		require.Nil(t, err)
		requireEqualToTmp(t, "testdata/out_offset100_limit1000.txt")
	})

	t.Run("file copy (offset 6000, limit 1000)", func(t *testing.T) {
		err := Copy("testdata/input.txt", tempFile, 6000, 1000)

		require.Nil(t, err)
		requireEqualToTmp(t, "testdata/out_offset6000_limit1000.txt")
	})
}
