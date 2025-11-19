package states

import (
	"fmt"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/skyflow-workflow/skyflow_backend/pkg/jsonpath"
)

// HeaderFieldNames ...
var HeaderFieldNames = struct {
	Version             string
	Type                string
	Comment             string
	QueryLanguage       string
	TimeoutSeconds      string
	AbortTimeoutSeconds string
	ExecutionInfoPath   string
	DefaultInput        string
}{
	Version:             "Version",
	Type:                "Type",
	Comment:             "Comment",
	QueryLanguage:       "QueryLanguage",
	TimeoutSeconds:      "TimeoutSeconds",
	AbortTimeoutSeconds: "AbortTimeoutSeconds",
	ExecutionInfoPath:   "ExecutionInfoPath",
	DefaultInput:        "DefaultInput",
}

// StateMachineHeader statemachine的header 定义
type StateMachineHeader struct {
	Version       string `json:"Version" mapstructure:"Version" ` // 语法版本， 默认1.0
	Type          string `json:"Type" mapstructure:"Type"`
	Comment       string `json:"Comment" mapstructure:"Comment"`
	QueryLanguage string `json:"QueryLanguage" mapstructure:"QueryLanguage"`
	// 任务执行超时时间， 过期自动设置任务状态为失败状态
	TimeoutSeconds int `json:"TimeoutSeconds" mapstructure:"TimeoutSeconds" validate:"gte=0"`
	// 任务执行超时终止时间，找过指定时间后，如果任务状态为非Success，任务自动终止并且设置为 Abort状态
	AbortTimeoutSeconds int `json:"AbortTimeoutSeconds" mapstructure:"AbortTimeoutSeconds" validate:"gte=0"`
	// Execution 的信息路径
	ExecutionInfoPath string `json:"ExecutionInfoPath" mapstructure:"ExecutionInfoPath"`
	// 默认input参数
	DefaultInput map[string]interface{} `json:"DefaultInput" mapstructure:"DefaultInput"`
	// defintion
	_Definition string
}

// StateMachineTimeout StateMachine 超时的结构
type StateMachineTimeout struct {
	Timeout      time.Duration
	TaskTimeout  time.Duration
	AbortTimeout time.Duration
}

func NewDefautStateMachineHeader() StateMachineHeader {
	return StateMachineHeader{
		Version:             "1.0",
		Type:                StateMachineType,
		TimeoutSeconds:      0,
		AbortTimeoutSeconds: 0,
		QueryLanguage:       string(QueryLanguages.JSONPath),
		ExecutionInfoPath:   "",
	}
}

// NewStateMachineHeaderFromMap NewStateMachineHeaderFromMap
func NewStateMachineHeaderFromMap(data map[string]interface{}) (*StateMachineHeader, error) {
	var err error
	header := NewDefautStateMachineHeader()
	headerptr := &header
	err = InitStateMachineHeaderByMap(headerptr, data)
	if err != nil {
		return nil, err
	}
	return headerptr, err
}

// InitByMap 通过Map类型初始化
func InitStateMachineHeaderByMap(header *StateMachineHeader, data map[string]interface{}) error {

	var err error

	// 解析header
	err = mapstructure.Decode(data, header)
	if err != nil {
		return err
	}
	err = header.Init()
	if err != nil {
		return err
	}
	return nil
}

// Init ...
func (header *StateMachineHeader) Init() error {
	var err error
	err = myvalidate.Struct(header)
	if err != nil {
		return err
	}
	if header.Version == "" {
		err = NewFieldPathError(ErrorInvalidFiledContent, HeaderFieldNames.Version)
		return err
	}
	if header.Type == "" {
		err = NewFieldPathError(ErrorInvalidFiledContent, HeaderFieldNames.Type)
		return err
	}
	// TODO  validate QueryLanguage later
	// if header.QueryLanguage == "" {
	// 	err = NewFieldError(ErrorInvalidFiledContent, HeaderFieldNames.QueryLanguage)
	// 	return err
	// }

	if header.ExecutionInfoPath != "" {
		_, err = jsonpath.JsonPathCompile(header.ExecutionInfoPath)
		if err != nil {
			return fmt.Errorf("field 'ExecutionInfoPath' parse error : %w", err)
		}
	}

	// definition
	def, err := ToString(header)
	if err != nil {
		return err
	}
	header._Definition = def

	return nil
}

// NewStateMachineHeaderFromString NewStateMachineHeaderFromString
func NewStateMachineHeaderFromString(definition string) (*StateMachineHeader, error) {

	data, err := StringToMap(definition)
	if err != nil {
		return nil, err
	}
	state, err := NewStateMachineHeaderFromMap(data)
	return state, err
}

// GetInput 计算新的输入
func (h StateMachineHeader) GetInput(input interface{}, executioninfo interface{}) (interface{}, error) {

	var outputdata = input
	var err error
	// 处理 ExecutionInfoPath
	if h.ExecutionInfoPath != "" {
		err = jsonpath.JsonPathSetValue(h.ExecutionInfoPath, input, executioninfo)
		if err != nil {
			return nil, err
		}
	}
	if h.DefaultInput != nil {
		outputdatamap, ok := outputdata.(map[string]interface{})
		if !ok {
			return outputdata, fmt.Errorf("input is not map")
		}
		MapUpdate(h.DefaultInput, outputdatamap)
		outputdata = h.DefaultInput

	}
	return outputdata, nil
}

// GetDefinition returan header definition
func (h StateMachineHeader) GetDefinition() string {
	return h._Definition
}

// GetTimeout return statemachine timeout
func (h StateMachineHeader) GetTimeout() StateMachineTimeout {
	timeout := StateMachineTimeout{
		Timeout:      time.Second * time.Duration(h.TimeoutSeconds),
		AbortTimeout: time.Second * time.Duration(h.AbortTimeoutSeconds),
	}
	return timeout
}
