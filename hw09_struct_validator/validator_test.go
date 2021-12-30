package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}

	Context struct {
		FirstResponse  Response `validate:"nested"`
		SecondResponse Response `validate:"nested"`
		Token          Token    `validate:"nested"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: App{"0.0.0"},
		},
		{
			in: App{"0"},
			expectedErr: ValidationErrors{
				ValidationError{"Version", ErrStringInvalidLength},
			},
		},
		{
			in: User{
				ID:   "ff285bd7-e473-4639-98d3-af5171790621",
				Name: "John", Age: 49, Email: "at@at.at",
				Role: "admin", Phones: []string{
					"+1111111111",
				}, meta: json.RawMessage{},
			},
		},
		{
			in: User{
				ID:   "ff285bd7-e473-4639-98d3-af5171790621",
				Name: "Alex", Age: 17, Email: "invalid",
				Role: "stuff", Phones: []string{}, meta: json.RawMessage{},
			},
			expectedErr: ValidationErrors{
				ValidationError{"Age", ErrNumberTooSmall},
				ValidationError{"Email", ErrStringRegexpMismatch},
			},
		},
		{
			in: User{
				ID:   "ff285bd7-e473-4639-98d3-af5171790621",
				Name: "Billy", Age: 52, Email: "at@at.at",
				Role: "pleb", Phones: []string{
					"+1111111111",
					"+222222222",
				}, meta: json.RawMessage{},
			},
			expectedErr: ValidationErrors{
				ValidationError{"Age", ErrNumberTooBig},
				ValidationError{"Role", ErrProhibitedValue},
				ValidationError{"Phones[1]", ErrStringInvalidLength},
			},
		},
		{
			in: Context{
				Response{Code: 200, Body: "{}"},
				Response{Code: 418, Body: "{}"},
				Token{},
			},
			expectedErr: ValidationErrors{
				ValidationError{"SecondResponse.Code", ErrProhibitedValue},
			},
		},
		{
			in: struct {
				TwoDigits string `validate:"regexp:^\\d+$|len:20"`
			}{"2a3a"},
			expectedErr: ValidationErrors{
				ValidationError{"TwoDigits", ErrStringRegexpMismatch},
				ValidationError{"TwoDigits", ErrStringInvalidLength},
			},
		},
		{
			in:          &Token{},
			expectedErr: ErrUnsupportedType,
		},
		{
			in: struct {
				Map map[string]string `validate:"len:10"`
			}{},
			expectedErr: ErrUnsupportedType,
		},
		{
			in: struct {
				V int `validate:"len:10"`
			}{},
			expectedErr: ErrUnsupportedRule,
		},
		{
			in: struct {
				S []string `validate:"max:10"`
			}{},
			expectedErr: ErrUnsupportedRule,
		},
		{
			in: struct {
				V int `validate:"len"`
			}{},
			expectedErr: ErrInvalidRuleSyntax,
		},
		{
			in: struct {
				S string `validate:"max"`
			}{},
			expectedErr: ErrInvalidRuleSyntax,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			require.Equal(t, tt.expectedErr, err)
		})
	}
}
