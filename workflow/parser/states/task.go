package states

import (
	"fmt"
	"math"
	"time"
)

// TaskBody ...
type TaskBody struct {
	Resource         string `validate:"required,gt=0"`
	TimeoutSeconds   uint   `validate:"gte=0"`
	HeartbeatSeconds uint   `validate:"gte=0"`
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

// TaskTimeout Describe Task Timeout demand
type TaskTimeout struct {
	TaskTimeout      time.Duration
	HeartBeatTimeout time.Duration
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
	return body.Validate()
}

// Task ...
type Task struct {
	*BaseState
	*TaskBody
}

// Init init task
func (t *Task) Init() error {

	err := t.BaseState.Init()
	if err != nil {
		return err
	}
	err = t.TaskBody.Init()
	if err != nil {
		return err
	}
	return nil
}

// Validate ...
func (t *Task) Validate() error {
	return t.TaskBody.Validate()
}

// GetBone get bone
func (t *Task) GetBone() StateBone {
	bone := t.BaseState.GetBone()

	if len(t.TaskBody.Retry) == 0 {
		return bone
	}
	for _, catch := range t.TaskBody.Catch {
		bone.Next = append(bone.Next, catch.Next)
	}
	return bone
}

// GetTaskTimeout  return task timeout
func (t *Task) GetTaskTimeout() (TaskTimeout, error) {
	var tt TaskTimeout
	if t.TimeoutSeconds > 0 {
		tt.TaskTimeout = time.Duration(t.TimeoutSeconds) * time.Second
	}
	if t.HeartbeatSeconds > 0 {
		tt.HeartBeatTimeout = time.Duration(t.HeartbeatSeconds) * time.Second
	}
	return tt, nil
}

// GetNextState get task next state
// @input state input data
// @taskdata task send data
// return next state
func (t *Task) GetNextState(input interface{}, taskdata TaskSendData) (NextState, error) {
	var err error

	var nextstate NextState
	if taskdata.Success {
		// task submit success
		output, err := t.GetOutput(input, taskdata.Output)
		if err != nil {
			return nextstate, err
		}
		nextstate = NextState{
			Name:   t.Next,
			Output: output,
		}
		return nextstate, nil
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
		return nextstate, nil

	}
	// try catch node
	for _, catchnode := range t.TaskBody.Catch {
		match := HasIntersection(catchnode.ErrorEquals, taskdata.Errors)
		if !match {
			continue
		}
		output, err := t.GetOutputWithPath(input, taskdata.Output, catchnode.ResultPath, "$")
		if err != nil {
			return nextstate, err
		}
		nextstate = NextState{
			Name:   catchnode.Next,
			Output: output,
		}
		return nextstate, nil
	}
	err = fmt.Errorf("can't match any strategy")
	return nextstate, err
}

// HasIntersection  return  if x and y have common elements
func HasIntersection(x []string, y []string) bool {

	for _, yi := range y {
		for _, xi := range x {
			if yi == xi {
				return true
			}
		}
	}
	return false

}
