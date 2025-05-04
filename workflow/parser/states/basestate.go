package states

// basestate is a struct that defines the base state of a state machine, with default values
// processing input/output flow defined in https://docs.amazonaws.cn/en_us/step-functions/latest/dg/workflow-variables.html

import (
	"fmt"
	"strings"

	"github.com/skyflow-workflow/skyflow_backbend/pkg/jsonpath"
)

// BaseState is a struct that defines the base state of a state machine, with default values
type BaseState struct {
	Name            string `json:"Name,omitempty"`
	Type            string `json:"Type,omitempty"`
	Comment         string `json:"Comment,omitempty"`
	InputPath       string `json:"InputPath,omitempty"`
	OutputPath      string `json:"OutputPath,omitempty"`
	ResultPath      string `json:"ResultPath,omitempty"`
	Parameters      any    `json:"Parameters,omitempty"`
	MaxExecuteTimes int    `json:"MaxExecuteTimes,omitempty"`
	End             bool   `json:"End,omitempty"`
	Next            string `json:"Next,omitempty"`
	Retry           any    `json:"Retry,omitempty"`
	Catch           any    `json:"Catch,omitempty"`
}

// GetName ...
func (s *BaseState) GetName() string {
	return s.Name
}

// SetName ...
func (s *BaseState) SetName(name string) {
	s.Name = name
}

// GetType ...
func (s *BaseState) GetType() string {
	return s.Type
}

// Validate ...
func (s *BaseState) Validate() error {

	var err error
	// 校验input结构
	if s.InputPath != "" {
		_, err = jsonpath.JsonPathCompile(s.InputPath)
		if err != nil {
			return NewFieldPathError(fmt.Errorf("%w:%w", ErrorInvalidFiledContent, err), StateFieldNames.InputPath)
		}
	}
	if s.OutputPath != "" {
		_, err = jsonpath.JsonPathCompile(s.OutputPath)
		if err != nil {
			return NewFieldPathError(fmt.Errorf("%w:%w", ErrorInvalidFiledContent, err), StateFieldNames.OutputPath)
		}
	}
	if s.ResultPath != "" {
		_, err = jsonpath.JsonPathCompile(s.ResultPath)
		if err != nil {
			return NewFieldPathError(fmt.Errorf("%w:%w", ErrorInvalidFiledContent, err), StateFieldNames.ResultPath)
		}
	}

	return nil
}

// Init ...
func (s *BaseState) Init() error {
	return s.Validate()
}

// GetBone ...
func (s *BaseState) GetBone() StateBone {
	// Bone
	bone := StateBone{
		BaseBone: BaseBone{
			Name:    s.Name,
			Comment: s.Comment,
			Type:    s.Type,
			Next:    []string{},
			End:     s.End,
		},
	}
	if s.Next != "" {
		bone.Next = []string{s.Next}
	}
	return bone
}

// GetOutput calculate task output
// cautions:
// 1. return nil if outputpath is empty
// 2. return output is not a copy of input, it will reference data from input and output, which is defined in resultpath and outputpath
// @input input state origin input
// @output output state produce data
// return output state produce data
func (s *BaseState) GetOutput(input any, output any) (any, error) {
	newoutput, err := s.GetOutputWithPath(input, output, s.ResultPath, s.OutputPath)
	if err != nil {
		return nil, err
	}
	return newoutput, nil
}

// GetOutputWithPath calculate output with resultpath and outputpath
// cautions:
// 1. return nil if outputpath is empty
// 2. return output is not a copy of input, it will reference data from input and output, which is defined in resultpath and outputpath
// @input input state origin input
// @output output state produce data
// @resultpath resultpath
// @outputpath outputpath
func (s *BaseState) GetOutputWithPath(input any, output any, resultpath string, outputpath string) (any, error) {

	var err error

	// if outputpath is empty, return nil
	if outputpath == "" {
		return nil, nil
	}
	// if output is not empty, calculate new output
	var result = input

	if resultpath != "" {
		err = jsonpath.JsonPathSetValue(resultpath, input, output)
		if err != nil {
			return nil, err
		}
	}
	if outputpath != "" {
		result, err = jsonpath.JsonPathGetValue(outputpath, result)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

// GetParametersInput get state parameters input
// input , origin input,
// output state input using inputpath , parameters
func (s *BaseState) GetParametersInput(input any) (any, error) {
	var result = input
	var err error
	if s.InputPath != "" {
		result, err = jsonpath.JsonPathGetValue(s.InputPath, input)
		if err != nil {
			return nil, err
		}
	}

	if s.Parameters != nil {
		result, err = s.GenParameters(result)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

// GenParameters calculate parameters by input data  and parameters
// @input input state origin input
// return parameters
func (s *BaseState) GenParameters(input any) (any, error) {

	parameters := s.Parameters
	newparameters, err := s.RenderParameters(input, parameters)
	return newparameters, err
}

// RenderParameters render parameters
// @input input state origin input
// @parameters parameters field parameters
// return parameters
func (s *BaseState) RenderParameters(input any, parameters any) (any, error) {

	// if parameters is nil, return input
	if parameters == nil {
		return input, nil
	}
	parametersMap, ok := parameters.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("parameters should be a map")
	}

	specialchar := ".$"

	// GetResult get execution result, recursive call
	type GetResult func(any) (any, error)
	var getResult GetResult

	getResult = func(param any) (any, error) {
		var result any

		// try map
		if paramap, ok := param.(map[string]any); ok {
			var mapresult = map[string]any{}
			for key, value := range paramap {
				key = strings.TrimSpace(key)
				if strings.HasSuffix(key, specialchar) {
					valuestr, valueStringOk := value.(string)
					if !valueStringOk {
						return nil, fmt.Errorf("key [ %s ] 's value should be string ", key)
					}
					newvalue, err := jsonpath.JsonPathGetValue(valuestr, input)
					if err != nil {
						return nil, err
					}
					newkey := strings.TrimSuffix(key, specialchar)
					mapresult[newkey] = newvalue

				} else {
					newvalue, err := getResult(value)
					if err != nil {
						return nil, err
					}
					mapresult[key] = newvalue
				}
			}
			result = mapresult
			return result, nil
		}
		// try array
		if paramArray, ok := param.([]any); ok {
			var arrayresult []any
			for _, item := range paramArray {
				newvalue, err := getResult(item)
				if err != nil {
					return nil, err
				}
				arrayresult = append(arrayresult, newvalue)

			}
			result = arrayresult
			return result, nil
		}
		return param, nil
	}
	genparameters, err := getResult(parametersMap)
	return genparameters, err

}

// GetNextState get next state info
// generate final output using inputpath parameters resultpath outputpath
// @input state origin input
// @output state produce data
// return next state info
func (s *BaseState) GetNextState(input any, output any) (NextState, error) {
	var err error
	var ns = NextState{}
	finaloutput, err := s.GetOutput(input, output)
	if err != nil {
		return ns, err
	}

	ns = NextState{
		Name:   s.Next,
		Output: finaloutput,
	}

	return ns, nil
}

// ValidateStateFieldOptional validate state field optional
func ValidateStateFieldOptional(data map[string]any) error {

	if data == nil {
		return NewFieldPathError(ErrorInvalidData)
	}
	// validate field "Type"
	stype, ok := data[StateFieldNames.Type]
	if !ok {
		return NewFieldPathError(fmt.Errorf("%w: %s", ErrorLackOfRequiredField, StateFieldNames.Type))
	}
	statetype, ok := stype.(string)
	if !ok {
		return NewFieldPathError(fmt.Errorf("%w: %s", ErrorInvalidData, StateFieldNames.Type))
	}
	stateRequired, ok := StateFieldRequiredMap[StateType(statetype)]
	if !ok {
		return NewFieldPathError(fmt.Errorf("%w: %s", ErrorInvalidStateType, statetype), StateFieldNames.Type)
	}

	// validate required fields

	// check required fields
	// check comment
	if stateRequired.Comment == FiledRequiredLevel.Required && data[StateFieldNames.Comment] == nil {
		return NewFieldPathError(fmt.Errorf("%w: %s", ErrorLackOfRequiredField, StateFieldNames.Comment),
			StateFieldNames.Comment)
	}
	// check inputpath
	if stateRequired.InputPath == FiledRequiredLevel.Required && data[StateFieldNames.InputPath] == nil {
		return NewFieldPathError(
			fmt.Errorf("%w: %s", ErrorLackOfRequiredField, StateFieldNames.InputPath),
			StateFieldNames.InputPath)
	} else if stateRequired.InputPath == FiledRequiredLevel.Deny && data[StateFieldNames.InputPath] != nil {
		return NewFieldPathError(
			fmt.Errorf("%w: %s", ErrorFiledDenied, StateFieldNames.InputPath),
			StateFieldNames.InputPath)
	}
	// check outputpath
	if stateRequired.OutputPath == FiledRequiredLevel.Required && data[StateFieldNames.OutputPath] == nil {
		return NewFieldPathError(
			fmt.Errorf("%w: %s", ErrorLackOfRequiredField, StateFieldNames.OutputPath),
			StateFieldNames.OutputPath)
	} else if stateRequired.OutputPath == FiledRequiredLevel.Deny && data[StateFieldNames.OutputPath] != nil {
		return NewFieldPathError(
			fmt.Errorf("%w: %s", ErrorFiledDenied, StateFieldNames.OutputPath),
			StateFieldNames.OutputPath)
	}
	// check parameters
	if stateRequired.Parameters == FiledRequiredLevel.Required && data[StateFieldNames.Parameters] == nil {
		return NewFieldPathError(
			fmt.Errorf("%w: %s", ErrorLackOfRequiredField, StateFieldNames.Parameters),
			StateFieldNames.Parameters)
	} else if stateRequired.Parameters == FiledRequiredLevel.Deny && data[StateFieldNames.Parameters] != nil {
		return NewFieldPathError(
			fmt.Errorf("%w: %s", ErrorFiledDenied, StateFieldNames.Parameters),
			StateFieldNames.Parameters)
	}
	// check resultpath
	if stateRequired.ResultPath == FiledRequiredLevel.Required && data[StateFieldNames.ResultPath] == nil {
		return NewFieldPathError(
			fmt.Errorf("%w: %s", ErrorLackOfRequiredField, StateFieldNames.ResultPath),
			StateFieldNames.ResultPath)
	} else if stateRequired.ResultPath == FiledRequiredLevel.Deny && data[StateFieldNames.ResultPath] != nil {
		return NewFieldPathError(
			fmt.Errorf("%w: %s", ErrorFiledDenied, StateFieldNames.ResultPath),
			StateFieldNames.ResultPath)
	}

	// check nextend
	// if nextend is deny , next and end should be nil
	// if nextend is required , one of next or end is not nil and other is nil
	if stateRequired.NextEnd == FiledRequiredLevel.Deny {
		if data[StateFieldNames.Next] != nil {
			return NewFieldPathError(fmt.Errorf("%w: %s", ErrorFiledDenied, StateFieldNames.Next),
				StateFieldNames.Next)
		}
		if data[StateFieldNames.End] != nil {
			return NewFieldPathError(fmt.Errorf("%w: %s", ErrorFiledDenied, StateFieldNames.End),
				StateFieldNames.End)
		}
	} else if stateRequired.NextEnd == FiledRequiredLevel.Required {
		// one of next or end is not nil and other is nil
		if data[StateFieldNames.Next] == nil && data[StateFieldNames.End] == nil {
			return NewFieldPathError(fmt.Errorf("%w: %s or %s", ErrorLackOfRequiredField,
				StateFieldNames.Next, StateFieldNames.End),
				StateFieldNames.Next)
		} else {
			// both next and end are not nil
			if data[StateFieldNames.Next] != nil && data[StateFieldNames.End] != nil {
				return NewFieldPathError(fmt.Errorf("%w: %s and %s should not be both defined", ErrorInvalidField,
					StateFieldNames.Next, StateFieldNames.End),
					StateFieldNames.Next)
			}
		}
	}

	return nil
}
