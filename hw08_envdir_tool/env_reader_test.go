package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestReadDir_positive(t *testing.T) {
	defer goleak.VerifyNone(t)
	//testDir, err := os.CreateTemp()
	var expRes1 Environment = map[string]EnvValue{
		"BAR":              {"bar", false},
		"EMPTY":            {"ss", false},
		"FOO":              {"   foo\nwith new line", false},
		"HELLO":            {`"hello"`, false},
		"UNSET":            {"", true},
		"lower_case":       {"123", false},
		"with_digits_0123": {"testdata/env/with_digits_0123", false},
	}

	for _, tc := range []struct {
		name   string
		path   string
		expRes Environment
	}{
		{
			name: "cases:" +
				"embedded empty dir" +
				"file with invalid name" +
				"valid files" +
				"embedded dir with duplicate files" +
				"lower case file name" +
				"file name consist from digits",
			path:   "testdata/env",
			expRes: expRes1,
		},
		{
			name:   "pass empty dir",
			path:   "testdata/env/empty_dir",
			expRes: make(Environment),
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			res, err := ReadDir(tc.path)
			require.NoError(t, err)
			assert.Equal(t, tc.expRes, res)
		})
	}
}

func TestReadDir_negative(t *testing.T) {
	defer goleak.VerifyNone(t)
	for _, tc := range []struct {
		name       string
		path       string
		errContain string
	}{
		{
			name:       "pass file",
			path:       "testdata/echo.sh",
			errContain: "not a directory",
		},
		{
			name:       "pass .",
			path:       ".",
			errContain: "empty dir",
		},
		{
			name:       "pass empty string",
			path:       "",
			errContain: "empty dir",
		},
		{
			name:       "pass empty string",
			path:       "",
			errContain: "empty dir",
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			res, err := ReadDir(tc.path)
			require.Error(t, err)
			require.True(t, strings.Contains(err.Error(), tc.errContain))
			assert.Nil(t, res)
		})
	}
}
