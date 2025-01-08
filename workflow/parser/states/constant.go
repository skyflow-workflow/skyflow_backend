package states

/*
# @Date    : Thu Mar 14 2019
# @Author  : mumangtao(mumangtao@gmail.com)
# @Link    : link
# @Version : 1.0.0
# @Function: file for
*/

import (
	"github.com/go-playground/validator/v10"
)

var myvalidate = validator.New()

// 						Pass		Task		Choice		Wait		Succeed		Fail		Parallel    Map
// Type					Required	Required	Required	Required	Required	Required	Required	Required
// Comment				Allowed		Allowed		Allowed		Allowed		Allowed		Allowed		Allowed		Allowed
// InputPath/OutputPath	Allowed		Allowed		Allowed		Allowed		Allowed					Allowed		Allowed
// Parameters			Allowed		Allowed														Allowed		Allowed
// ResultPath			Allowed		Allowed														Allowed		Allowed
// One of: Next/End		Required	Required				Required							Required	Required
// Retry, Catch						Allowed														Allowed  	Allowed

// StateType constant State type define here
var StateType = struct {
	Task     string
	Choice   string
	Fail     string
	Succeed  string
	Wait     string
	Pass     string
	Parallel string
	Map      string
	Suspend  string
	// 子流程类型, Parallel/Map 生成的中间步骤
	StateGroup string
}{
	Task:       "Task",
	Choice:     "Choice",
	Fail:       "Fail",
	Succeed:    "Succeed",
	Wait:       "Wait",
	Pass:       "Pass",
	Parallel:   "Parallel",
	Map:        "Map",
	Suspend:    "Suspend",
	StateGroup: "StateGroup",
}

var requiredlevel = struct {
	Required int
	Allowed  int
	Deny     int
}{
	Required: 0,
	Allowed:  1,
	Deny:     2,
}

// ElementRequiredMap  元素依赖需求字段映射
var ElementRequiredMap = map[string]map[string]int{
	StateType.Pass: {
		"Comment":    requiredlevel.Allowed,
		"InputPath":  requiredlevel.Allowed,
		"OutputPath": requiredlevel.Allowed,
		"Parameters": requiredlevel.Allowed,
		"ResultPath": requiredlevel.Allowed,
		"NextEnd":    requiredlevel.Required,
	},
	StateType.Task: {
		"Comment":    requiredlevel.Allowed,
		"InputPath":  requiredlevel.Allowed,
		"OutputPath": requiredlevel.Allowed,
		"Parameters": requiredlevel.Allowed,
		"ResultPath": requiredlevel.Allowed,
		"NextEnd":    requiredlevel.Required,
	},
	StateType.Choice: {
		"Comment":    requiredlevel.Allowed,
		"InputPath":  requiredlevel.Allowed,
		"OutputPath": requiredlevel.Allowed,
		"Parameters": requiredlevel.Deny,
		"ResultPath": requiredlevel.Deny,
		"NextEnd":    requiredlevel.Deny,
	},
	StateType.Wait: {
		"Comment":    requiredlevel.Allowed,
		"InputPath":  requiredlevel.Allowed,
		"OutputPath": requiredlevel.Allowed,
		"Parameters": requiredlevel.Deny,
		"ResultPath": requiredlevel.Deny,
		"NextEnd":    requiredlevel.Required,
	},
	StateType.Succeed: {
		"Comment":    requiredlevel.Allowed,
		"InputPath":  requiredlevel.Allowed,
		"OutputPath": requiredlevel.Allowed,
		"Parameters": requiredlevel.Deny,
		"ResultPath": requiredlevel.Deny,
		"NextEnd":    requiredlevel.Deny,
	},
	StateType.Fail: {
		"Comment":    requiredlevel.Allowed,
		"InputPath":  requiredlevel.Deny,
		"OutputPath": requiredlevel.Deny,
		"Parameters": requiredlevel.Deny,
		"ResultPath": requiredlevel.Deny,
		"NextEnd":    requiredlevel.Deny,
	},
	StateType.Parallel: {
		"Comment":    requiredlevel.Allowed,
		"InputPath":  requiredlevel.Allowed,
		"OutputPath": requiredlevel.Allowed,
		"Parameters": requiredlevel.Allowed,
		"ResultPath": requiredlevel.Allowed,
		"NextEnd":    requiredlevel.Required,
	},
	StateType.Map: {
		"Comment":    requiredlevel.Allowed,
		"InputPath":  requiredlevel.Allowed,
		"OutputPath": requiredlevel.Allowed,
		"Parameters": requiredlevel.Allowed,
		"ResultPath": requiredlevel.Allowed,
		"NextEnd":    requiredlevel.Required,
	},
	StateType.Suspend: {
		"Comment":    requiredlevel.Allowed,
		"InputPath":  requiredlevel.Allowed,
		"OutputPath": requiredlevel.Deny,
		"Parameters": requiredlevel.Deny,
		"ResultPath": requiredlevel.Deny,
		"NextEnd":    requiredlevel.Required,
	},
	StateType.StateGroup: {
		"Comment":    requiredlevel.Allowed,
		"InputPath":  requiredlevel.Allowed,
		"OutputPath": requiredlevel.Allowed,
		"Parameters": requiredlevel.Allowed,
		"ResultPath": requiredlevel.Allowed,
		"NextEnd":    requiredlevel.Allowed,
	},
}

// Fields Comment Field
var Fields = struct {
	Type             string
	Comment          string
	InputPath        string
	OutputPath       string
	Parameters       string
	ResultPath       string
	Next             string
	End              string
	NextEnd          string
	Resource         string
	Retry            string
	Catch            string
	TimeoutSeconds   string
	HeartbeatSeconds string
}{
	Type:             "Type",
	Comment:          "Comment",
	InputPath:        "InputPath",
	OutputPath:       "OutputPath",
	Parameters:       "Parameters",
	ResultPath:       "ResultPath",
	Next:             "Next",
	End:              "End",
	NextEnd:          "NextEnd",
	Resource:         "Resource",
	Retry:            "Retry",
	Catch:            "Catch",
	TimeoutSeconds:   "TimeoutSeconds",
	HeartbeatSeconds: "HeartbeatSeconds",
}

// CommonFields 常用字段
var CommonFields = []string{Fields.Comment, Fields.InputPath, Fields.OutputPath,
	Fields.Parameters, Fields.ResultPath, Fields.Next, Fields.End}

// 变量字段的前缀
var VariablePrefex = "$."

var StateMachineType = "statemachine"

// IsExecutableStateType 是否是可执行的步骤类型
func IsExecutableStateType(stype string) bool {
	return stype == StateType.Task ||
		stype == StateType.Wait ||
		stype == StateType.Succeed ||
		stype == StateType.Fail ||
		stype == StateType.Suspend ||
		stype == StateType.Choice
}
