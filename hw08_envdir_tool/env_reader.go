package main

import (
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

		if len(content) == 0 {
			env[entry.Name()] = EnvValue{NeedRemove: true}
		} else {
			//content = bytes.Split(content, []byte{32})[0]
			first_line := strings.Split(string(content), "\n")[0]

			value := strings.TrimRight(first_line, "\t ")
			value = strings.ReplaceAll(value, string(rune(0)), "\n")

			env[entry.Name()] = EnvValue{Value: value}
		}
	}

	return env, nil
}
