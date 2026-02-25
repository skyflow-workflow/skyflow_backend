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
	ErrorFiledDenied         = errors.New("field is denied")
	ErrorFiledRequired       = errors.New("field is required")
	ErrorInvalidData         = errors.New("invalid data")
	ErrorInvalidField        = errors.New("invalid field")
)

// FieldPathError is an error that occurred
// in a field at a specific path or line number and column number or offset in the file.
// FieldPathError is final state error for FieldPathError.
type FieldPathError struct {
	// The error that occurred
	RawError error
	Line     int64
	Column   int64
	Offset   int64
	Paths    []string
	// state name
	StateName string
	// state type
	StateType string
}

// Error string format: error message, path: path1.path2.path3
func (e *FieldPathError) Error() string {
	msg := fmt.Sprintf(
		"error: %s, state: %s, type: %s, path: %s",
		e.RawError.Error(), e.StateName, e.StateType, strings.Join(e.Paths, "."))
	return msg
}

// NewFieldPathError NewFieldPathError is a constructor for FieldPathError
// err is the error that occurred
// stateType is the type of the state
// stateName is the name of the state
// paths is the paths of the error
func NewFieldPathError(err error, stateType string, paths ...string) *FieldPathError {
	return &FieldPathError{
		RawError:  err,
		Paths:     paths,
		StateType: stateType,
	}
}
