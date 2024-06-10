package entity

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"project-helper/internal/utils"
)

func TestDynamicFlagValueGetString(t *testing.T) {
	//func GetString(d *DynamicFlagValue) (string, error) {
	//	if d == nil {
	//		return "", domainerrors.ErrorNilInput
	//	}
	//
	//	switch d.Type {
	//	case String:
	//		value, ok := d.Value.(*string)
	//		if !ok {
	//			return "", errors.Wrap(errInvalidFlagTypeValue, "flag is not a string")
	//		}
	//
	//		return *value, nil
	//	case Array:
	//		value, ok := d.Value.(*[]string)
	//		if !ok {
	//			return "", errors.Wrap(errInvalidFlagTypeValue, "flag is not an array")
	//		}
	//
	//		if value == nil {
	//			return "", errors.Wrap(errInvalidFlagTypeValue, "flag is nil")
	//		}
	//
	//		if len(*value) == 0 {
	//			return "", errors.Wrap(errInvalidFlagTypeValue, "flag is empty")
	//		}
	//
	//		return strings.Join(*value, ","), nil
	//	default:
	//		return "", errors.Wrapf(errInvalidFlagTypeValue, "flag type '%s' is not supported", d.Type)
	//	}
	//}

	t.Parallel()

	tests := map[string]struct {
		dynamicFlagValue *DynamicFlagValue
		expectedValue    string
		expectedError    error
	}{
		"success string": {
			dynamicFlagValue: &DynamicFlagValue{
				Name:  "flag",
				Type:  String,
				Value: utils.MakePointer("value"),
			},
			expectedValue: "value",
		},
		"success array": {
			dynamicFlagValue: &DynamicFlagValue{
				Name:  "flag",
				Type:  Array,
				Value: utils.MakePointer([]string{"value1", "value2"}),
			},
			expectedValue: "value1,value2",
		},
		"nil": {
			dynamicFlagValue: nil,
			expectedError:    errors.New("nil input"),
		},
		"not a string": {
			dynamicFlagValue: &DynamicFlagValue{
				Name:  "flag",
				Type:  String,
				Value: utils.MakePointer(42),
			},
			expectedError: errors.New("flag is not a string: invalid flag type value"),
		},
		"not an array": {
			dynamicFlagValue: &DynamicFlagValue{
				Name:  "flag",
				Type:  Array,
				Value: utils.MakePointer("value"),
			},
			expectedError: errors.New("flag is not an array: invalid flag type value"),
		},
		"array nil": {
			dynamicFlagValue: &DynamicFlagValue{
				Name:  "flag",
				Type:  Array,
				Value: nil,
			},
			expectedError: errors.New("flag is not an array: invalid flag type value"),
		},
		"array empty": {
			dynamicFlagValue: &DynamicFlagValue{
				Name:  "flag",
				Type:  Array,
				Value: utils.MakePointer([]string{}),
			},
			expectedError: errors.New("flag is empty: invalid flag type value"),
		},
		"unsupported type": {
			dynamicFlagValue: &DynamicFlagValue{
				Name:  "flag",
				Type:  "unsupported",
				Value: utils.MakePointer("value"),
			},
			expectedError: errors.New("flag type 'unsupported' is not supported: invalid flag type value"),
		},
	}

	for name, testCase := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			value, err := GetString(testCase.dynamicFlagValue)

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
				DynamicFlags: map[string]*DynamicFlagValue{
					"flag": {
						Name:  "flag",
						Type:  String,
						Value: utils.MakePointer("value"),
					},
				},
			},
			expectedValue: "value",
		},
		"without value": {
			runFlags: &Flags{DynamicFlags: map[string]*DynamicFlagValue{}},
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
				DynamicFlags: map[string]*DynamicFlagValue{
					"flag": {
						Name:  "flag",
						Type:  String,
						Value: utils.MakePointer("value"),
					},
				},
			},
			expectedValue: "value",
		},
		"without value": {
			runFlags:      &Flags{DynamicFlags: map[string]*DynamicFlagValue{}},
			expectedError: errors.New("flag flag not found"),
		},
		"nil value": {
			runFlags: &Flags{
				DynamicFlags: map[string]*DynamicFlagValue{
					"flag": nil,
				},
			},
			expectedError: errors.New("flag flag is nil"),
		},
		"not a string": {
			runFlags: &Flags{
				DynamicFlags: map[string]*DynamicFlagValue{
					"flag": {
						Name:  "flag",
						Type:  String,
						Value: utils.MakePointer(42),
					},
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
			flags: &Flags{
				DynamicFlags: map[string]*DynamicFlagValue{
					"flag": {
						Name:  "flag",
						Type:  Array,
						Value: utils.MakePointer([]string{"value1", "value2"}),
					},
				},
			},
			expectedValue: []string{"value1", "value2"},
		},
		"without value": {
			flags:         &Flags{DynamicFlags: map[string]*DynamicFlagValue{}},
			expectedError: errors.New("flag flag not found"),
		},
		"nil value": {
			flags: &Flags{
				DynamicFlags: map[string]*DynamicFlagValue{
					"flag": nil,
				},
			},
			expectedError: errors.New("flag flag is nil"),
		},
		"not an array": {
			flags: &Flags{
				DynamicFlags: map[string]*DynamicFlagValue{
					"flag": {
						Name:  "flag",
						Type:  Array,
						Value: utils.MakePointer("value"),
					},
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
			flags: &Flags{
				DynamicFlags: map[string]*DynamicFlagValue{
					"flag": {
						Name:  "flag",
						Type:  Array,
						Value: &[]string{"value1", "value2"},
					},
				},
			},
			expectedValue: []string{"value1", "value2"},
		},
		"without value": {
			flags:         &Flags{DynamicFlags: map[string]*DynamicFlagValue{}},
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
