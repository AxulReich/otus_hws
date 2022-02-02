package hw02unpackstring

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Pack(inStr string) (string, error) {
	// TODO: implement me
	return inStr, nil
}


// Unpack checking & building in same loop.
func Unpack(inStr string) (string, error) {
	const shieldingChar rune = 92 // stands for '\' character.

	var (
		metShieldingChar bool
		multiplied       bool
		resultRune       []rune
	)

	for idx, r := range inStr {
		var (
			runeToAdd = r
			mult      = 1
		)
		switch {
		case unicode.IsDigit(r):
			switch {
			case metShieldingChar:
				metShieldingChar = false
			default:
				if idx > 0 && !multiplied {
					multiplied = true
					mult = int(r-'0') - 1
					runeToAdd = resultRune[len(resultRune)-1]
				} else {
					return "", fmt.Errorf("digit at [0] position in string or digits in a row without \\ err: %w", ErrInvalidString)
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
					mult = 0
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

func add(target []rune, runeToCopy rune, n int) []rune {
	switch {
	case n < 0:
		target = target[:len(target)-1]
	case n > 0:
		for i := 0; i < n; i++ {
			target = append(target, runeToCopy)
		}
	default:
	}
	return target
}

func MasterUnpack(inputString string) (string, error) {
	const escapeSymbol string = "\\"

	var (
		resultBuilder strings.Builder
		targetToRepeat string
		nextSymbolEscaped bool
	)

	for _, symbolRune := range inputString {
		currentSymbol := string(symbolRune)
		switch {
		case nextSymbolEscaped:
			if !(unicode.IsDigit(symbolRune) || currentSymbol == escapeSymbol) {
				return "", ErrInvalidString
			}
			targetToRepeat = currentSymbol
			nextSymbolEscaped = false

		case currentSymbol == escapeSymbol:
			resultBuilder.WriteString(targetToRepeat)
			targetToRepeat = ""
			nextSymbolEscaped = true

		case unicode.IsDigit(symbolRune):
			if targetToRepeat == "" {
				return "", ErrInvalidString
			}
			repeatCount, _ := strconv.Atoi(currentSymbol)
			resultBuilder.WriteString(strings.Repeat(targetToRepeat, repeatCount))
			targetToRepeat = ""

		default:
			resultBuilder.WriteString(targetToRepeat)
			targetToRepeat = currentSymbol
		}
	}
	if nextSymbolEscaped {
		return "", ErrInvalidString
	}
	resultBuilder.WriteString(targetToRepeat)
	return resultBuilder.String(), nil
}