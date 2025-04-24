package stepfunction

import (
	"context"

	"github.com/skyflow-workflow/skyflow_backbend/workflow/parser/states"
)

// TaskBody task body for stepfunction
type TaskBody struct {
	Resource         string           `mapstructure:"Resource"`
	TimeoutSeconds   uint             `mapstructure:"TimeoutSeconds"`
	HeartbeatSeconds uint             `mapstructure:"HeartbeatSeconds"`
	Retry            []map[string]any `mapstructure:"Retry"`
	Catch            []map[string]any `mapstructure:"Catch"`
}

// DefaultTaskBody ...
var DefaultTaskBody = TaskBody{
	Resource:         "",
	TimeoutSeconds:   0,
	HeartbeatSeconds: 0,
}

// DecodeTaskState ...
func (decoder *StepfuncionDecoder) DecodeTaskState(ctx context.Context, basestate *states.BaseState, data map[string]any) (
	states.State, error) {
	taskbody, err := decoder.DecodeTaskBody(ctx, data)
	if err != nil {
		return nil, err
	}
	task := &states.Task{
		BaseState: basestate,
		TaskBody:  taskbody,
	}
	return task, nil
}

// DecodeTaskBody ...
func (decoder *StepfuncionDecoder) DecodeTaskBody(ctx context.Context, data map[string]any) (
	*states.TaskBody, error) {
	taskbody := states.DefaultTaskBody
	err := decoder.MapDecode(data, &taskbody)
	if err != nil {
		return nil, err
	}
	err = taskbody.Validate()
	if err != nil {
		return nil, err
	}
	return &taskbody, nil
}
