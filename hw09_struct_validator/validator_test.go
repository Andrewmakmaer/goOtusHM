package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:5"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
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
		Code int    `validate:"in:200,404,500|len:3"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in:          "It is not struct",
			expectedErr: NoStructTypeError{"type error: given object of type 'string', expected struct"},
		},
		{
			in: User{
				ID:     "12345",
				Name:   "John Doe",
				Age:    30,
				Email:  "john@example.com",
				Role:   "admin",
				Phones: []string{"12345678901", "90987654321"},
			},
			expectedErr: nil,
		},
		{
			in: User{
				ID:     "123456789012345678901234567890123456",
				Name:   "Jane Doe",
				Age:    17,
				Email:  "invalid-email",
				Role:   "user",
				Phones: []string{"1234", "53587654321"},
			},
			expectedErr: ValidationErrors{
				{Field: "ID", Err: fmt.Errorf("%w, length of the 123456789012345678901234567890123456 is not equal 5", ErrValidate)},
				{Field: "Age", Err: fmt.Errorf("%w, number 17 less that 18", ErrValidate)},
				{Field: "Email", Err: fmt.Errorf("%w, invalid-email is not match for ^\\w+@\\w+\\.\\w+$ expression", ErrValidate)},
				{Field: "Role", Err: fmt.Errorf("%w, value user not in admin,stuff", ErrValidate)},
				{Field: "Phones", Err: fmt.Errorf("%w, length of the 1234 is not equal 11", ErrValidate)},
				// {Field: "Phones", Err: fmt.Errorf("%w, length of the 535 is not equal 11", ErrValidate)},
			},
		},
		{
			in: Response{
				Code: 200,
				Body: "Hello",
			},
			expectedErr: ValidatorSetError{
				{Field: "Code", Err: fmt.Errorf("program fail: len validator is not supported for type int")},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			// t.Parallel()
			err := Validate(tt.in)
			require.Truef(t, errors.Is(tt.expectedErr, err), "expected error: %q, but actual error %q", tt.expectedErr, err)
		})
	}
}
