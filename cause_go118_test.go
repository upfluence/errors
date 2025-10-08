//go:build go1.18

package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockError struct{}

func (mockError) Error() string { return "foo" }

func TestIsOfType(t *testing.T) {
	for _, tt := range []struct {
		input error

		isError     bool
		isMockError bool
	}{
		{}, // nil is not of any compilable type
		{
			input:   New("foo"),
			isError: true,
		},
		{
			input:       mockError{},
			isError:     true,
			isMockError: true,
		},
		{
			input:       Wrap(mockError{}, "wrapping"),
			isError:     true,
			isMockError: true,
		},
	} {
		assert.Equal(t, tt.isError, IsOfType[error](tt.input))
		assert.Equal(t, tt.isMockError, IsOfType[mockError](tt.input))
	}
}

func TestAsType(t *testing.T) {
	for _, tt := range []struct {
		name     string
		input    error
		expected error
	}{
		{
			name: "nil error",
		},
		{
			name:     "mockError",
			input:    mockError{},
			expected: mockError{},
		},
		{
			name:     "wrapped mockError",
			input:    Wrap(mockError{}, "wrapping"),
			expected: mockError{},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			m, ok := AsType[mockError](tt.input)

			assert.Equal(t, tt.expected != nil, ok)

			if ok {
				assert.Equal(t, tt.expected, m)
			}
		})
	}
}
