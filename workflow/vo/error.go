package vo

import "fmt"

var (
	// ErrorExecutionUUIDExisted  常用错误类型 之一， uuid已经存在
	ErrorExecutionUUIDExisted = fmt.Errorf("execution uuid  has been existed")
	// ErrorTaskTokenNotFound  常用错误类型 之一， TaskToken没找到
	ErrorTaskTokenNotFound = fmt.Errorf("TaskTokenNotFound")
	// ErrorStepNotFound  常用错误类型 之一， Step没找到
	ErrorStepNotFound = fmt.Errorf("StepNotFound ")
	// ErrorStepGroupNotFound  常用错误类型 之一， StepGroup没找到
	ErrorStepGroupNotFound = fmt.Errorf("StepgroupNotFound ")
	// ErrorExecutionNotFound  常用错误类型 之一， Execution没找到
	ErrorExecutionNotFound = fmt.Errorf("ExecutionNotFound ")
	ErrorUUIDExisted       = fmt.Errorf("uuid  has been existed")
	ErrorGenUUIDFailed     = fmt.Errorf("generate UUID failed")
	ErrUnrecognizeEvent    = fmt.Errorf("event type  unrecognize ")

	ErrorStepStatus           = fmt.Errorf("StepStatusError")
	ErrorExecutionStatus      = fmt.Errorf("execution status is incorrect")
	ErrorNotAllStepGoupFinish = fmt.Errorf("not all step group succeed")
	// activity not found error
	ErrorActivityTaskNotFound = fmt.Errorf("ActivityTaskNotFound")
	//unsupport operation for step
	ErrorUnsupportOperationForStep = fmt.Errorf("unsupport operation for step")

	ErrorParamterInvalid             = fmt.Errorf("parameter is invalild")
	ErrorParseWorkflow               = fmt.Errorf("parse workflow failed")
	ErrorUnrecognizeStatemachineType = fmt.Errorf("unrecognize statemachine type")

	// ErrorOutputSizeLimitExceeded Output超过限额
	ErrorOutputSizeLimitExceeded = fmt.Errorf("output size limit exceeded")

	// ErrorInputSizeLimitExceeded Input超过限额
	ErrorInputSizeLimitExceeded = fmt.Errorf("input size limit exceeded")

	// ErrorTaskInputSizeLimitExceeded Input超过限额
	ErrorTaskInputSizeLimitExceeded = fmt.Errorf("task input size limit exceeded")

	// ErrorParameterExceedLimit 请求参数超过限额
	ErrorParamterLimitExceeded = fmt.Errorf("parameter limit exceeded")

	// ErrorStartExecutionInputSizeLimitExceeded StartExecution输入超过限额, 基于ErrorParamterExceedLimit
	ErrorStartExecutionInputSizeLimitExceeded = fmt.Errorf("%w: input size limit exceeded", ErrorParamterLimitExceeded)

	// ErrorWorkflowSizeLimitExceeded Workflow超过限额, 基于ErrorParamterExceedLimit
	ErrorWorkflowSizeLimitExceeded = fmt.Errorf("%w:workflow size limit exceeded", ErrorParamterLimitExceeded)

	// ErrorMapBranchNumberLimitExceeded Map分支数超过限额, 基于ErrorParamterExceedLimit
	ErrorMapBranchNumberLimitExceeded = fmt.Errorf("%w:map branch number limit exceeded", ErrorParamterLimitExceeded)

	// ErrorTitleSizeLimitExceeded Title超过限额, 基于ErrorParamterExceedLimit
	ErrorTitleSizeLimitExceeded = fmt.Errorf("%w:title size limit exceeded", ErrorParamterLimitExceeded)

	// ErrorDescriptionSizeExceedLimit uuid 超过限额, 基于ErrorParamterExceedLimit
	ErrorExecutionUUIDSizeLimitExceed = fmt.Errorf("%w:execution uuid size limit exceeded", ErrorParamterLimitExceeded)

	// workflow definition
	// ErrorWorkflowDefinitionInvalid Workflow定义不合法
	ErrorWorkflowDefinitionInvalid = fmt.Errorf("workflow definition is invalid")
	// ErrorMapConncurrencyLimitExceeded Map并发数超过限额, 基于ErrorParamterExceedLimit
	ErrorMapConncurrencyLimitExceeded = fmt.Errorf("%w:map conncurrency limit exceeded", ErrorWorkflowDefinitionInvalid)

	// ErrorParallelBranchNumberLimitExceeded 并行分支数超过限额, 基于ErrorParamterExceedLimit
	ErrorParallelBranchNumberLimitExceeded = fmt.Errorf("%w:parallel branch number limit exceeded", ErrorWorkflowDefinitionInvalid)

	// ErrorUnspportedOperationForExecution 不支持的操作
	ErrorUnspportedOperationForExecution = fmt.Errorf("unsupport operation for execution")

	ErrorAKSKInvalid = fmt.Errorf("aksk is invalid")
)
