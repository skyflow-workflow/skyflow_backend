package vo

import (
	"github.com/mmtbak/microlibrary/paging"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/po"
)

// StartExecutionRequest 新建Execution 的请求
type StartExecutionRequest struct {
	ExecutionUUID      string        `json:"execution_uuid" `      //关联的 uuid
	WorkflowURI        string        `json:"workflow_uri" `        //statemachine uri
	WorkflowDefinition string        `json:"workflow_definition" ` // statemachine content
	Title              string        `json:"title"`                // execution title
	Data               ExecutionData // 执行数据
	Input              string        `json:"input"` //输入
}

// ExecutionData Execution Data Field
type ExecutionData struct {
	TaskData string `json:"task_data"` // Task Data
}

// StartExecutionResponse ResponseCreateExecution  createexecution 返回的struct
type StartExecutionResponse struct {
	Data   po.Execution
	Events []ExecutionEvent
}

// RestartExecutionRequest 重新创建 Execution
type RestartExecutionRequest struct {
	// 当前需要重新发起的execution uuid
	UUID  string `json:"uuid" post:"required notzero"`
	Cause string `json:"cause" `
	Error string `json:"error" `
	// 是否关闭当前任务
	Abort bool `json:"abort"`
}

// StopExecutionRequest StopExecution 终止 Execution 请求结构
type StopExecutionRequest struct {
	ExecutionUUID string `json:"execution_uuid"  post:"required notzero"`
	Cause         string `json:"cause"`
	Error         string `json:"error"`
}

// StopExecutionResponse stopexecution 返回的struct
type StopExecutionResponse struct {
	Events []ExecutionEvent
}

// StateWakeupMessage  在messagequeue 中表达此次节点唤醒的信息
type StateWakeupMessage struct {
	Counter int    `json:"counter"`
	Token   string `json:"token"`
}

// GetActivityTaskRequest ...
type GetActivityTaskRequest struct {
	ActivityURI string
}

// GetActivityTaskResponse ...
type GetActivityTaskResponse struct {
	Step             *po.State
	Execution        *po.Execution
	Input            string
	TaskToken        string
	TimeoutSeconds   int
	HeartbeatSeconds int
}

// SendTaskSuccessRequest ...
type SendTaskSuccessRequest struct {
	TaskToken string
	Output    string
}

// SendTaskFailureRequest ...
type SendTaskFailureRequest struct {
	TaskToken string
	Error     string
	Cause     string
}

// SendTaskHeartbeatRequest ...
type SendTaskHeartbeatRequest struct {
	TaskToken string
	Message   string
}

// SendTaskReferenceRequest ...
type SendTaskReferenceRequest struct {
	TaskToken string
	Title     string
	URL       string
}

// SendStepSkipRequest ...
type SendStepSkipRequest struct {
	StepID       int
	NextStepName string
	Output       string
}

// DescribeExecutionRequest ...
type DescribeExecutionRequest struct {
	ExecutionUUID string
}

// DescribeExecutionBoneResponse ...
type DescribeExecutionBoneResponse struct {
	Data po.Execution
	Bone interface{}
}

// DescribeStepRequest ...
type DescribeStepRequest struct {
	StepID int64
}

// ExecutionName ...
type ExecutionName struct {
	URI   string
	Name  string
	Count int
}

// ListExecutionsRequest 查询Execution 请求
type ListExecutionsRequest struct {
	WorkflowURI    string
	Title          string
	Status         string
	ExecutionUUIDs []string
	PageRequest    paging.PageRequest
}

// ListExecutionsResponse 查询Execution 返回
type ListExecutionsResponse struct {
	Executions   []po.Execution
	PageResponse paging.PageResponse
}

// ListExecutionEventsRequest ...
type ListExecutionEventsRequest struct {
	ExecutionUUID string
	ExecutionID   int
	PageRequest   paging.PageRequest
}

// ListStepEventsRequest ...
type ListStepEventsRequest struct {
	StepID      int
	PageRequest paging.PageRequest
}

// ListExecutionEventsResponse ...
type ListExecutionEventsResponse struct {
	Events       []po.ExecutionEvent
	PageResponse paging.PageResponse
}

// Reference ...
type Reference struct {
	Title string
	URL   string
}

// StoreTaskDataRequest ...
type StoreTaskDataRequest struct {
	TaskToken string
	Data      string
}

// DescribeExecutionBoneRequest pipeline bone
type DescribeExecutionBoneRequest struct {
	ExecutionID   int
	ExecutionUUID string
	PipelineMode  bool
}

// SkipBlockedTaskRequest ...
type SkipBlockedTaskRequest struct {
	StepID int
}

// UnblockTaskRequest ...
type UnblockTaskRequest struct {
	StepID int
}

// IdentifyUser ...
type IdentifyUser struct{}
