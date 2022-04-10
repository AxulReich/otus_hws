package main

import (
	"fmt"
	"io"
	"os"

	//"os/exec"
	exec "golang.org/x/sys/execabs"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(in io.Reader, out io.Writer, cmd []string, env Environment) (int, []error) {
	var (
		returnCode int
		errors     []error
	)

	if in == nil || out == nil {
		errors = append(errors, fmt.Errorf("RunCmd error: passed in:%v or out:%v is nil", in, out))
		returnCode = 1
	}

	if len(cmd) < 1 || env == nil {
		errors = append(errors, fmt.Errorf("RunCmd error: passed cmd: %v is empty or env: %v is nil", cmd, env))
		returnCode = 1
	}

	for envName, envValue := range env {
		if envValue.NeedRemove {
			err := os.Unsetenv(envName)
			if err != nil {
				errors = append(errors, fmt.Errorf("can't unset env variabe: %v, err: %v", envName, err))
				returnCode = 1
			}
			continue
		}

		err := os.Setenv(envName, envValue.Value)
		if err != nil {
			errors = append(errors, fmt.Errorf("can't set env variabe: %v, value: %v, err: %v", envName, envValue.Value, err))
			returnCode = 1
		}
	}

	if len(errors) == 0 {
		cmdExec := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec
		cmdExec.Stdin = in
		cmdExec.Stdout = out
		err := cmdExec.Run()
		if err != nil {
			errors = append(errors, fmt.Errorf("can't run: %v, args: %v, err: %v", cmd[0], cmd[1:], err))
		}
		returnCode = cmdExec.ProcessState.ExitCode()
	}

	return returnCode, errors
}
