package main

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

func makeEnvValue(value []byte) EnvValue {
	if len(value) == 0 {
		return EnvValue{NeedRemove: true}
	}

	firstLine := string(bytes.Split(value, []byte{10})[0])
	envValue := strings.TrimRight(firstLine, "\t ")
	envValue = strings.ReplaceAll(envValue, string(rune(0)), "\n")

	return EnvValue{Value: envValue}
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	env := Environment{}
	for _, entry := range entries {
		if entry.IsDir() || strings.ContainsAny(entry.Name(), "=") {
			continue
		}

		content, err := ioutil.ReadFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			return env, err
		}

		env[entry.Name()] = makeEnvValue(content)
	}

	return env, nil
}
