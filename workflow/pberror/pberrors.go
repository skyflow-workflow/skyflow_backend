package pberror

import (
	"fmt"

	pbv1 "github.com/skyflow-workflow/skyflow_backend/api/v1"
)

type PBError struct {
	Code     pbv1.ErrorCode
	Message  string
	RawError error
}

func (e *PBError) Error() string {
	return fmt.Sprintf("%s:%s", e.Code.String(), e.Message)
}

func NewPBError(code pbv1.ErrorCode, message string) *PBError {
	pbErr := &PBError{
		Code:    code,
		Message: message,
	}
	return pbErr
}

func (err *PBError) WithMesage(message string) *PBError {
	err.Message = message
	return err
}

func NewFromError(code pbv1.ErrorCode, err error) *PBError {
	pbErr := &PBError{
		Code:     code,
		Message:  err.Error(),
		RawError: err,
	}
	return pbErr
}

func IsPBError(err error) (*PBError, bool) {
	pbError, ok := err.(*PBError)
	return pbError, ok
}
