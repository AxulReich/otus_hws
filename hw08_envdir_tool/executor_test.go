package main

import (
	"os"
	"testing"
)

func TestRunCmd_Negative(t *testing.T) {
	for _, tc := range []struct {
		name string
		cmd  []string
		env  Environment
	}{
		{
			name: "pass nil cmd, expect err",
			cmd:  nil,
			env:  make(Environment),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, _ = RunCmd(os.Stdin, os.Stdout, tc.cmd, tc.env)
		})
	}
}
