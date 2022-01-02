package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var (
	ErrUnsupportedType   = errors.New("unsupported type for validation")
	ErrUnsupportedRule   = errors.New("unsupported validation rule")
	ErrInvalidRuleSyntax = errors.New("invalid validation rule syntax")

	ErrProhibitedValue      = errors.New("value is not allowed")
	ErrStringInvalidLength  = errors.New("invalid string length")
	ErrStringRegexpMismatch = errors.New("value does not match regular expression")
	ErrNumberTooSmall       = errors.New("number value is too small")
	ErrNumberTooBig         = errors.New("number value is too big")
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	errorStrings := make([]string, len(v))
	for i, err := range v {
		errorStrings[i] = err.Err.Error()
	}
	return strings.Join(errorStrings, "; ")
}

type validator func(string, reflect.Value) error

func mergeValidators(validators []validator) validator {
	return func(name string, value reflect.Value) error {
		var mergedErrors ValidationErrors
		for _, v := range validators {
			if err := v(name, value); err != nil {
				var localErrors ValidationErrors
				if errors.As(err, &localErrors) {
					mergedErrors = append(mergedErrors, localErrors...)
				}
			}
		}

		if len(mergedErrors) == 0 {
			return nil
		}

		return mergedErrors
	}
}

func makeStringValidator(tag string) (validator, error) {
	var validators []validator

	for _, rawRule := range strings.Split(tag, "|") {
		rule := strings.SplitN(rawRule, ":", 2)
		if len(rule) != 2 {
			return nil, ErrInvalidRuleSyntax
		}

		switch rule[0] {
		case "len":
			requiredLength, err := strconv.Atoi(rule[1])
			if err != nil {
				return nil, err
			}

			validators = append(validators, func(n string, v reflect.Value) error {
				if len(v.String()) != requiredLength {
					return ValidationErrors{
						ValidationError{n, ErrStringInvalidLength},
					}
				}
				return nil
			})

		case "regexp":
			re, err := regexp.Compile(rule[1])
			if err != nil {
				return nil, err
			}

			validators = append(validators, func(n string, v reflect.Value) error {
				if !re.MatchString(v.String()) {
					return ValidationErrors{
						ValidationError{n, ErrStringRegexpMismatch},
					}
				}
				return nil
			})

		case "in":
			allowedValues := strings.Split(rule[1], ",")
			sort.Strings(allowedValues)

			validators = append(validators, func(n string, v reflect.Value) error {
				value := v.String()
				i := sort.SearchStrings(allowedValues, value)
				if i >= len(allowedValues) || allowedValues[i] != value {
					return ValidationErrors{
						ValidationError{n, ErrProhibitedValue},
					}
				}
				return nil
			})

		default:
			return nil, ErrUnsupportedRule
		}
	}

	return mergeValidators(validators), nil
}

func parseInts(rawValues string) ([]int, error) {
	values := []int{}
	for _, rawValue := range strings.Split(rawValues, ",") {
		value, err := strconv.Atoi(rawValue)
		if err != nil {
			return nil, err
		}

		values = append(values, value)
	}
	return values, nil
}

func makeIntValidator(tag string) (validator, error) {
	var validators []validator

	for _, rawRule := range strings.Split(tag, "|") {
		rule := strings.SplitN(rawRule, ":", 2)
		if len(rule) != 2 {
			return nil, ErrInvalidRuleSyntax
		}

		switch rule[0] {
		case "min":
			minValue, err := strconv.ParseInt(rule[1], 10, 0)
			if err != nil {
				return nil, err
			}

			validators = append(validators, func(n string, v reflect.Value) error {
				if v.Int() < minValue {
					return ValidationErrors{
						ValidationError{n, ErrNumberTooSmall},
					}
				}
				return nil
			})

		case "max":
			maxValue, err := strconv.ParseInt(rule[1], 10, 0)
			if err != nil {
				return nil, err
			}

			validators = append(validators, func(n string, v reflect.Value) error {
				if v.Int() > maxValue {
					return ValidationErrors{
						ValidationError{n, ErrNumberTooBig},
					}
				}
				return nil
			})

		case "in":
			allowedValues, err := parseInts(rule[1])
			if err != nil {
				return nil, err
			}
			sort.Ints(allowedValues)

			validators = append(validators, func(n string, v reflect.Value) error {
				value := int(v.Int())
				i := sort.SearchInts(allowedValues, value)
				if i >= len(allowedValues) || allowedValues[i] != value {
					return ValidationErrors{
						ValidationError{n, ErrProhibitedValue},
					}
				}
				return nil
			})

		default:
			return nil, ErrUnsupportedRule
		}
	}

	return mergeValidators(validators), nil
}

func makeValidator(t reflect.Type, tag string) (validator, error) {
	switch t.Kind() {
	case reflect.String:
		return makeStringValidator(tag)

	case reflect.Int:
		return makeIntValidator(tag)

	case reflect.Slice:
		elemValidator, err := makeValidator(t.Elem(), tag)
		if err != nil {
			return nil, err
		}

		return func(name string, v reflect.Value) error {
			var sliceErrors ValidationErrors
			for i := 0; i < v.Len(); i++ {
				elemName := fmt.Sprintf("%s[%d]", name, i)
				if err := elemValidator(elemName, v.Index(i)); err != nil {
					var fieldErrors ValidationErrors
					if errors.As(err, &fieldErrors) {
						sliceErrors = append(sliceErrors, fieldErrors...)
					} else {
						return err
					}
				}
			}
			return sliceErrors
		}, nil

	case reflect.Struct:
		if tag == "nested" {
			return func(name string, v reflect.Value) error {
				return validateStruct(name+".", v)
			}, nil
		}

	// Ugly code to silence switch statement exhaustiveness check
	case reflect.Array, reflect.Bool, reflect.Chan,
		reflect.Complex128, reflect.Complex64,
		reflect.Float32, reflect.Float64,
		reflect.Func, reflect.Int16, reflect.Int32,
		reflect.Int64, reflect.Int8, reflect.Interface,
		reflect.Invalid, reflect.Map, reflect.Ptr,
		reflect.Uint, reflect.Uint16, reflect.Uint32,
		reflect.Uint64, reflect.Uint8, reflect.Uintptr,
		reflect.UnsafePointer:
	default:
	}

	return nil, ErrUnsupportedType
}

func validateStruct(prefix string, s reflect.Value) error {
	var structErrors ValidationErrors

	sType := s.Type()
	for i := 0; i < s.NumField(); i++ {
		sf := sType.Field(i)

		tag, found := sf.Tag.Lookup("validate")
		if !found {
			continue
		}

		validator, err := makeValidator(sf.Type, tag)
		if err != nil {
			return err
		}

		if err := validator(prefix+sf.Name, s.Field(i)); err != nil {
			var fieldErrors ValidationErrors
			if errors.As(err, &fieldErrors) {
				structErrors = append(structErrors, fieldErrors...)
			} else {
				return err
			}
		}
	}

	if len(structErrors) != 0 {
		return structErrors
	}

	return nil
}

func Validate(v interface{}) error {
	s := reflect.ValueOf(v)
	if s.Kind() != reflect.Struct {
		return ErrUnsupportedType
	}

	return validateStruct("", s)
}
