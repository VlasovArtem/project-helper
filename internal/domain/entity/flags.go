package entity

import (
	"github.com/pkg/errors"
)

var errOperationNotProvided = errors.New("operation not provided")

type Flags struct {
	Operation    string
	DynamicFlags map[string]any
}

func NewFlags() *Flags {
	return &Flags{
		DynamicFlags: make(map[string]any),
	}
}

func (f *Flags) Validate() error {
	if f.Operation == "" {
		return errOperationNotProvided
	}

	return nil
}

func (f *Flags) GetFlagStringValue(flag string) string {
	value, _ := f.GetRequiredFlagStringValue(flag)

	return value
}

func (f *Flags) GetRequiredFlagStringValue(flag string) (string, error) {
	value, ok := f.DynamicFlags[flag]

	if !ok {
		return "", errors.Errorf("flag %s not found", flag)
	}
	if value == nil {
		return "", errors.Errorf("flag %s is nil", flag)
	}

	assertedValue, ok := value.(*string)
	if !ok {
		return "", errors.Errorf("flag %s is not a string", flag)
	}

	return *assertedValue, nil
}

func (f *Flags) GetRequiredFlagArrayValue(flag string) ([]string, error) {
	value, ok := f.DynamicFlags[flag]

	if !ok {
		return nil, errors.Errorf("flag %s not found", flag)
	}
	if value == nil {
		return nil, errors.Errorf("flag %s is nil", flag)
	}

	assertedValue, ok := value.(*[]string)
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
