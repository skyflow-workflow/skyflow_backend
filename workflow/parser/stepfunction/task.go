package stepfunction

import "github.com/skyflow-workflow/skyflow_backbend/workflow/parser/states"

// TaskBody task body umashanall
type TaskBody struct {
	Resource         string                   `mapstructure:"Resource"`
	TimeoutSeconds   uint                     `mapstructure:"TimeoutSeconds"`
	HeartbeatSeconds uint                     `mapstructure:"HeartbeatSeconds"`
	Retry            []map[string]interface{} `mapstructure:"Retry"`
	Catch            []map[string]interface{} `mapstructure:"Catch"`
}

// DefaultTaskBody ...
var DefaultTaskBody = TaskBody{
	Resource:         "",
	TimeoutSeconds:   0,
	HeartbeatSeconds: 0,
}

// DecodeTaskState ...
func (sfdecoder *StepfuncionDecoder) DecodeTaskState(basestate *states.BaseState, data map[string]interface{}) (
	states.State, error) {
	var err error

	// states.taskbody
	taskbody := states.TaskBody{}

	// stepfunction taskbody
	sftaskbody := DefaultTaskBody
	err = sfdecoder.MapDecode(data, &sftaskbody)
	if err != nil {
		return nil, err
	}

	var retrynodes = []states.RetryNode{}
	if len(sftaskbody.Retry) > 0 {
		// retry 不为空
		for _, retrynodedata := range sftaskbody.Retry {
			node := states.DefaultRetryNode
			err = sfdecoder.MapDecode(retrynodedata, &node)
			if err != nil {
				return nil, err
			}
			retrynodes = append(retrynodes, node)
		}
	}

	// Catch
	var catchnodes = []states.CatchNode{}
	if len(sftaskbody.Catch) > 0 {
		for _, catchnodedata := range sftaskbody.Catch {
			node := states.DefaultCatchNode
			err = sfdecoder.MapDecode(catchnodedata, &node)
			if err != nil {
				return nil, err
			}
			catchnodes = append(catchnodes, node)
		}
	}
	taskbody = states.TaskBody{
		Resource:         sftaskbody.Resource,
		HeartbeatSeconds: sftaskbody.HeartbeatSeconds,
		TimeoutSeconds:   sftaskbody.TimeoutSeconds,
		Retry:            retrynodes,
		Catch:            catchnodes,
	}

	taskstate := &states.Task{
		BaseState: basestate,
		TaskBody:  &taskbody,
	}
	err = taskstate.Init()
	if err != nil {
		return nil, err
	}
	return taskstate, nil
}
