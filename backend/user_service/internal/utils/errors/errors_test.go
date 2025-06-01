package errpkg

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAppError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *AppError
		expected string
	}{
		{
			name:     "with message",
			err:      New("CODE", "test message", nil),
			expected: "test message",
		},
		{
			name:     "empty message",
			err:      New("CODE", "", nil),
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.err.Error())
		})
	}
}

func TestAppError_Unwrap(t *testing.T) {
	baseErr := errors.New("base error")
	tests := []struct {
		name     string
		err      *AppError
		expected error
	}{
		{
			name:     "with wrapped error",
			err:      New("CODE", "msg", baseErr),
			expected: baseErr,
		},
		{
			name:     "without wrapped error",
			err:      New("CODE", "msg", nil),
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.err.Unwrap())
		})
	}
}

func TestNew(t *testing.T) {
	baseErr := errors.New("base error")
	tests := []struct {
		name     string
		code     string
		message  string
		err      error
		expected *AppError
	}{
		{
			name:    "full fields",
			code:    "TEST_CODE",
			message: "test message",
			err:     baseErr,
			expected: &AppError{
				Code:    "TEST_CODE",
				Message: "test message",
				Err:     baseErr,
			},
		},
		{
			name:    "empty fields",
			code:    "",
			message: "",
			err:     nil,
			expected: &AppError{
				Code:    "",
				Message: "",
				Err:     nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := New(tt.code, tt.message, tt.err)
			require.Equal(t, tt.expected.Code, result.Code)
			require.Equal(t, tt.expected.Message, result.Message)
			require.Equal(t, tt.expected.Err, result.Err)
		})
	}
}

func TestPredefinedErrors(t *testing.T) {
	tests := []struct {
		name     string
		err      *AppError
		expected *AppError
	}{
		{
			name:     "ErrNotFound",
			err:      ErrNotFound,
			expected: New("NOT_FOUND", "not found", nil),
		},
		{
			name:     "ErrInvalidArgument",
			err:      ErrInvalidArgument,
			expected: New("INVALID_ARGUMENT", "invalid argument", nil),
		},
		{
			name:     "ErrInternal",
			err:      ErrInternal,
			expected: New("INTERNAL", "internal server error", nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected.Code, tt.err.Code)
			require.Equal(t, tt.expected.Message, tt.err.Message)
			require.Equal(t, tt.expected.Err, tt.err.Err)
		})
	}
}
