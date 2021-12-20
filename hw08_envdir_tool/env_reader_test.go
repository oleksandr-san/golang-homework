package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("testdata", func(t *testing.T) {
		env, err := ReadDir("testdata/env")
		require.Nil(t, err)

		expectedEnv := Environment{
			"BAR":   EnvValue{"bar", false},
			"EMPTY": EnvValue{"", false},
			"FOO":   EnvValue{"   foo\nwith new line", false},
			"HELLO": EnvValue{"\"hello\"", false},
			"UNSET": EnvValue{"", true},
		}

		if !cmp.Equal(env, expectedEnv) {
			t.Errorf("Invalid environment map: %s", cmp.Diff(expectedEnv, env))
		}
	})
}
