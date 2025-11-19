package executor

// struct.go defines the structures used in the workflow executor package.
import (
	"github.com/skyflow-workflow/skyflow_backbend/workflow/parser/states"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/po"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/repository/queue"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/vo"
)

// ExceptionData Step异常数据结构
type ExceptionData struct {
	Cause      string      `json:"cause"`       // Detail of the failure
	Error      string      `json:"error"`       // Error Code of the failure
	ErrorMatch []string    `json:"error_match"` // Error Code match list
	Output     interface{} `json:"output"`      // Output of state
	Extra      string      `json:"extra"`       // Extra information about exception
}

// StartExecutionResponse function StartExecution Return struct
type StartExecutionResponse struct {
	Data     *po.Execution
	Events   []vo.ExecutionEvent
	Messages []queue.InnerMessageBody
}

// StopExecutionResponse function StopExecution Return struct
type StopExecutionResponse struct {
	Events []vo.ExecutionEvent
}

// StepWakeupMessage  在messagequeue 中表达此次节点唤醒的信息
type StepWakeupMessage struct {
	// step execute counter when send message
	ExecuteCount int    `json:"execute_count"`
	Token        string `json:"token"`
}

type StepExecuteMessage struct {
	// if block step execute
	Block bool `json:"block"`
}

// FindNextStep  for find next step
type FindNextStep struct {
	Name    string
	GroupID int
}

// NextStep  for find next step
type NextStep struct {
	states.NextState
	GroupID int // state group id
}

// GroupStatus
type GroupStatus struct {
	GroupID int
	Success bool
	Status  string
}

// ActivityTaskData   data  attach to activity task
type ActivityTaskData struct {
	TimeoutSeconds   int `json:"timeout_seconds"`   // execute timeout in seconds
	HeartbeatSeconds int `json:"heartbeat_seconds"` // heartbeat timeout in seconds
	HeartbeatCount   int `json:"heartbeat_count"`   // current heartbeat count
}

// ExecutionData Execution Data Field
type ExecutionData struct {
	TaskData string `json:"task_data"` // Task Data
}

// TaskFailure task state failed struct
type TaskFailure struct {
	Error string
	Cause string
}

// ChangeStateStatusRequest 修改State状态请求
type ChangeStepStatusRequest struct {
	StateStatus _StepStatusType
	GroupStatus _StepStatusType
}

// InsertStateMachineOption 插入状态机选项
type InsertStateMachineOption struct {
	// StartGroupID 可以使用的状态组ID
	StartGroupID int
	//StartDeIndex 可以使用的开始状态递减索引
	StartDeIndex int
	// StartDepth 可以使用的开始状态深度
	StartDepth int
}

// InsertStateMachineResponse 插入状态机响应
type InsertStateMachineResponse struct {
	// 已插入的状态机 第一个 StepID
	StartStepID int
	// 已插入的状态机 第一个 GroupID
	StartGroupID int
	// 已插入的状态机 最大的 GroupID
	MaxGroupID int
	// 已插入的状态机 最小的 Deindex
	MinDeIndex int
}

var DefaultStepExecuteMessage = StepExecuteMessage{
	Block: false,
}
