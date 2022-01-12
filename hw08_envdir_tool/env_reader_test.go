package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

func TestMakeEnvValue(t *testing.T) {
	tests := []struct {
		id            string
		input         []byte
		expectedValue EnvValue
	}{
		{
			"number value", []byte("1"), EnvValue{Value: "1"},
		},
		{
			"empty value", []byte{}, EnvValue{NeedRemove: true},
		},
		{
			"multiline value", []byte("a\nb\nc"), EnvValue{Value: "a"},
		},
		{
			"value with spaces", []byte("a\t \nb"), EnvValue{Value: "a"},
		},
		{
			"value with 0 bytes", []byte{42, 0, 42}, EnvValue{Value: "*\n*"},
		},
	}

	for _, test := range tests {
		t.Run(test.id, func(t *testing.T) {
			t.Parallel()
			tt := test

			actualValue := makeEnvValue(tt.input)
			require.Equal(t, tt.expectedValue, actualValue)
		})
	}
}

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
