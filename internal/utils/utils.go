package utils

import (
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

var Validate = validator.New(validator.WithRequiredStructEnabled())

func MakePointer[v any](s v) *v {
	return &s
}

func EscapeValue(value string) (string, error) {
	unquote, err := strconv.Unquote(value)
	if err != nil && errors.Is(err, strconv.ErrSyntax) {
		return value, nil
	} else if err != nil {
		return "", errors.Wrap(err, "failed to unquote value")
	}

	return unquote, nil
}
