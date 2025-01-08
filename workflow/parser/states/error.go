package states

import "errors"

type FieldError struct {
	// The error that occurred
	RawError error
	Line     int64
	Offset   int64
	Paths    []string
}

func (e *FieldError) Error() string {
	return e.RawError.Error()
}

func NewFieldError(err error, paths ...string) *FieldError {
	return &FieldError{
		RawError: err,
		Paths:    paths,
	}
}

var (
	ErrInvalidStateType = errors.New("invalid state type")
)
