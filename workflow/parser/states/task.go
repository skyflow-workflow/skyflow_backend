package states

import (
	"fmt"
	"math"
	"slices"
	"time"
)

// TaskBody ...
type TaskBody struct {
	// Block : if block execution for task, default false
	// if true, the task will block the workflow execution until the manual call to resume the task
	// if false, the task will be executed automatically
	Block            bool   `mapstructure:"Block"`
	Resource         string `mapstructure:"Resource" validate:"required,gt=0"`
	TimeoutSeconds   uint   `mapstructure:"TimeoutSeconds" validate:"gte=0"`
	HeartbeatSeconds uint   `mapstructure:"HeartbeatSeconds" validate:"gte=0"`
	// Retry for decode
	Retry []TaskRetryNode `mapstructure:"Retry"`
	// Catch for decode
	Catch []TaskCatchNode `mapstructure:"Catch"`
}

// TaskRetryNode struct for retry
type TaskRetryNode struct {
	ErrorEquals     []string `mapstructure:"ErrorEquals"`
	IntervalSeconds uint     `mapstructure:"IntervalSeconds"`
	MaxAttempts     uint     `mapstructure:"MaxAttempts"`
	BackoffRate     float64  `mapstructure:"BackoffRate"`
}

// TaskCatchNode struct for catch
type TaskCatchNode struct {
	ErrorEquals []string `mapstructure:"ErrorEquals"`
	Next        string   `mapstructure:"Next"`
	ResultPath  string   `mapstructure:"ResultPath"`
}

// TaskSendData  task state execute and send data
type TaskSendData struct {
	Success bool     // is task submit success
	Retry   []int    // current retry times, index is retry times for each Retry branch
	Errors  []string // send error list
	Output  any      // task state execute output
}

// DefaultRetryNode default retry config for task
var DefaultRetryNode = TaskRetryNode{
	ErrorEquals:     []string{},
	IntervalSeconds: 1,
	MaxAttempts:     3,
	BackoffRate:     1.5,
}

// DefaultCatchNode default catch config for task
var DefaultCatchNode = TaskCatchNode{
	ErrorEquals: []string{},
	Next:        "",
	ResultPath:  "$",
}

// DefaultTaskBody ...
var DefaultTaskBody = TaskBody{
	// default false, means task will not block the workflow execution
	Block:    false,
	Resource: "",
	// default 0, 0 means no timeout limit
	TimeoutSeconds: 0,
	// default 0, 0 means no heartbeat timeout limit
	HeartbeatSeconds: 0,
	// default empty, no retry
	Retry: []TaskRetryNode{},
	// default empty, no catch
	Catch: []TaskCatchNode{},
}

// NewTaskStateFromString  Create New Wait State
func NewTaskStateFromString(definition string) (state *TaskState, err error) {
	// state, err = NewGenericStateFromString(definition, NewTaskStateFromMap)
	mapdata, err := StringToMap(definition)
	if err != nil {
		return
	}

	state, err = NewTaskStateFromMap(mapdata)
	if err != nil {
		return
	}
	return
}

// NewTaskStateFromMap NewTaskStateFromMap
func NewTaskStateFromMap(data map[string]interface{}) (state *TaskState, err error) {

	// basestate
	bs, err := NewBaseStateFromMap(data)
	if err != nil {
		return
	}
	taskbody, err := NewTaskBodyFromMap(data)
	if err != nil {
		return
	}
	task := &TaskState{
		BaseState: bs,
		TaskBody:  taskbody,
	}
	err = task.Init()
	if err != nil {
		return
	}
	return task, err
}

// NewTaskBodyFromMap  NewTaskBodyFromMap
func NewTaskBodyFromMap(data map[string]interface{}) (*TaskBody, error) {

	taskbody := DefaultTaskBody
	err := InitTaskBodyByMap(&taskbody, data)
	if err != nil {
		return nil, err
	}
	return &taskbody, nil
}

// InitByMap Inititalize TaskState Content
func InitTaskBodyByMap(body *TaskBody, data map[string]interface{}) error {

	// 数据初始化
	var err error
	// 初始化自身
	err = DecodeMapToStruct(data, body)
	if err != nil {
		return err
	}
	err = myvalidate.Struct(body)
	if err != nil {
		return err
	}
	// Retry

	if len(body.Retry) > 0 {
		// retry 不为空
		var retrynodes []TaskRetryNode
		inputretrynodes := data[StateFieldNames.Retry].([]interface{})
		for idx, retrynodedata := range inputretrynodes {
			retrynodeMapData, ok := retrynodedata.(map[string]interface{})
			if !ok {
				return fmt.Errorf("'Retry' Node [ %d ] data should be map", idx+1)
			}
			node := DefaultRetryNode
			err = DecodeMapToStruct(retrynodeMapData, &node)
			if err != nil {
				return err
			}
			retrynodes = append(retrynodes, node)
		}
		body.Retry = retrynodes
	}

	// Catch
	if len(body.Catch) > 0 {
		// retry 不为空
		var catchnodes []TaskCatchNode
		inputcatchnodes := data[StateFieldNames.Catch].([]interface{})
		for idx, catchnodedata := range inputcatchnodes {
			catchnodedata, ok := catchnodedata.(map[string]interface{})
			if !ok {
				return fmt.Errorf("'Catch' Node [ %d ] data should be map", idx+1)
			}
			node := DefaultCatchNode
			err = DecodeMapToStruct(catchnodedata, &node)
			if err != nil {
				return err
			}
			catchnodes = append(catchnodes, node)
		}
		body.Catch = catchnodes
	}

	err = body.Init()
	return err

}

// TaskTimeout Describe Task Timeout demand
type TaskTimeout struct {
	TaskTimeout      time.Time
	HeartBeatTimeout time.Time
}

// Validate ...
func (body *TaskBody) Validate() error {
	var err error
	err = myValidate.Struct(body)
	if err != nil {
		return err
	}
	_, err = ParseResource(body.Resource)
	if err != nil {
		return err
	}
	return err
}

// Init ...
func (body *TaskBody) Init() error {
	return nil
}

// TaskState ...
type TaskState struct {
	*BaseState `json:",inline"`
	*TaskBody  `json:",inline"`
}

func (t *TaskState) GetBaseState() *BaseState {
	return t.BaseState
}

// Init init task
func (t *TaskState) Init() error {

	err := t.TaskBody.Init()
	if err != nil {
		return err
	}
	return nil
}

// Validate ...
func (t *TaskState) Validate() error {
	return t.TaskBody.Validate()
}

// GetBone get bone
func (t *TaskState) GetBone() StateBone {
	bone := t.BaseState.GetBone()
	for _, catch := range t.TaskBody.Catch {
		bone.Next = append(bone.Next, catch.Next)
	}
	return bone
}

// GetTaskTimeout  return task timeout
func (t *TaskState) GetTaskTimeout() (TaskTimeout, error) {
	var tt TaskTimeout
	if t.TimeoutSeconds > 0 {
		tt.TaskTimeout = time.Now().Add(time.Duration(t.TimeoutSeconds) * time.Second)
	}
	if t.HeartbeatSeconds > 0 {
		tt.HeartBeatTimeout = time.Now().Add(time.Duration(t.HeartbeatSeconds) * time.Second)
	}
	return tt, nil
}

// GetNextState get task next state
// @input state input data
// @taskdata task send data
// return next state
func (t *TaskState) GetNextState(input any, taskdata TaskSendData) (*NextState, error) {
	var err error

	var nextstate NextState
	if taskdata.Success {
		// task submit success
		output, err := t.GetOutput(input, taskdata.Output)
		if err != nil {
			return nil, err
		}
		nextstate = NextState{
			Name:   t.Next,
			Output: output,
		}
		return &nextstate, nil
	}
	// task submit failed, find Retry/Catch strategy
	for index, retry := range t.TaskBody.Retry {

		match := HasIntersection(retry.ErrorEquals, taskdata.Errors)
		if !match {
			continue
		}

		if taskdata.Retry[index] >= int(retry.MaxAttempts) {
			// reach max attempts, continue to next retry
			continue
		}
		durationSecond := int(float64(retry.IntervalSeconds) * math.Pow(retry.BackoffRate, float64(taskdata.Retry[index])))
		var nextstate = NextState{
			// when retry, retry = true, name = ""
			Name:       "",
			Delay:      time.Second * time.Duration(durationSecond),
			Output:     nil,
			RetryIndex: index,
			Retry:      true,
		}
		return &nextstate, nil

	}
	// try catch node
	for _, catchnode := range t.TaskBody.Catch {
		match := HasIntersection(catchnode.ErrorEquals, taskdata.Errors)
		if !match {
			continue
		}
		output, err := t.GetOutputWithPath(input, taskdata.Output, catchnode.ResultPath, "$")
		if err != nil {
			return nil, err
		}
		nextstate = NextState{
			Name:   catchnode.Next,
			Output: output,
		}
		return &nextstate, nil
	}
	err = fmt.Errorf("can't match any strategy")
	return nil, err
}

func (t *TaskState) GetDefinition() (string, error) {
	data, err := ToString(t)
	return data, err
}

// HasIntersection  return  if x and y have common elements
func HasIntersection(x []string, y []string) bool {

	for _, yi := range y {
		if slices.Contains(x, yi) {
			return true
		}
	}
	return false

}
