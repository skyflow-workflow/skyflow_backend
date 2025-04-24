package states

// PassBody ...
type PassBody struct {
	Result map[string]any `mapstructure:"Result" validate:"required"`
}

// Pass ...
type Pass struct {
	*BaseState
	*PassBody
}

func (p *PassBody) GetOutput(input any) (any, error) {
	return p.Result, nil
}

// GetResult render result with input and parameters
func (p *Pass) GetResult(input any) (any, error) {
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
func (p *Pass) GetNextState(input any) (NextState, error) {
	var err error
	var ns NextState
	result, err := p.GetResult(input)
	if err != nil {
		return ns, err
	}
	ns, err = p.BaseState.GetNextState(input, result)
	return ns, err
}
