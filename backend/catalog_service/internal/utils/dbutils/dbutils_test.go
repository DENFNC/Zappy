package dbutils_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/DENFNC/Zappy/catalog_service/internal/utils/dbutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockScanner struct {
	mock.Mock
}

func (ms *MockScanner) Scan(dest ...any) error {
	args := ms.Called(dest...)
	return args.Error(0)
}

type Person struct {
	Name string
	Age  int
}

func TestScanStruct_MultipleCases(t *testing.T) {
	tests := []struct {
		name     string
		values   []any
		wantName string
		wantAge  int
	}{
		{"Alice-30", []any{"Alice", 30}, "Alice", 30},
		{"Bob-17", []any{"Bob", 17}, "Bob", 17},
		{"Empty-0", []any{"", 0}, "", 0},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockScanner := new(MockScanner)
			mockScanner.
				On("Scan", mock.Anything, mock.Anything).
				Return(nil).
				Run(func(args mock.Arguments) {
					*args.Get(0).(*string) = tc.values[0].(string)
					*args.Get(1).(*int) = tc.values[1].(int)
				})

			var p Person
			err := dbutils.ScanStruct(mockScanner, &p)
			assert.NoError(t, err)

			mockScanner.AssertExpectations(t)
			assert.Equal(t, tc.wantName, p.Name)
			assert.Equal(t, tc.wantAge, p.Age)
		})
	}
}

func TestScanStruct_ErrorCases(t *testing.T) {
	tests := []struct {
		name    string
		dest    any
		wantErr string
	}{
		{
			name:    "nil dest",
			dest:    nil,
			wantErr: "dest must be",
		},
		{
			name:    "not-pointer",
			dest:    Person{},
			wantErr: "dest must be",
		},
		{
			name:    "pointer to non-struct",
			dest:    new(int),
			wantErr: "dest must point to a struct",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockScanner := new(MockScanner)
			err := dbutils.ScanStruct(mockScanner, tc.dest)

			assert.ErrorContains(t, err, tc.wantErr)
			mockScanner.AssertNotCalled(t, "Scan")
		})
	}
}

func TestScanStruct_ScanError(t *testing.T) {
	mockErr := errors.New("underlying scan failure")
	mockScanner := new(MockScanner)
	mockScanner.
		On("Scan", mock.Anything, mock.Anything).
		Return(mockErr)

	var p Person
	err := dbutils.ScanStruct(mockScanner, &p)

	assert.Error(t, err)
	assert.ErrorContains(t, err, "scan into struct failed: underlying scan failure")
	assert.True(t, errors.Is(err, mockErr))

	mockScanner.AssertExpectations(t)
}

func TestScanStruct_UnexportedFieldsSkipped(t *testing.T) {
	type Mixed struct {
		Name  string
		age   int
		Email string
	}

	var m Mixed
	mockScanner := new(MockScanner)
	mockScanner.
		On("Scan", mock.Anything, mock.Anything).
		Return(nil).
		Run(func(args mock.Arguments) {
			require.Len(t, args, 2)
			dest0 := args.Get(0)
			dest1 := args.Get(1)
			require.Equal(t, reflect.TypeOf(dest0), reflect.TypeOf(new(string)))
			require.Equal(t, reflect.TypeOf(dest1), reflect.TypeOf(new(string)))
			*dest0.(*string) = "Ivan"
			*dest1.(*string) = "ivan@example.com"
		})

	err := dbutils.ScanStruct(mockScanner, &m)
	assert.NoError(t, err)

	mockScanner.AssertExpectations(t)
	assert.Equal(t, "Ivan", m.Name)
	assert.Equal(t, "ivan@example.com", m.Email)
	assert.Equal(t, 0, m.age)
}
