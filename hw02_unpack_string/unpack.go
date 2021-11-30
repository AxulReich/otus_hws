package hw02unpackstring

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

// Unpack checking & building in same loop
func Unpack(inStr string) (string, error) {
	const shieldingChar rune = 92 // stands for '/' character
	var (
		metShieldingChar bool
		inStrRuned       = []rune(inStr)

		result     strings.Builder
		resultRune = make([]rune, 0, 2*cap(inStrRuned))
	)

	for idx, r := range inStrRuned {
		if unicode.IsDigit(r) {
			if (idx > 1 && (inStrRuned[idx-2] == shieldingChar)) || (idx > 0 && !unicode.IsDigit(inStrRuned[idx-1])) {
				if metShieldingChar {
					result.WriteRune(r)
					resultRune = append(resultRune, r)
					metShieldingChar = false

				} else {
					rInt, err := strconv.Atoi(string(r))
					if err != nil {
						return "", fmt.Errorf("error while Atoi: %w", err)
					}

					if rInt == 0 {
						resultRune = resultRune[:idx-1]
						continue
					}

					if rInt == 1 {
						continue
					}

					runeToCopy := resultRune[len(resultRune)-1]
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

	return string(resultRune), nil
}
