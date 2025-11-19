package executor

/*
 * @Author: mumangtao@gmail.com
 * @Date: 2020-08-08 22:12:01
 * @Last Modified by: mumangtao@gmail.com
 * @Last Modified time: 2020-08-08 22:18:13
 */

import (
	"time"

	"github.com/skyflow-workflow/skyflow_backend/workflow/vo"
)

// ExecutionEvent event结构
type ExecutionEvent struct {
	ID            int
	ExecutionID   int
	ExecutionUUID string
	ExecutionURI  string
	StepName      string
	StepID        int
	EventType     string
	NanoSeconds   string
	Data          interface{}
	StartTime     time.Time
	FinishTime    time.Time
}

// EventContent Event Content Base Struct
type EventContent struct {
}

// EventContentPrefix 事件内容的前缀匹配
var EventContentPrefix = "EventContent_"

// EventContent_ExecutionStart Event Content for Execution start
type EventContent_ExecutionStart struct {
	Input string `json:"input"`
	URI   string `json:"uri"`
}

// EventContent_ExecutionCreated  Event Content for Execution start
type EventContent_ExecutionCreated struct {
	Input string `json:"input"`
}

// EventContent_ExecutionSucceeded  Event Content for Execution start
type EventContent_ExecutionSucceeded struct {
	Output string `json:"output"`
}

// EventContent_ExecutionFailed  Event Content for Execution start
type EventContent_ExecutionFailed struct {
	EventType string `json:"event_type"`
	Error     string `json:"error"`
	Cause     string `json:"cause"`
}

// EventContent_ExecutionSuspend  Event Content for Execution Suspend
type EventContent_ExecutionSuspend struct {
}

// EventContent_ExecutionBlocked  Event Content for Execution Suspend
type EventContent_ExecutionBlocked struct {
}

// EventContent_ExecutionContinue  Event Content for Execution Suspend
type EventContent_ExecutionContinue struct {
}

// EventContent_ExecutionRetry  Event Content for Execution start
type EventContent_ExecutionRetry struct {
}

// EventContent_ExecutionInfoModified  Event Content for Execution start
type EventContent_ExecutionInfoModified struct {
	Field  string      `json:"field"`
	Before interface{} `json:"before"  `
	After  interface{} `json:"after" `
}

// EventContent_ExecutionAbort  Event Content for Execution start
type EventContent_ExecutionAbort struct {
	Error string `json:"error"`
	Cause string `json:"cause"`
}

// EventContent_TaskSubmitFailed  Event Content for Execution start
type EventContent_TaskSubmitFailed struct {
	Error       string         `json:"error"`
	Cause       string         `json:"cause"`
	RequestInfo vo.RequestInfo `json:"requestinfo"`
}

// EventContent_TaskSubmitted  Event Content for Execution start
type EventContent_TaskSubmitted struct {
	Result      interface{}    `json:"result"`
	RequestInfo vo.RequestInfo `json:"requestinfo"`
}

// EventContent_TaskSendRefer  Event Content for Execution start
type EventContent_TaskSendReference struct {
	Reference   string         `json:"reference"`
	RequestInfo vo.RequestInfo `json:"requestinfo"`
}

// EventContent_StateEntered  Event content for State Entered
type EventContent_StateEntered struct {
	Input string `json:"input"`
}

// EventContent_StateInit  Event content for State Entered
type EventContent_StateInit struct {
	Action         string
	ExecutionIndex int
	ExecutionCount int
}

// EventContent_StateExited  Event content for State Exited
type EventContent_StateExited struct {
	Output interface{} `json:"output"`
}

// EventContent_StateFailed   state faild
type EventContent_StateFailed struct {
	EventType string `json:"event_type"`
	Error     string `json:"error"`
	Cause     string `json:"cause"`
}

// EventContent_StateInfoModified  Event content for State Exited
type EventContent_StateInfoModified struct {
	Field  string      `json:"field"`
	Before interface{} `json:"before"  `
	After  interface{} `json:"after" `
}

// EventContent_TaskStateExited  Event content for State Exited
type EventContent_TaskStateExited struct {
	Output interface{} `json:"output"`
}

// EventContent_TaskStateRetry  Event content for WaitStateExecuted
type EventContent_TaskStateRetry struct {
	RetryIndex int       `json:"retry_index"`
	Error      string    `json:"error"`
	RetryTime  time.Time `json:"retry_time"`
}

// EventContent_PassStateExecuted  Event content for State Exited
type EventContent_PassStateExecuted struct {
	Result interface{} `json:"result"`
}

// EventContent_WaitStateExecuted  Event content for WaitStateExecuted
type EventContent_WaitStateExecuted struct {
	ExecuteCount int       `json:"execute_count"`
	WakeupTime   time.Time `json:"wakeup_time"`
}

// EventContent_WaitStateWakeup EventContent_WaitStateWakeup
type EventContent_WaitStateWakeup struct {
}

// EventContent_ChoiceStateExecuted  Event content for WaitStateExecuted
type EventContent_ChoiceStateExecuted struct {
	Next   string      `json:"next"`
	Output interface{} `json:"output"`
}

// EventContent_ParallelStateExecuted  Event content for WaitStateExecuted
type EventContent_ParallelStateExecuted struct {
}
type EventContent_StepGroupExecuted struct {
	Index int `json:"index"`
}

// EventContent_MapIterationStarted  Event Content for map item start
type EventContent_MapIterationStarted struct {
	Index int    `json:"index"`
	Name  string `json:"name"`
}

// EventContent_MapIterationSucceeded   Event Content for map item succeed
type EventContent_MapIterationSucceeded struct {
	Index  int         `json:"index"`
	Output interface{} `json:"output"`
}

// EventContent_ParallelBranchSucceeded   Event Content for map item succeed
type EventContent_ParallelBranchSucceeded struct {
	Index  int         `json:"index"`
	Output interface{} `json:"output"`
}

// EventContent_StepGroupSucceed EventContent_StepGroupSucceed
type EventContent_StepGroupSucceed struct {
	Index  int
	Output interface{}
}

// EventContent_StepGroupSuspend EventContent_StepGroupSuspend
type EventContent_StepGroupSuspend struct {
	Index int
}

// EventContent_StepGroupBlocked EventContent_StepGroupBlocked
type EventContent_StepGroupBlocked struct {
	Index int
}

// EventContent_StepGroupRecovery EventContent_StepGroupRecovery
type EventContent_StepGroupRecovery struct {
	Index int
}

// EventContent_StepGroupFailed EventContent_StepGroupFailed
type EventContent_StepGroupFailed struct {
	Index int
	Error string
	Cause string
}

// EventContent_ActivityScheduled   activity scheduled
type EventContent_ActivityScheduled struct {
	Input       interface{}      `json:"input"`
	Resource    string           `json:"resource"`
	Timeout     ActivityTaskData `json:"timeout"`
	RequestInfo vo.RequestInfo   `json:"requestinfo"`
}

// EventContent_TaskInitialized   activity scheduled
type EventContent_TaskInitialized struct {
	TimeoutSeconds   int
	HeartBeatSeconds int
	Input            interface{}
	Resource         string
}

// EventContent_TaskSendHeartbeatd   activity scheduled
type EventContent_TaskSendHeartbeat struct {
	Message     string         `json:"message"`
	RequestInfo vo.RequestInfo `json:"requestinfo"`
}

// EventContent_TaskBlocked   activity scheduled
type EventContent_TaskBlocked struct {
}

// EventContent_TaskRecovery   activity scheduled
type EventContent_TaskRecovery struct {
}

// EventContent_StepRetry   activity scheduled
type EventContent_StepRetry struct {
	ExecuteCount int
}

// EventContent_SendStepFailed   step failedevent
type EventContent_SendStepFailed struct {
}

// EventContent_StoreTaskData   store task data event
type EventContent_StoreTaskData struct {
	Data        string         `json:"data"`
	RequestInfo vo.RequestInfo `json:"requestinfo"`
}

// EventContent_StoreTaskData   store task data event
type EventContent_UnblockTask struct {
	RequestInfo vo.RequestInfo `json:"requestinfo"`
}

// EventContent_StoreTaskData   store task data event
type EventContent_UnblockExecution struct {
}

// EventContent_RedoStep   step failedevent
type EventContent_RedoStep struct {
	ExecuteCount int
}

// EventContent_StepSkip   activity scheduled
type EventContent_StepSkip struct {
	NextStepName string
	Output       string
	RequestInfo  vo.RequestInfo
}

// EventContent_SkipBlockedTask   activity scheduled
type EventContent_SkipBlockedTask struct {
	NextStepName string
	Output       string
	RequestInfo  vo.RequestInfo
}

// EventContent_TaskPush   activity scheduled
type EventContent_TaskPush struct {
	Input            interface{}
	Resource         string
	Timeout          time.Time
	HeartBeatTimeout time.Time
	ExecuteCount     int
}

// EventContent_SucceedStateExited   activity scheduled
type EventContent_SucceedStateExited struct {
	Output interface{} `json:"output"`
}

// EventContent_SuspendStateExecuted  Event content for
type EventContent_SuspendStepExecuted struct {
	// ResumeTimeout time.Time `json:"resume_timeout"`
}

// EventContent_SuspendStateTimeout EventContent_SuspendStepResumeTimeout
type EventContent_SuspendStepResume struct {
}

// EventContent_SuspendStateTimeout EventContent_SuspendStepResumeTimeout
type EventContent_SuspendStepResumeTimeout struct {
}

type EventContent_ProcessEventFailed struct {
	Error     string `json:"error"`
	EventType string `json:"event_type"`
}
