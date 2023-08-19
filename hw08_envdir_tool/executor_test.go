package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("common case; expect no error", func(t *testing.T) {
		var (
			env Environment = map[string]EnvValue{
				"BAR":   {"bar", false},
				"EMPTY": {"", false},
				"FOO":   {"   foo\nwith new line", false},
				"HELLO": {`"hello"`, false},
				"UNSET": {"", true},
			}
			expected = `HELLO is ("hello")
BAR is (bar)
FOO is (   foo
with new line)
UNSET is ()
ADDED is ()
EMPTY is ()
arguments are arg1=1 arg2=2
`
		)

		path, err := filepath.Abs("testdata/echo.sh")
		require.NoError(t, err)
		cmd := []string{path, "arg1=1", "arg2=2"}

		out := new(bytes.Buffer)
		code, errors := RunCmd(os.Stdin, out, cmd, env)
		require.Zero(t, code)
		require.Empty(t, errors)
		assert.Equal(t, expected, out.String())
	})

	t.Run("return valid code", func(t *testing.T) {
		code, errors := RunCmd(os.Stdin, os.Stdout, []string{"/bin/bash", "-c", "exit 5"}, Environment{})
		require.NotEmpty(t, errors)
		require.Equal(t, 5, code)
	})
}
