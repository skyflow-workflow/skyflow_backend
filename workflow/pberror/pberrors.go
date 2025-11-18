package pberror

import (
	"fmt"

	v1pb "github.com/skyflow-workflow/skyflow_backbend/api/v1"
)

type PBError struct {
	Code    v1pb.ErrorCode
	Message string
}

func (e *PBError) Error() string {
	return fmt.Sprintf("%s:%s", e.Code.String(), e.Message)
}

func NewPBError(code v1pb.ErrorCode, message string) *PBError {
	pberr := &PBError{
		Code:    code,
		Message: message,
	}
	return pberr
}

func IsPBError(err error) (*PBError, bool) {
	pberror, ok := err.(*PBError)
	return pberror, ok
}
