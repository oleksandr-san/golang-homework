package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

var (
	ErrUnsupportedType     = errors.New("unsupported type for validation")
	ErrInvalidStringLength = errors.New("invalid string length")
)

type RuleError struct {
	Field string
	Rule  string
	Err   error
}

func (e *RuleError) Error() string {
	panic("implement me")
}

func (e *RuleError) Unwrap() error {
	return e.Err
}

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	panic("implement me")
}

func Validate(v interface{}) error {
	s := reflect.ValueOf(v)
	if s.Kind() != reflect.Struct {
		return ErrUnsupportedType
	}

	sType := s.Type()

	var validationErrors ValidationErrors
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		fType := sType.Field(i)

		rawRules, found := fType.Tag.Lookup("validate")
		if !found {
			continue
		}

		for _, rule := range strings.Split(rawRules, "|") {
			ruleComponents := strings.SplitN(rule, ":", 2)
			if ruleComponents[0] == "len" {
				requiredLength, err := strconv.Atoi(ruleComponents[1])
				if err != nil {
					return &RuleError{
						Field: fType.Name,
						Rule:  ruleComponents[0],
						Err:   err,
					}
				}

				if len(f.String()) != requiredLength {
					validationError := ValidationError{
						Field: fType.Name,
						Err:   ErrInvalidStringLength,
					}
					validationErrors = append(validationErrors, validationError)
				}
			}
		}

		fmt.Println(f, fType, rawRules)
	}

	if len(validationErrors) != 0 {
		return validationErrors
	}

	return nil
}
