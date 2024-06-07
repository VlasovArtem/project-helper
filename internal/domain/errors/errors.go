package errors

import "errors"

var (
	ErrorOperationNotFound      = errors.New("operation not found")
	ErrorTagValueNotFound       = errors.New("tag value not found")
	ErrorPredefinedArgNotFound  = errors.New("predefined arg not found")
	ErrorPatternTagTypeNotFound = errors.New("pattern tag type not found")
	ErrorNilInput               = errors.New("nil input")
	ErrorAdditionalArgNotFound  = errors.New("additional arg not found")
	ErrorObjectIsNil            = errors.New("object is nil")
)
