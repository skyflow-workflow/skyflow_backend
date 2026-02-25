package states

// PassBody ...
type PassBody struct {
	Result map[string]any `mapstructure:"Result" validate:"required"`
}

// PassState ...
type PassState struct {
	*BaseState `json:",inline" mapstructure:"BaseState"`
	*PassBody  `json:",inline" mapstructure:"PassBody"`
}

func (p *PassBody) GetOutput(input any) (any, error) {
	return p.Result, nil
}

func (p *PassBody) GetDefinitionMap() (map[string]any, error) {
	data, err := DecodeStructToMap(p)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// NewPassStateFromString  Create New Pass State From String
func NewPassStateFromString(definition string) (state *PassState, err error) {

	data, err := StringToMap(definition)
	if err != nil {
		return
	}
	state, err = NewPassStateFromMap(data)
	return
}

// NewPassStateFromMap Create New Pass State From Map
func NewPassStateFromMap(data map[string]interface{}) (state *PassState, err error) {

	state = &PassState{}
	bs, err := NewBaseStateFromMap(data)
	if err != nil {
		return
	}

	body, err := NewPassBodyFromMap(data)
	if err != nil {
		return
	}
	state = &PassState{
		BaseState: bs,
		PassBody:  body,
	}

	return
}

func NewPassBodyFromMap(data map[string]interface{}) (*PassBody, error) {
	body := &PassBody{}
	err := DecodeMapToStruct(data, body)
	if err != nil {
		return nil, err
	}
	err = myValidate.Struct(body)
	if err != nil {
		return nil, err
	}
	err = InitPassBodyByMap(body, data)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// InitPassBodyByMap Initialize PassBody Content
func InitPassBodyByMap(body *PassBody, data map[string]interface{}) error {
	return nil
}

func (p *PassState) GetBaseState() *BaseState {
	return p.BaseState
}

// GetResult render result with input and parameters
func (p *PassState) GetResult(input any) (any, error) {
	var err error
	var result any
	if p.Result == nil {
		return nil, nil
	}
	stateinput, err := p.GetParametersInput(input)
	if err != nil {
		return nil, err
	}
	result, err = p.RenderParameters(stateinput, p.Result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetNextState Get Next State
func (p *PassState) GetNextState(input any) (NextState, error) {
	var err error
	var ns NextState
	result, err := p.GetResult(input)
	if err != nil {
		return ns, err
	}
	ns, err = p.BaseState.GetNextState(input, result)
	return ns, err
}

// GetDefinition 使用 DecodeStructToMap + FlattenMap + 按拒绝名单删 key，与 Fail 策略一致。
func (p *PassState) GetDefinition() (string, error) {
	baseData, err := p.BaseState.GetDefinitionMap()
	if err != nil {
		return "", err
	}
	bodyData, err := p.PassBody.GetDefinitionMap()
	if err != nil {
		return "", err
	}
	MapUpdate(baseData, bodyData)
	dataStr, err := ToString(baseData)
	return dataStr, err
}
