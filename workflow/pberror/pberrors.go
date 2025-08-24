package pberror

import (
	"fmt"

	"github.com/skyflow-workflow/skyflow_backbend/gen/pb"
	"trpc.group/trpc-go/trpc-go/errs"
)

type PBError struct {
	Code    pb.ErrorCode
	Message string
}

func (e *PBError) Error() string {
	return fmt.Sprintf("%s:%s", e.Code.String(), e.Message)
}

func NewPBError(code pb.ErrorCode, message string) error {
	pberr := &PBError{
		Code:    code,
		Message: message,
	}
	return errs.New(int32(code), pberr.Error())
}
