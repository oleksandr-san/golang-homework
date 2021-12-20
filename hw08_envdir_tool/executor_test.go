package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func captureStdout(t *testing.T, f func()) string {
	t.Helper()

	rescueStdout := os.Stdout
	r, w, err := os.Pipe()
	require.Nil(t, err)
	os.Stdout = w

	f()

	w.Close()
	out, err := ioutil.ReadAll(r)
	require.Nil(t, err)
	os.Stdout = rescueStdout

	return string(out)
}

func TestRunCmd(t *testing.T) {
	t.Run("successful run", func(t *testing.T) {
		cmd := []string{"/bin/bash", "testdata/echo.sh", "arg1=1", "arg2=2"}
		env := Environment{
			"HELLO": EnvValue{"\"hello\"", false},
			"BAR":   EnvValue{"bar", false},
			"FOO":   EnvValue{"   foo\nnew line", false},
			"EMPTY": EnvValue{"", false},
			"UNSET": EnvValue{"", true},
		}
		expectedCode := 0
		expectedStdout := `HELLO is ("hello")
BAR is (bar)
FOO is (   foo
new line)
UNSET is ()
ADDED is ()
EMPTY is ()
arguments are arg1=1 arg2=2

`

		returnCode := -1
		stdout := captureStdout(t, func() {
			returnCode = RunCmd(cmd, env)
		})

		require.Equal(t, expectedCode, returnCode)
		require.Equal(t, expectedStdout, stdout)
	})

	t.Run("failed run", func(t *testing.T) {
		cmd := []string{"/bin/bash", "testdata/env"}
		env := Environment{}
		expectedCode := 126
		expectedStdout := "testdata/env: testdata/env: Is a directory\n\n"

		returnCode := -1
		stdout := captureStdout(t, func() {
			returnCode = RunCmd(cmd, env)
		})

		require.Equal(t, expectedCode, returnCode)
		require.Equal(t, expectedStdout, stdout)
	})
}
