package main

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

const (
	rootPath         = "testdata/env"
	fileNameToCopy   = "BAR"
	dirNameDupl      = "dir_with_dupl"
	emptyDirName     = "empty_dir"
	invalidNameFile  = "INVALID=FILE=NAME"
	lowCaseFileName  = "lower_case"
	lowCaseFileValue = "123"
	withDigFileName  = "with_digits_0123"
)

func closeFileAndDelete(t *testing.T, file *os.File) {
	t.Helper()

	_ = file.Close()
	_ = os.Remove(file.Name())
}

func TestReadDir_positive(t *testing.T) {
	defer goleak.VerifyNone(t)

	invalidFile, err := os.Create(filepath.Join(rootPath, invalidNameFile))
	require.NoError(t, err)
	defer closeFileAndDelete(t, invalidFile)
	lowCaseFile, err := os.Create(filepath.Join(rootPath, lowCaseFileName))
	require.NoError(t, err)
	nBytes, err := lowCaseFile.Write([]byte(lowCaseFileValue))
	require.NotZero(t, nBytes)
	require.NoError(t, err)
	defer closeFileAndDelete(t, lowCaseFile)
	withDigFile, err := os.Create(filepath.Join(rootPath, withDigFileName))
	require.NoError(t, err)
	defer closeFileAndDelete(t, withDigFile)

	duplDir, err := os.MkdirTemp(rootPath, dirNameDupl)
	require.NoError(t, err)
	duplFile, err := os.Create(filepath.Join(duplDir, fileNameToCopy))
	require.NoError(t, err)
	source, err := os.Open(filepath.Join(rootPath, fileNameToCopy))
	require.NoError(t, err)
	nBytesCopy, err := io.Copy(duplFile, source)
	require.NoError(t, err)
	require.NotZero(t, nBytesCopy)
	emptyDir, err := os.MkdirTemp(rootPath, emptyDirName)
	require.NoError(t, err)

	defer func() {
		_ = source.Close()
		_ = duplFile.Close()
		_ = os.Remove(duplFile.Name())
		_ = os.Remove(duplDir)
		_ = os.Remove(emptyDir)
	}()

	var expRes1 Environment = map[string]EnvValue{
		"BAR":              {"bar", false},
		"EMPTY":            {"ss", false},
		"FOO":              {"   foo\nwith new line", false},
		"HELLO":            {`"hello"`, false},
		"UNSET":            {"", true},
		"lower_case":       {lowCaseFileValue, false},
		"with_digits_0123": {"", true},
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
			path:   emptyDir,
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
