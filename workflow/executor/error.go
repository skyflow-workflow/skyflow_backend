package executor

/*
 * @Author: mumangtao@gmail.com
 * @Date: 2020-08-12 21:31:09
 * @Last Modified by: mumangtao@gmail.com
 * @Last Modified time: 2020-08-12 21:37:22
 */

import (
	"fmt"
	"slices"
)

var (
	// ErrorCreateExecutionUUIDFailed ...
	ErrorCreateExecutionUUIDFailed = fmt.Errorf("create execution new uuid  failed")
	// ErrorUnrecognizedEvent ...
	ErrorUnrecognizedEvent = fmt.Errorf("event type  unrecognized")

	// ErrorStepNotFound  常用错误类型 之一， step
	ErrorStepNotFound = fmt.Errorf("step not found ")

	// ErrorStepGroupNotFound  常用错误类型 之一， StepGroup没找到
	ErrorStepGroupNotFound = fmt.Errorf("stepgroup not found")

	// ErrorExecutionNotFound  常用错误类型 之一， Execution没找到
	ErrorExecutionNotFound = fmt.Errorf("execution not found")

	// ErrorExecutionStatus ...
	ErrorExecutionStatus = fmt.Errorf("execution status is incorrect")
	// ErrorNotAllStepGroupSucceed ...
	ErrorNotAllStepGroupSucceed = fmt.Errorf("not all stepgroup succeed")

	// ErrorMapBranchNumberExceedLimit Map分支数超过限额
	ErrorMapBranchNumberExceedLimit = fmt.Errorf("map branch number exceed limit")

	// ErrorOutputSizeExceedLimit Output超过限额
	ErrorOutputSizeExceedLimit = fmt.Errorf("output size exceed limit")

	// ErrorInputSizeExceedLimit Input超过限额
	ErrorInputSizeExceedLimit = fmt.Errorf("input size exceed limit")
	// ErrorUnrecognizedEventType
	ErrorUnrecognizedEventType = fmt.Errorf("unrecognized event type")
	// ErrorStepTypeIsNotMatch
	ErrorStepTypeIsNotMatch = fmt.Errorf("step type is not match")
)

// StandardErrorNames  name for handling in retry and catch
var StandardErrorNames = struct {
	StatesALL             string
	StatesRuntime         string
	StatesTimeout         string
	StatesHearbeatTimeout string
	StatesTaskFailed      string
	StatesPermissions     string
	StatesGroupFailed     string
}{
	StatesALL:             "States.ALL",
	StatesRuntime:         "States.Runtime",
	StatesTimeout:         "States.Timeout",
	StatesHearbeatTimeout: "States.HearbeatTimeout",
	StatesTaskFailed:      "States.TaskFailed",
	StatesPermissions:     "States.Permissions",
	StatesGroupFailed:     "States.GroupFailed", // 一个组内的节点失败了
}

var StandardErrorBuiltinNamesList = []string{
	StandardErrorNames.StatesRuntime,
	StandardErrorNames.StatesTimeout,
	StandardErrorNames.StatesHearbeatTimeout,
	StandardErrorNames.StatesTaskFailed,
	StandardErrorNames.StatesPermissions,
	StandardErrorNames.StatesGroupFailed,
}

// ExtendErrorNames 扩展错误匹配队列
func ExtendErrorNames(errorname string) []string {

	if slices.Contains(StandardErrorBuiltinNamesList, errorname) {
		return []string{errorname, StandardErrorNames.StatesALL}
	}
	return []string{errorname, StandardErrorNames.StatesTaskFailed, StandardErrorNames.StatesALL}
}
