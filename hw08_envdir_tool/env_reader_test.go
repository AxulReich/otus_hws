package main

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestReadDir_positive(t *testing.T) {
	var expRes1 Environment = map[string]EnvValue{
		"BAR":   {"bar", false},
		"EMPTY": {"ss", false},
		"FOO":   {"   foo\nwith new line", false},
		"HELLO": {`"hello"`, false},
		"UNSET": {"", true},
	}

	for _, tc := range []struct {
		name   string
		path   string
		expRes Environment
	}{
		{
			name:   "common case with embedded empty dir, file with invalid name , valid files, embedded dir with duplicate files",
			path:   "testdata/env",
			expRes: expRes1,
		},
		{
			name:   "pass empty dir",
			path:   "testdata/env/empty_dir",
			expRes: make(Environment),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			res, err := ReadDir(tc.path)
			require.NoError(t, err)
			assert.Equal(t, res, tc.expRes)
		})
	}
}

func TestReadDir_negative(t *testing.T) {
	const envDirPath = "testdata/echo.sh"
	res, err := ReadDir(envDirPath)
	require.Error(t, err)
	assert.Nil(t, res)
}
