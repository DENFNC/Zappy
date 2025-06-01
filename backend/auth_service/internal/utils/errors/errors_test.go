package errpkg_test

import (
	"errors"
	"fmt"
	"testing"

	errpkg "github.com/DENFNC/Zappy/auth_service/internal/utils/errors"
	"github.com/stretchr/testify/assert"
)

func TestNewAppError_BasicFields(t *testing.T) {
	excuse := errors.New("inner")
	ae := errpkg.New("CODE_X", "message text", excuse)

	assert.Equal(t, "CODE_X", ae.Code)
	assert.Equal(t, "message text", ae.Message)
	assert.Same(t, excuse, ae.Err)
}

func TestAppError_ErrorReturnsMessage(t *testing.T) {
	ae := errpkg.New("C", "some error occurred", nil)
	assert.Equal(t, "some error occurred", ae.Error())
}

func TestAppError_Unwrap(t *testing.T) {
	inner := errors.New("root cause")
	ae := errpkg.New("C2", "wrapped", inner)
	assert.Same(t, inner, ae.Unwrap())
}

func TestAppError_ErrorsIsAndAs(t *testing.T) {
	inner := errors.New("root")
	a := errpkg.New("E1", "foo", inner)
	err := fmt.Errorf("context: %w", a)

	assert.True(t, errors.Is(err, a))
	assert.True(t, errors.Is(err, inner))

	var appErr *errpkg.AppError
	assert.True(t, errors.As(err, &appErr))
	assert.Equal(t, "foo", appErr.Message)
}

func TestPredefinedErrors_Constants(t *testing.T) {
	cases := []struct {
		name    string
		e       *errpkg.AppError
		wantC   string
		wantMsg string
	}{
		{"NotFound", errpkg.ErrNotFound, "NOT_FOUND", "not found"},
		{"InvalidArgument", errpkg.ErrInvalidArgument, "INVALID_ARGUMENT", "invalid argument"},
		{"Internal", errpkg.ErrInternal, "INTERNAL", "internal server error"},
		{"Constraint", errpkg.ErrConstraint, "ERR_FOREIGN_KEY_VIOLATION", "external switch violation"},
		{"UniqueViolation", errpkg.ErrUniqueViolation, "ERR_UNIQUE_VIOLATION", "violation of uniqueness"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.wantC, tc.e.Code)
			assert.Equal(t, tc.wantMsg, tc.e.Message)
			assert.Nil(t, tc.e.Err)
		})
	}
}
