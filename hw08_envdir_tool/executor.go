package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

func updateEnv(env []string, changes Environment) []string {
	var changedEnv []string

	for _, envItem := range env {
		envKey := strings.SplitN(envItem, "=", 2)[0]
		if _, found := changes[envKey]; !found {
			changedEnv = append(changedEnv, envItem)
		}
	}

	for envKey, envVal := range changes {
		if !envVal.NeedRemove {
			changedEnv = append(changedEnv, envKey+"="+envVal.Value)
		}
	}

	return changedEnv
}

func RunCmd(cmd []string, env Environment) (returnCode int) {
	name := cmd[0]
	args := []string{}
	if len(cmd) > 1 {
		args = cmd[1:]
	}

	command := exec.Command(name, args...)
	command.Env = updateEnv(os.Environ(), env)

	stdin, err := command.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		defer stdin.Close()
		io.Copy(stdin, os.Stdin)
	}()

	output, err := command.CombinedOutput()
	if output != nil {
		fmt.Printf("%s\n", output)
	}

	var ee *exec.ExitError
	if errors.As(err, &ee) {
		return ee.ExitCode()
	} else if err != nil {
		return 1
	}

	return 0
}
