package pberror

import (
	"fmt"
)

type PBError struct {
	Code    int32
	Message string
}

func (e *PBError) Error() string {
	return fmt.Sprintf("code: %d, message: %s", e.Code, e.Message)
}

func NewPBError(code int32, message string) *PBError {
	return &PBError{
		Code:    code,
		Message: message,
	}
}
