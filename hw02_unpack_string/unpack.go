package hw02unpackstring

import (
	"errors"
	"fmt"
	"strconv"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

// Unpack checking & building in same loop.
func Unpack(inStr string) (string, error) {
	const shieldingChar rune = 92 // stands for '\' character.
	var (
		metShieldingChar bool
		inStrRuned       = []rune(inStr)
		multiplied       bool
		resultRune       = make([]rune, 0, 2*cap(inStrRuned))
	)

	for idx, r := range inStrRuned {
		if unicode.IsDigit(r) {
			if idx > 0 && !multiplied {
				if metShieldingChar {
					resultRune = append(resultRune, r)
					metShieldingChar = false
					continue
				}
				rInt, err := strconv.Atoi(string(r))
				if err != nil {
					return "", fmt.Errorf("error while Atoi err: %w", err)
				}
				curIdxInResult := len(resultRune) - 1
				multiplied = true

				if rInt == 0 {
					resultRune = resultRune[:curIdxInResult]
					continue
				}

				runeToCopy := resultRune[curIdxInResult]
				for i := 0; i < rInt-1; i++ {
					resultRune = append(resultRune, runeToCopy)
				}
				continue
			}
			return "", fmt.Errorf("digit at [0] position in string or in a row waithout \\ err: %w", ErrInvalidString)
		}

		if r == shieldingChar {
			if !metShieldingChar {
				metShieldingChar = true
				continue
			}
			metShieldingChar = false
		}
		if metShieldingChar {
			return "", fmt.Errorf("\\ char shield only digits err: %w", ErrInvalidString)
		}
		resultRune = append(resultRune, r)
		multiplied = false
	}

	if metShieldingChar {
		return "", fmt.Errorf("extra \\ err: %w", ErrInvalidString)
	}

	return string(resultRune), nil
}
