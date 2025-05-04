package states

import (
	"errors"
	"fmt"
	"strings"
)

// ErrorLackOfRequiredField ...
var (
	ErrorLackOfRequiredField = errors.New("lack of required field")
	ErrorInvalidStateType    = errors.New("invalid state type")
	ErrorInvalidFiledContent = errors.New("field content is invalid")
	ErrorFiledDenied         = errors.New("field is dentied")
	ErrorFiledRequired       = errors.New("field is required")
	ErrorInvalidData         = errors.New("invalid data")
	ErrorInvalidField        = errors.New("invalid field")
)

// FieldError is an error that occurred
// in a field at a specific path or line number and column number or offset in the file.
// FieldError is final state error for FieldPathError.
type FieldError struct {
	// The error that occurred
	RawError error
	Line     int64
	Column   int64
	Offset   int64
	Paths    []string
}

// Error ...
func (e *FieldError) Error() string {
	msg := fmt.Sprintf("%s, path: %s", e.RawError.Error(), strings.Join(e.Paths, "."))
	return msg
}

// FiledPathError FieldPathError is an error that occurred in a field at a specific path.
type FiledPathError struct {
	// The error that occurred
	RawError error
	Paths    []string
}

// Error ...
func (e *FiledPathError) Error() string {
	msg := fmt.Sprintf("%s, path: %s", e.RawError.Error(), strings.Join(e.Paths, "."))
	return msg
}

// NewFieldPathError ...
func NewFieldPathError(err error, paths ...string) *FieldError {
	return &FieldError{
		RawError: err,
		Paths:    paths,
	}
}
