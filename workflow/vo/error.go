package vo

import "fmt"

// ErrorExecutionUUIDExisted ...
var (
	// ErrorExecutionUUIDExisted  常用错误类型 之一， uuid已经存在
	ErrorExecutionUUIDExisted = fmt.Errorf("execution uuid  has been existed")
	// ErrorTaskTokenNotFound  常用错误类型 之一， TaskToken没找到
	ErrorTaskTokenNotFound = fmt.Errorf("TaskTokenNotFound")
	// ErrorStepNotFound  常用错误类型 之一， Step没找到
	ErrorStepNotFound = fmt.Errorf("step not found")
	// ErrorStepGroupNotFound  常用错误类型 之一， StepGroup没找到
	ErrorStepGroupNotFound = fmt.Errorf("stepgroup not found")
	// ErrorUserStepDataNotFound 常用错误类型 之一， UserStepData没找到
	ErrorUserStepDataNotFound = fmt.Errorf("user step data not found")
	// ErrorExecutionNotFound  常用错误类型 之一， Execution没找到
	ErrorExecutionNotFound     = fmt.Errorf("ExecutionNotFound ")
	ErrorUUIDExisted           = fmt.Errorf("uuid  has been existed")
	ErrorGenUUIDFailed         = fmt.Errorf("generate UUID failed")
	ErrorUnrecognizedEventType = fmt.Errorf("event type  unrecognized ")

	ErrorStepStatus            = fmt.Errorf("StepStatusError")
	ErrorExecutionStatus       = fmt.Errorf("execution status is incorrect")
	ErrorNotAllStepGroupFinish = fmt.Errorf("not all step group succeed")
	// activity not found error
	ErrorActivityTaskNotFound = fmt.Errorf("ActivityTaskNotFound")
	//unsupported operation for step
	ErrorUnsupportedOperationForStep = fmt.Errorf("unsupported operation for step")

	ErrorParameterInvalid             = fmt.Errorf("parameter is invalid")
	ErrorParseWorkflow                = fmt.Errorf("parse statemachine failed")
	ErrorUnrecognizedStatemachineType = fmt.Errorf("unrecognized statemachine type")
	ErrorUnsupportedStateType         = fmt.Errorf("unsupported state type")

	// ErrorOutputSizeLimitExceeded Output超过限额
	ErrorOutputSizeLimitExceeded = fmt.Errorf("output size limit exceeded")

	// ErrorInputSizeLimitExceeded Input超过限额
	ErrorInputSizeLimitExceeded = fmt.Errorf("input size limit exceeded")

	// ErrorTaskInputSizeLimitExceeded Input超过限额
	ErrorTaskInputSizeLimitExceeded = fmt.Errorf("task input size limit exceeded")

	// ErrorParameterExceedLimit 请求参数超过限额
	ErrorParameterExceedLimit = fmt.Errorf("parameter limit exceeded")

	// ErrorStartExecutionInputSizeLimitExceeded StartExecution输入超过限额, 基于ErrorParameterExceedLimit
	ErrorStartExecutionInputSizeLimitExceeded = fmt.Errorf("%w: input size limit exceeded", ErrorParameterExceedLimit)

	// ErrorStateMachineSizeLimitExceeded StateMachine超过限额, 基于ErrorParameterExceedLimit
	ErrorStateMachineSizeLimitExceeded = fmt.Errorf("%w:statemachine size limit exceeded", ErrorParameterExceedLimit)

	// ErrorMapBranchNumberLimitExceeded Map分支数超过限额, 基于ErrorParameterExceedLimit
	ErrorMapBranchNumberLimitExceeded = fmt.Errorf("%w:map branch number limit exceeded", ErrorParameterExceedLimit)

	// ErrorTitleSizeLimitExceeded Title超过限额, 基于ErrorParameterExceedLimit
	ErrorTitleSizeLimitExceeded = fmt.Errorf("%w:title size limit exceeded", ErrorParameterExceedLimit)

	// ErrorExecutionUUIDSizeLimitExceed uuid 超过限额, 基于ErrorParameterExceedLimit
	ErrorExecutionUUIDSizeLimitExceed = fmt.Errorf("%w:execution uuid size limit exceeded", ErrorParameterExceedLimit)

	// statemachine definition
	// ErrorStateMachineDefinitionInvalid Workflow定义不合法
	ErrorStateMachineDefinitionInvalid = fmt.Errorf("statemachine definition is invalid")
	// ErrorMapConncurrencyLimitExceeded Map并发数超过限额, 基于ErrorParameterExceedLimit
	ErrorMapConncurrencyLimitExceeded = fmt.Errorf("%w:map conncurrency limit exceeded", ErrorStateMachineDefinitionInvalid)

	// ErrorParallelBranchNumberLimitExceeded 并行分支数超过限额, 基于ErrorParameterExceedLimit
	ErrorParallelBranchNumberLimitExceeded = fmt.Errorf("%w:parallel branch number limit exceeded", ErrorStateMachineDefinitionInvalid)

	// ErrorUnsupportedOperationForExecution 不支持的操作
	ErrorUnsupportedOperationForExecution = fmt.Errorf("unsupported operation for execution")

	// ErrorAKSKInvalid AKSK 不合法
	ErrorAKSKInvalid      = fmt.Errorf("aksk is invalid")
	ErrorParameterMissing = fmt.Errorf("parameter is missing")
	ErrorParamterTooLong  = fmt.Errorf("parameter is too long")
	// ErrorNonsupportOperationForStep 不支持的操作
	ErrorNonsupportOperationForStep = fmt.Errorf("nonsupport operation for step")
)
