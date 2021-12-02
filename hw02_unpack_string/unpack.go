package hw02unpackstring

import (
	"errors"
	"fmt"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

const notAddMultiplier int64 = 0

// Unpack checking & building in same loop.
func Unpack(inStr string) (string, error) {
	const shieldingChar rune = 92 // stands for '\' character.

	var (
		metShieldingChar bool
		multiplied       bool
		inStrRuned       = []rune(inStr)
		resultRune       []rune
	)

	for idx, r := range inStrRuned {
		var (
			runeToAdd       = r
			mult      int64 = 1
		)
		switch {
		case unicode.IsDigit(r):
			switch {
			case metShieldingChar:
				metShieldingChar = false
			default:
				if idx > 0 && !multiplied {
					multiplied = true
					mult = int64(r-'0') - 1
					runeToAdd = resultRune[len(resultRune)-1]
				} else {
					return "", fmt.Errorf("digit at [0] position in string or digits in a row waithout \\ err: %w", ErrInvalidString)
				}
			}
		default:
			multiplied = false
			if r == shieldingChar {
				switch {
				case metShieldingChar:
					metShieldingChar = false
				default:
					metShieldingChar = true
					mult = notAddMultiplier
				}
			} else if metShieldingChar {
				return "", fmt.Errorf("\\ char shield only digits err: %w", ErrInvalidString)
			}
		}
		resultRune = add(resultRune, runeToAdd, mult)
	}
	if metShieldingChar {
		return "", fmt.Errorf("extra \\ err: %w", ErrInvalidString)
	}
	return string(resultRune), nil
}

func add(target []rune, runeToCopy rune, n int64) []rune {
	switch {
	case n < 0:
		target = target[:len(target)-1]
	case n > 0:
		for i := 0; i < int(n); i++ {
			target = append(target, runeToCopy)
		}
	case n == notAddMultiplier:
	}
	return target
}
