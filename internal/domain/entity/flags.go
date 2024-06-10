package entity

import (
	"strings"

	"github.com/pkg/errors"
	domainerrors "project-helper/internal/domain/errors"
)

var errOperationNotProvided = errors.New("operation not provided")
var errInvalidFlagTypeValue = errors.New("invalid flag type value")

type DynamicFlagValue struct {
	Name  string
	Type  Type
	Value any
}

func GetString(d *DynamicFlagValue) (string, error) {
	if d == nil {
		return "", domainerrors.ErrorNilInput
	}

	switch d.Type {
	case String:
		value, ok := d.Value.(*string)
		if !ok {
			return "", errors.Wrap(errInvalidFlagTypeValue, "flag is not a string")
		}

		return *value, nil
	case Array:
		value, ok := d.Value.(*[]string)
		if !ok {
			return "", errors.Wrap(errInvalidFlagTypeValue, "flag is not an array")
		}

		if value == nil {
			return "", errors.Wrap(errInvalidFlagTypeValue, "flag is nil")
		}

		if len(*value) == 0 {
			return "", errors.Wrap(errInvalidFlagTypeValue, "flag is empty")
		}

		return strings.Join(*value, ","), nil
	default:
		return "", errors.Wrapf(errInvalidFlagTypeValue, "flag type '%s' is not supported", d.Type)
	}
}

type Flags struct {
	Operation    string
	DynamicFlags map[string]*DynamicFlagValue
}

func NewFlags() *Flags {
	return &Flags{
		DynamicFlags: make(map[string]*DynamicFlagValue),
	}
}

func (f *Flags) Validate() error {
	if f.Operation == "" {
		return errOperationNotProvided
	}

	return nil
}

func (f *Flags) GetFlag(flag string) (*DynamicFlagValue, error) {
	dynamicFlagValue, ok := f.DynamicFlags[flag]

	if !ok {
		return &DynamicFlagValue{}, errors.Errorf("flag %s not found", flag)
	}

	return dynamicFlagValue, nil
}

func (f *Flags) GetFlagStringValue(flag string) string {
	value, _ := f.GetRequiredFlagStringValue(flag)

	return value
}

func (f *Flags) GetRequiredFlagStringValue(flag string) (string, error) {
	dynamicFlagValue, ok := f.DynamicFlags[flag]

	if !ok {
		return "", errors.Errorf("flag %s not found", flag)
	}
	if dynamicFlagValue == nil {
		return "", errors.Errorf("flag %s is nil", flag)
	}

	value, ok := dynamicFlagValue.Value.(*string)
	if !ok {
		return "", errors.Errorf("flag %s is not a string", flag)
	}

	return *value, nil
}

func (f *Flags) GetRequiredFlagArrayValue(flag string) ([]string, error) {
	dynamicFlagValue, ok := f.DynamicFlags[flag]

	if !ok {
		return nil, errors.Errorf("flag %s not found", flag)
	}
	if dynamicFlagValue == nil {
		return nil, errors.Errorf("flag %s is nil", flag)
	}

	assertedValue, ok := dynamicFlagValue.Value.(*[]string)
	if !ok {
		return nil, errors.Errorf("flag %s is not an array", flag)
	}

	if assertedValue == nil {
		return nil, errors.Errorf("flag %s is nil", flag)
	}

	if len(*assertedValue) == 0 {
		return nil, errors.Errorf("flag %s is empty", flag)
	}

	return *assertedValue, nil
}

func (f *Flags) GetFlagArrayValue(flag string) []string {
	value, err := f.GetRequiredFlagArrayValue(flag)
	if err != nil {
		return make([]string, 0)
	}

	return value
}
