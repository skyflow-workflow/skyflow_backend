package states

/*
# @Date    : Thu Mar 14 2019
# @Author  : mmtbak(mumangtao@gmail.com)
# @Link    : link
# @Version : 1.0.0
# @Function: file for
*/

import (
	"github.com/go-playground/validator/v10"
)

// myValidate self define validator
var myValidate = validator.New()

// doc: https://states-language.net/spec.html#state-type-table-jsonpath
// 						Pass		Task		Choice		Wait		Succeed		Fail		Parallel    Map
// Type					Required	Required	Required	Required	Required	Required	Required	Required
// Comment				Allowed		Allowed		Allowed		Allowed		Allowed		Allowed		Allowed		Allowed
// InputPath/OutputPath	Allowed		Allowed		Allowed		Allowed		Allowed					Allowed		Allowed
// Parameters			Allowed		Allowed														Allowed		Allowed
// ResultPath			Allowed		Allowed														Allowed		Allowed
// One of: Next/End		Required	Required				Required							Required	Required
// Retry, Catch						Allowed														Allowed  	Allowed

// StateType state type
type StateType string

// StateTypes constant State type define here
var StateTypes = struct {
	Task     StateType
	Choice   StateType
	Fail     StateType
	Succeed  StateType
	Wait     StateType
	Pass     StateType
	Parallel StateType
	Map      StateType
	Suspend  StateType
	// 子流程类型, Parallel/Map 生成的中间步骤
	StateGroup StateType
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

// FiledRequiredLevel field required level
var FiledRequiredLevel = struct {
	Allowed  int
	Required int
	Deny     int
}{
	// default is allowed
	Allowed:  0,
	Required: 1,
	Deny:     2,
}

// CommonFields common fields
var CommonFields = []string{
	StateFieldNames.Comment,
	StateFieldNames.InputPath,
	StateFieldNames.OutputPath,
	StateFieldNames.Parameters,
	StateFieldNames.ResultPath,
}

// StateFieldRequired state field required
type StateFieldRequired struct {
	Comment    int
	Type       int
	InputPath  int
	OutputPath int
	Parameters int
	ResultPath int
	NextEnd    int
	Retry      int
	Catch      int
}

// StateFieldRequiredMap field required map
var StateFieldRequiredMap = map[StateType]StateFieldRequired{

	StateTypes.Task: {
		// default is allowed
		// Comment, InputPath, OutputPath, Parameters, ResultPath,
		Type:    FiledRequiredLevel.Required,
		NextEnd: FiledRequiredLevel.Required,
	},
	StateTypes.Parallel: {
		// default is allowed
		// Comment, InputPath, OutputPath, Parameters, ResultPath,
		Type:    FiledRequiredLevel.Required,
		NextEnd: FiledRequiredLevel.Required,
	},

	StateTypes.Map: {
		// default is allowed
		// Comment, InputPath, OutputPath, Parameters, ResultPath,
		Type:    FiledRequiredLevel.Required,
		NextEnd: FiledRequiredLevel.Required,
	},
	StateTypes.Pass: {
		// default is allowed
		// Comment, InputPath, OutputPath, Parameters, ResultPath
		Type:    FiledRequiredLevel.Required,
		NextEnd: FiledRequiredLevel.Required,
		// Pass 不支持 Retry, Catch
		Retry: FiledRequiredLevel.Deny,
		Catch: FiledRequiredLevel.Deny,
	},
	StateTypes.Wait: {
		// default  Comment, InputPath, OutputPath is allowed
		Type:       FiledRequiredLevel.Required,
		NextEnd:    FiledRequiredLevel.Required,
		ResultPath: FiledRequiredLevel.Deny,
		Parameters: FiledRequiredLevel.Deny,
		Retry:      FiledRequiredLevel.Deny,
		Catch:      FiledRequiredLevel.Deny,
	},
	StateTypes.Choice: {
		// default  Comment, InputPath, OutputPath is allowed
		Type:       FiledRequiredLevel.Required,
		NextEnd:    FiledRequiredLevel.Deny,
		ResultPath: FiledRequiredLevel.Deny,
		Parameters: FiledRequiredLevel.Deny,
		Retry:      FiledRequiredLevel.Deny,
		Catch:      FiledRequiredLevel.Deny,
	},

	StateTypes.Succeed: {
		// default  Comment, InputPath, OutputPath is allowed
		Type:       FiledRequiredLevel.Required,
		NextEnd:    FiledRequiredLevel.Deny,
		ResultPath: FiledRequiredLevel.Deny,
		Parameters: FiledRequiredLevel.Deny,
		Retry:      FiledRequiredLevel.Deny,
		Catch:      FiledRequiredLevel.Deny,
	},
	StateTypes.Fail: {
		// default  Comment is allowed
		Type:       FiledRequiredLevel.Required,
		InputPath:  FiledRequiredLevel.Deny,
		OutputPath: FiledRequiredLevel.Deny,
		NextEnd:    FiledRequiredLevel.Deny,
		ResultPath: FiledRequiredLevel.Deny,
		Parameters: FiledRequiredLevel.Deny,
		Retry:      FiledRequiredLevel.Deny,
		Catch:      FiledRequiredLevel.Deny,
	},

	StateTypes.Suspend: {
		// default  Comment is allowed
		Type:       FiledRequiredLevel.Required,
		InputPath:  FiledRequiredLevel.Deny,
		OutputPath: FiledRequiredLevel.Deny,
		NextEnd:    FiledRequiredLevel.Deny,
		ResultPath: FiledRequiredLevel.Deny,
		Parameters: FiledRequiredLevel.Deny,
		Retry:      FiledRequiredLevel.Deny,
		Catch:      FiledRequiredLevel.Deny,
	},
	StateTypes.StateGroup: {
		// default  Comment is allowed
	},
}

// StateFieldNames common field names
var StateFieldNames = struct {
	Type             string
	Comment          string
	InputPath        string
	OutputPath       string
	Parameters       string
	ResultPath       string
	Next             string
	End              string
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
	Resource:         "Resource",
	Retry:            "Retry",
	Catch:            "Catch",
	TimeoutSeconds:   "TimeoutSeconds",
	HeartbeatSeconds: "HeartbeatSeconds",
}

var StateMachineFieldNames = struct {
	StartAt string
	States  string
}{
	StartAt: "StartAt",
	States:  "States",
}

// VariablePrefix 变量字段的前缀
var VariablePrefix = "$."

// StateMachineType ...
var StateMachineType = "statemachine"

// IsExecutableStateType 是否是可执行的步骤类型
func IsExecutableStateType(stype string) bool {
	sttype := StateType(stype)
	isExecuteable := sttype == StateTypes.Task ||
		sttype == StateTypes.Wait ||
		sttype == StateTypes.Succeed ||
		sttype == StateTypes.Fail ||
		sttype == StateTypes.Suspend ||
		sttype == StateTypes.Choice
	return isExecuteable
}

// QueryLanguageType 查询语言类型
type QueryLanguageType string

// QueryLanguages ...
var QueryLanguages = struct {
	JSONPath QueryLanguageType
	JSONata  QueryLanguageType
}{
	JSONPath: "JSONPath",
	JSONata:  "JSONata",
}
