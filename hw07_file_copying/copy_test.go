package main

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCopy_PositiveCases(t *testing.T) {
	for _, tc := range []struct {
		name   string
		from   string
		offset int64
		limit  int64
		result string
	}{
		{
			name:   "offset & limit = 0; expect full copy",
			from:   "testdata/input.txt",
			offset: 0,
			limit:  0,
			result: "testdata/out_offset0_limit0.txt",
		},
		{
			name:   "offset = 0, limit = 10; expect copy 10 bytes",
			from:   "testdata/input.txt",
			offset: 0,
			limit:  10,
			result: "testdata/out_offset0_limit10.txt",
		},
		{
			name:   "offset = 0, limit = 1000; expect copy 1000 bytes",
			from:   "testdata/input.txt",
			offset: 0,
			limit:  1000,
			result: "testdata/out_offset0_limit1000.txt",
		},
		{
			name:   "offset = 0, limit = 10000 > file from size; expect copy full file",
			from:   "testdata/input.txt",
			offset: 0,
			limit:  10000,
			result: "testdata/out_offset0_limit10000.txt",
		},
		{
			name:   "offset = 100, limit = 1000; expect copy 1000 bytes",
			from:   "testdata/input.txt",
			offset: 100,
			limit:  1000,
			result: "testdata/out_offset100_limit1000.txt",
		},
		{
			name:   "offset = 6000, limit = 1000; expect copy fileSize - offset bytes",
			from:   "testdata/input.txt",
			offset: 6000,
			limit:  1000,
			result: "testdata/out_offset6000_limit1000.txt",
		},
		{
			name:   "offset == fileSize == 6617, limit = 0; expect copy 0 bytes",
			from:   "testdata/input.txt",
			offset: 6617,
			limit:  0,
			result: "testdata/out_offset6617_limit0.txt",
		},
		{
			name:   "offset = 0, limit = 6617; expect copy full file",
			from:   "testdata/input.txt",
			offset: 0,
			limit:  6617,
			result: "testdata/out_offset0_limit0.txt",
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			tempFile, err := filepath.Abs("testdata")
			require.NoError(t, err)
			fileTo, err := os.CreateTemp(tempFile, "test_case")
			require.NoError(t, err)
			expResFile, err := os.Open(tc.result)
			require.NoError(t, err)
			defer func() {
				_ = fileTo.Close()
				_ = expResFile.Close()
				_ = os.Remove(fileTo.Name())
			}()
			expRes, err := ioutil.ReadAll(expResFile)
			require.NoError(t, err)

			err = Copy(tc.from, fileTo.Name(), tc.offset, tc.limit)
			require.NoError(t, err)

			res, err := ioutil.ReadAll(fileTo)
			require.NoError(t, err)
			assert.Equal(t, expRes, res)
		})
	}
}

func TestCopy_NegativeCases(t *testing.T) {
	for _, tc := range []struct {
		name     string
		from     string
		offset   int64
		limit    int64
		expError error
	}{
		{
			name:     "file with no size",
			from:     "/dev/urandom",
			offset:   0,
			limit:    0,
			expError: ErrUnsupportedFile,
		},
		{
			name:     "negative offset",
			from:     "testdata/input.txt",
			offset:   -1,
			limit:    0,
			expError: ErrOffsetNegativeValue,
		},
		{
			name:     "over offset",
			from:     "testdata/input.txt",
			offset:   6672,
			limit:    0,
			expError: ErrOffsetExceedsFileSize,
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := Copy(tc.from, "/tmp/", tc.offset, tc.limit)
			assert.True(t, errors.Is(err, tc.expError))
		})
	}
}
