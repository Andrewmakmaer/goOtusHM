package hw09structvalidator

import (
	"encoding/json"
	"errors"
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
		in interface{}
	}{
		{
			in: "It is not struct",
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
		},
		{
			in: User{
				Name:   "Jane Doe",
				Age:    17,
				Email:  "invalid-email",
				Role:   "user",
				Phones: []string{"1234", "535"},
			},
		},
		{
			in: Response{
				Code: 200,
				Body: "Hello",
			},
		},
	}

	t.Run("Sent not structure", func(t *testing.T) {
		err := Validate(tests[0].in)
		require.Truef(t, errors.As(err, &NoStructTypeError{}), "expected error: %q, but actual error %q", err, NoStructTypeError{})
	})

	t.Run("Run correct struct", func(t *testing.T) {
		err := Validate(tests[1].in)
		require.Truef(t, errors.Is(err, nil), "expected error: %q, but actual error %q", err, nil)
	})

	t.Run("Run wrong struct", func(t *testing.T) {
		err := Validate(tests[2].in)
		require.Truef(t, errors.As(err, &ValidationErrors{}), "expected error: %q, but actual error %q", err, nil)
	})

	t.Run("Run wrong struct", func(t *testing.T) {
		err := Validate(tests[3].in)
		require.Truef(t, errors.As(err, &ValidationErrors{}), "expected error: %q, but actual error %q", err, nil)
	})
}
