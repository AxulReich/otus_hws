package hw02unpackstring

import (
	"errors"
	"fmt"
	"strconv"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

// Unpack checking & building in same loop
func Unpack(inStr string) (string, error) {
	const shieldingChar rune = 92 // stands for '/' character
	var (
		metShieldingChar bool
		inStrRuned       = []rune(inStr)

		resultRune = make([]rune, 0, 2*cap(inStrRuned))
	)

	for idx, r := range inStrRuned {
		if unicode.IsDigit(r) {
			if (idx > 1 && (inStrRuned[idx-2] == shieldingChar)) || (idx > 0 && !unicode.IsDigit(inStrRuned[idx-1])) {
				if metShieldingChar {
					resultRune = append(resultRune, r)
					metShieldingChar = false

				} else {
					rInt, err := strconv.Atoi(string(r))
					if err != nil {
						return "", fmt.Errorf("error while Atoi: %w", err)
					}
					curIdxInResult := len(resultRune)-1

					if rInt == 0 {
						resultRune = resultRune[:curIdxInResult]
						continue
					}

					if rInt == 1 {
						continue
					}

					runeToCopy := resultRune[curIdxInResult]

					for i := 0; i < rInt-1; i++ {
						resultRune = append(resultRune, runeToCopy)
					}

				}
			} else {
				return "", ErrInvalidString
			}
		} else {
			if r == shieldingChar {
				if !metShieldingChar {
					metShieldingChar = true
					continue
				}
				metShieldingChar = false
			}
			resultRune = append(resultRune, r)
		}
	}
	if metShieldingChar {
		return "", fmt.Errorf("extra \\: %w", ErrInvalidString)
	}
	return string(resultRune), nil
}

// Un2pack checking & building in same loop
func Un2pack(inStr string) (string, error) {
	const shieldingChar rune = 92 // stands for '/' character
	var (
		metShieldingChar bool
		inStrRuned       = []rune(inStr)
		multiplied bool
		resultRune = make([]rune, 0, 2*cap(inStrRuned))
	)

	for idx, r := range inStrRuned {
		if unicode.IsDigit(r) {
			if idx > 0 {
				if metShieldingChar {
					resultRune = append(resultRune, r)
					metShieldingChar = false
					continue
				}
				if multiplied {
					return "", ErrInvalidString
				}
				rInt, err := strconv.Atoi(string(r))
				if err != nil {
					return "", fmt.Errorf("error while Atoi: %w", err)
				}
				curIdxInResult := len(resultRune)-1

				if rInt == 0 {
					resultRune = resultRune[:curIdxInResult]
					multiplied = true
					continue
				}

				if rInt == 1 {
					multiplied = true
					continue
				}

				runeToCopy := resultRune[curIdxInResult]

				for i := 0; i < rInt-1; i++ {
					resultRune = append(resultRune, runeToCopy)
				}
				multiplied = true

			} else {
				return "", ErrInvalidString
			}
		} else {
			if r == shieldingChar {
				if !metShieldingChar {
					metShieldingChar = true
					continue
				}
				metShieldingChar = false
			}
			if metShieldingChar {
				return "", ErrInvalidString
			}
			resultRune = append(resultRune, r)
			multiplied = false
		}
	}

	if metShieldingChar {
		return "", fmt.Errorf("extra \\: %w", ErrInvalidString)
	}
	return string(resultRune), nil
}