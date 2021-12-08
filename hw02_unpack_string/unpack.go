package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var (
	ErrInvalidString     = errors.New("invalid string")
	ErrUnknownStateError = errors.New("unknown state error")
)

type unpackState int8

const (
	expectAny unpackState = iota
	escapeCharFound
	prevCharStored
)

func Unpack(packed string) (string, error) {
	var sb strings.Builder
	var state unpackState
	var prevChar rune

	for _, char := range packed {
		switch state {
		case expectAny:
			switch {
			case char == '\\':
				state = escapeCharFound
			case unicode.IsDigit(char):
				return "", ErrInvalidString
			default:
				prevChar = char
				state = prevCharStored
			}
		case escapeCharFound:
			prevChar = char
			state = prevCharStored
		case prevCharStored:
			switch {
			case char == '\\':
				sb.WriteRune(prevChar)
				state = escapeCharFound
			case unicode.IsDigit(char):
				multiplier, err := strconv.Atoi(string(char))
				if err != nil {
					return "", err
				}
				sb.WriteString(strings.Repeat(string(prevChar), multiplier))
				state = expectAny
			default:
				sb.WriteRune(prevChar)
				prevChar = char
			}
		default:
			return "", ErrUnknownStateError
		}
	}

	switch state {
	case prevCharStored:
		sb.WriteRune(prevChar)
	case escapeCharFound:
		return "", ErrInvalidString
	case expectAny:
	}

	return sb.String(), nil
}
