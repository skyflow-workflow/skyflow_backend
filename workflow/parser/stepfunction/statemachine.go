package stepfunction

import (
	"context"

	"github.com/skyflow-workflow/skyflow_backbend/workflow/parser/states"
)

// StateMachineBody ...
type StateMachineBody struct {
	StartAt string
	States  map[string]map[string]any
}

// DecodeStateMachine ...
func (decoder *StepfuncionDecoder) DecodeStateMachine(ctx context.Context, data map[string]any) (*states.StateMachine, error) {
	// Parse the state machine
	var err error
	var header *states.StateMachineHeader
	var body *states.StateMachineBody
	header, err = decoder.DecodeStateMachineHeader(ctx, data)
	if err != nil {
		return nil, err
	}

	body, err = decoder.DecodeStateMachineBody(ctx, data)
	if err != nil {
		return nil, err
	}
	sm := &states.StateMachine{
		StateMachineHeader: header,
		StateMachineBody:   body,
	}
	return sm, nil
}

// DecodeStateMachineBody ...
func (decoder *StepfuncionDecoder) DecodeStateMachineBody(ctx context.Context, data map[string]any) (
	*states.StateMachineBody, error) {

	var err error
	bodydata := StateMachineBody{}
	err = decoder.MapDecode(data, &bodydata)
	if err != nil {
		return nil, err
	}
	body := &states.StateMachineBody{
		States: make(map[string]states.State, len(bodydata.States)),
	}
	body.SetStartAt(bodydata.StartAt)

	ctx = decoder.AddCtxDecodePath(ctx, StateMachineFields.States)
	for name, stateData := range bodydata.States {
		if name == "" {
			err = decoder.NewFieldPathError(ctx, ErrorStateNameEmpty)
			return nil, err
		}
		stateCtx := decoder.AddCtxDecodePath(ctx, name)
		newstate, err := decoder.DecodeState(stateCtx, stateData)
		if err != nil {
			return nil, err
		}
		newstate.SetName(name)
		body.AddState(newstate)
	}
	err = body.Validate()
	if err != nil {
		return nil, err
	}

	return body, nil
}
