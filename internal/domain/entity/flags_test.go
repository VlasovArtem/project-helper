package entity

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"project-helper/internal/utils"
)

func TestNewFlags(t *testing.T) {
	t.Parallel()

	flags := NewFlags()

	assert.NotNil(t, flags)
	assert.Empty(t, flags.DynamicFlags)
	assert.Empty(t, flags.Operation)
}

func TestFlagsValidate(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		flags         *Flags
		expectedError error
	}{
		"success": {
			flags: &Flags{
				Operation: "operation",
			},
		},
		"error operation not provided": {
			flags:         &Flags{},
			expectedError: errors.New("operation not provided"),
		},
	}

	for name, testCase := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := testCase.flags.Validate()

			if testCase.expectedError != nil {
				require.Error(t, err)
				assert.ErrorContains(t, err, testCase.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRunFlagsGetFlagStringValue(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		runFlags      *Flags
		expectedValue string
	}{
		"with value": {
			runFlags: &Flags{
				DynamicFlags: map[string]any{
					"flag": utils.MakePointer("value"),
				},
			},
			expectedValue: "value",
		},
		"without value": {
			runFlags: &Flags{DynamicFlags: map[string]any{}},
		},
	}

	for name, testCase := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			value := testCase.runFlags.GetFlagStringValue("flag")

			assert.Equal(t, testCase.expectedValue, value)
		})
	}
}

func TestGetRequiredFlagStringValue(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		runFlags      *Flags
		expectedValue string
		expectedError error
	}{
		"with value": {
			runFlags: &Flags{
				DynamicFlags: map[string]any{
					"flag": utils.MakePointer("value"),
				},
			},
			expectedValue: "value",
		},
		"without value": {
			runFlags:      &Flags{DynamicFlags: map[string]any{}},
			expectedError: errors.New("flag flag not found"),
		},
		"nil value": {
			runFlags: &Flags{
				DynamicFlags: map[string]any{
					"flag": nil,
				},
			},
			expectedError: errors.New("flag flag is nil"),
		},
		"not a string": {
			runFlags: &Flags{
				DynamicFlags: map[string]any{
					"flag": utils.MakePointer(42),
				},
			},
			expectedError: errors.New("flag flag is not a string"),
		},
	}

	for name, testCase := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			value, err := testCase.runFlags.GetRequiredFlagStringValue("flag")

			if testCase.expectedError != nil {
				require.Error(t, err)
				assert.ErrorContains(t, err, testCase.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expectedValue, value)
			}
		})
	}
}

func TestGetRequiredFlagArrayValue(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		flags         *Flags
		expectedValue []string
		expectedError error
	}{
		"with value": {
			flags: &Flags{DynamicFlags: map[string]any{
				"flag": utils.MakePointer([]string{"value1", "value2"}),
			},
			},
			expectedValue: []string{"value1", "value2"},
		},
		"without value": {
			flags:         &Flags{DynamicFlags: map[string]any{}},
			expectedError: errors.New("flag flag not found"),
		},
		"nil value": {
			flags: &Flags{DynamicFlags: map[string]any{
				"flag": nil,
			},
			},
			expectedError: errors.New("flag flag is nil"),
		},
		"not an array": {
			flags: &Flags{
				DynamicFlags: map[string]any{
					"flag": utils.MakePointer("value"),
				},
			},
			expectedError: errors.New("flag flag is not an array"),
		},
	}

	for name, testCase := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			value, err := testCase.flags.GetRequiredFlagArrayValue("flag")

			if testCase.expectedError != nil {
				require.Error(t, err)
				assert.ErrorContains(t, err, testCase.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expectedValue, value)
			}
		})
	}
}

func TestGetFlagArrayValue(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		flags         *Flags
		expectedValue []string
	}{
		"with value": {
			flags: &Flags{DynamicFlags: map[string]any{
				"flag": utils.MakePointer([]string{"value1", "value2"}),
			},
			},
			expectedValue: []string{"value1", "value2"},
		},
		"without value": {
			flags:         &Flags{DynamicFlags: map[string]any{}},
			expectedValue: make([]string, 0),
		},
	}

	for name, testCase := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			value := testCase.flags.GetFlagArrayValue("flag")

			assert.Equal(t, testCase.expectedValue, value)
		})
	}
}
