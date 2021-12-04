package hw02unpackstring

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnpack(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "a4bc2d5e", expected: "aaaabccddddde"},
		{input: "abccd", expected: "abccd"},
		{input: "", expected: ""},
		{input: "aaa0b", expected: "aab"},
		{input: "З=3>", expected: "З===>"},
		{input: " ", expected: " "},
		{input: " 2", expected: "  "},
		{input: " 0", expected: ""},
		{input: "s0", expected: ""},
		{input: "\n3", expected: "\n\n\n"},
		{input: "d\n5abc", expected: "d\n\n\n\n\nabc"},

		{input: `a-4b`, expected: `a----b`},
		{input: `''2'a'`, expected: `''''a'`},
		{input: `qwe\4\5`, expected: `qwe45`},
		{input: `qwe\45`, expected: `qwe44444`},
		{input: `qwe\\5`, expected: `qwe\\\\\`},
		{input: `\\`, expected: `\`},
		{input: `qwe\\\3`, expected: `qwe\3`},
		{input: `\3`, expected: "3"},
		{input: `\32`, expected: "33"},
		{input: `\00`, expected: ""},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			result, err := Unpack(tc.input)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestUnpackInvalidString(t *testing.T) {
	invalidStrings := []string{"3abc", "45", "aaa10b", `\s4`, `\`, `qwe\\\`, `w\e3`, `\444`, `3`, `qw\ne\`}
	for _, tc := range invalidStrings {
		tc := tc
		t.Run(tc, func(t *testing.T) {
			_, err := Unpack(tc)
			assert.Truef(t, errors.Is(err, ErrInvalidString), "actual error %q", err)
		})
	}
}
