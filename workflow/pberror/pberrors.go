package pberror

import (
	"fmt"

	v1pb "github.com/skyflow-workflow/skyflow_backbend/api/v1"
	"trpc.group/trpc-go/trpc-go/errs"
)

type PBError struct {
	Code    v1pb.ErrorCode
	Message string
}

func (e *PBError) Error() string {
	return fmt.Sprintf("%s:%s", e.Code.String(), e.Message)
}

func NewPBError(code v1pb.ErrorCode, message string) error {
	pberr := &PBError{
		Code:    code,
		Message: message,
	}
	return errs.New(int32(code), pberr.Error())
}
