package states

import "fmt"

// BaseStateFieldRequired is a struct that validates the required fields of a state
type BaseStateFieldRequired struct {
	Type       *string
	Comment    *string
	Parameters interface{}
	InputPath  interface{}
	OutputPath interface{}
	ResultPath interface{}
	Next       interface{}
	End        interface{}
	Retry      interface{}
	Catch      interface{}
}

// Validate ...
func (s *BaseStateFieldRequired) Validate() error {

	stype := s.Type
	if stype == nil {
		return NewFieldPathError(fmt.Errorf("%w: Type", ErrorLackOfRequiredField))
	}
	sttype := StateType(*stype)
	staterequire, ok := StateFeildRequiredMap[sttype]
	if !ok {
		return NewFieldPathError(fmt.Errorf("%w: %s", ErrorInvalidStateType, *stype), StateFields.Type)
	}
	// check required fields
	// check comment
	if staterequire.Comment == FiledRequiredLevel.Required && s.Comment == nil {
		return NewFieldPathError(fmt.Errorf("%w: %s", ErrorLackOfRequiredField, StateFields.Comment),
			StateFields.Comment)
	}
	// check inputpath
	if staterequire.InputPath == FiledRequiredLevel.Required && s.InputPath == nil {
		return NewFieldPathError(fmt.Errorf("%w: InputPath", ErrorLackOfRequiredField))
	} else if staterequire.InputPath == FiledRequiredLevel.Deny && s.InputPath != nil {
		return NewFieldPathError(fmt.Errorf("%w: InputPath", ErrorFiledDenied), StateFields.InputPath)
	}
	// check outputpath
	if staterequire.OutputPath == FiledRequiredLevel.Required && s.OutputPath == nil {
		return NewFieldPathError(fmt.Errorf("%w: OutputPath", ErrorLackOfRequiredField))
	} else if staterequire.OutputPath == FiledRequiredLevel.Deny && s.OutputPath != nil {
		return NewFieldPathError(fmt.Errorf("%w: OutputPath", ErrorFiledDenied), StateFields.OutputPath)
	}
	// check parameters
	if staterequire.Parameters == FiledRequiredLevel.Required && s.Parameters == nil {
		return NewFieldPathError(fmt.Errorf("%w: Parameters", ErrorLackOfRequiredField))
	} else if staterequire.Parameters == FiledRequiredLevel.Deny && s.Parameters != nil {
		return NewFieldPathError(fmt.Errorf("%w: Parameters", ErrorFiledDenied), StateFields.Parameters)
	}
	// check resultpath
	if staterequire.ResultPath == FiledRequiredLevel.Required && s.ResultPath == nil {
		return NewFieldPathError(fmt.Errorf("%w: ResultPath", ErrorLackOfRequiredField))
	} else if staterequire.ResultPath == FiledRequiredLevel.Deny && s.ResultPath != nil {
		return NewFieldPathError(fmt.Errorf("%w: ResultPath", ErrorFiledDenied), StateFields.ResultPath)
	}
	// check nextend
	// if nextend is deny , next and end should be nil
	if staterequire.NextEnd == FiledRequiredLevel.Deny {
		if s.Next != nil {
			return NewFieldPathError(fmt.Errorf("%w: Next", ErrorFiledDenied), StateFields.Next)
		}
		if s.End != nil {
			return NewFieldPathError(fmt.Errorf("%w: End", ErrorFiledDenied), StateFields.End)
		}
	}
	if staterequire.NextEnd == FiledRequiredLevel.Required {
		if s.Next == nil && s.End == nil {
			return NewFieldPathError(fmt.Errorf("%w: Next or End", ErrorLackOfRequiredField))
		}
	}
	return nil
}

// BaseState is a struct that defines the base state of a state machine, with default values
type BaseState struct {
	Name               string      `json:"Name,omitempty"`
	Type               string      `json:"Type,omitempty"`
	Comment            string      `json:"Comment,omitempty"`
	InputPath          string      `json:"InputPath,omitempty"`
	OutputPath         string      `json:"OutputPath,omitempty"`
	ResultPath         string      `json:"ResultPath,omitempty"`
	Parameters         interface{} `json:"Parameters,omitempty"`
	MaxExecuteTimes    int         `json:"MaxExecuteTimes,omitempty"`
	End                bool        `json:"End,omitempty"`
	Next               string      `json:"Next,omitempty"`
	Retry              interface{} `json:"Retry,omitempty"`
	Catch              interface{} `json:"Catch,omitempty"`
	_Bone              *StateBone
	_InputPathPattern  JSONPathCompiled
	_OutputPathPattern JSONPathCompiled
	_ResultPathPattern JSONPathCompiled
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

// Init ...
func (s *BaseState) Init() error {

	var err error
	// 校验input结构
	if s.InputPath != "" {
		s._InputPathPattern, err = JSONPathParse(s.InputPath)
		if err != nil {
			return NewFieldPathError(fmt.Errorf("%w:%w", ErrorInvalidFiledContent, err), StateFields.InputPath)
		}
	}
	if s.OutputPath != "" {
		s._OutputPathPattern, err = JSONPathParse(s.OutputPath)
		if err != nil {
			return NewFieldPathError(fmt.Errorf("%w:%w", ErrorInvalidFiledContent, err), StateFields.OutputPath)
		}
	}
	if s.ResultPath != "" {
		s._ResultPathPattern, err = JSONPathParse(s.ResultPath)
		if err != nil {
			return NewFieldPathError(fmt.Errorf("%w:%w", ErrorInvalidFiledContent, err), StateFields.ResultPath)
		}
	}
	if s.MaxExecuteTimes <= 0 {
		return fmt.Errorf("field 'MaxExecuteTimes' must gt 0 ")
	}
	// Bone
	s._Bone = &StateBone{
		BaseBone: BaseBone{
			Name:    s.Name,
			Comment: s.Comment,
			Type:    s.Type,
			Next:    []string{},
			End:     s.End,
		},
	}
	if s.Next != "" {
		s._Bone.Next = []string{s.Next}
	}

	return nil
}

// GetBone ...
func (s *BaseState) GetBone() StateBone {
	return *s._Bone
}
