package main

import (
	"fmt"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	for envName, envValue := range env {
		if envValue.NeedRemove {
			err := os.Unsetenv(envName)
			if err != nil {
				fmt.Printf("can't unset env variabe: %v, err: %v", envName, err)
				returnCode = 1
			}
			continue
		}

		err := os.Setenv(envName, envValue.Value)
		if err != nil {
			fmt.Printf("can't set env variabe: %v, value: %v, err: %v", envName, envValue.Value, err)
			returnCode = 1
		}
	}

	if returnCode == 0 {
		cmdExec := exec.Command(cmd[0], cmd[1:]...) // nolint:gosec
		cmdExec.Stdin = os.Stdin
		cmdExec.Stdout = os.Stdout
		err := cmdExec.Run()
		if err != nil {
			fmt.Printf("can't run: %v, args: %v, err: %v", cmd[0], cmd[1:], err)
		}
		returnCode = cmdExec.ProcessState.ExitCode()
	}

	return
}
