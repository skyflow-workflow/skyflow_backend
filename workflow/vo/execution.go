package vo

import (
	"gopkg.mihoyo.com/plat/cloudflow/pkg/paging"
	"gopkg.mihoyo.com/plat/cloudflow/workflow/po"
	"gopkg.mihoyo.com/plat/cloudflow/workflow/repository/queue"
)

// StartExecutionRequest 新建Execution 的请求
type StartExecutionRequest struct {
	ExecutionUUID      string        `json:"execution_uuid" `      //关联的 uuid
	WorkflowURI        string        `json:"workflow_uri" `        //statemachine uri
	WorkflowDefinition string        `json:"workflow_definition" ` // statemachine content
	Title              string        `json:"title"`                // execution title
	Data               ExecutionData // 执行数据
	Input              string        `json:"input" xorm:"MEDIUMTEXT"` //输入
}

// ExecutionData Execution Data Field
type ExecutionData struct {
	TaskData string `json:"task_data"` // Task Data
}

// ResponseCreateExecution  createexecution 返回的struct
type StartExecutionResponse struct {
	Data     po.Execution
	Events   []ExecutionEvent
	Messages []queue.InnerMessage
}

// RestartExecutionRequest 重新创建 Execution
type RestartExecutionRequest struct {
	// 当前需要重新发起的execution uuid
	UUID  string `json:"uuid" jpath:"uuid" post:"required notzero"`
	Cause string `json:"cause" jpath:"cause"`
	Error string `json:"error" jpath:"error"`
	// 是否关闭当前任务
	Abort bool `json:"abort" jpath:"cause"`
}

// StopExecution 终止 Execution 请求结构
type StopExecutionRequest struct {
	ExecutionUUID string `json:"execution_uuid" jpath:"uuid" post:"required notzero"`
	Cause         string `json:"cause" jpath:"cause"`
	Error         string `json:"error" jpath:"error"`
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

type GetActivityTaskRequest struct {
	ActivityURI string
}

type GetActivityTaskResponse struct {
	Step             *po.Step
	Execution        *po.Execution
	Input            string
	TaskToken        string
	TimeoutSeconds   int
	HeartbeatSeconds int
}

type SendTaskSuccessRequest struct {
	TaskToken string
	Output    string
}

type SendTaskFailureRequest struct {
	TaskToken string
	Error     string
	Cause     string
}

type SendTaskHeartbeatRequest struct {
	TaskToken string
	Message   string
}
type SendTaskReferenceRequest struct {
	TaskToken string
	Title     string
	URL       string
}

type SendStepSkipRequest struct {
	StepID       int
	NextStepName string
	Output       string
}

type DescribeExecutionReqeust struct {
	ExecutionUUID string
}
type DescribeExecutionBoneResponse struct {
	Data po.Execution
	Bone interface{}
}

type DescribeStepReqeust struct {
	StepID int64
}

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

type ListExecutionEventsRequest struct {
	ExecutionUUID string
	ExecutionID   int
	PageRequest   paging.PageRequest
}

type ListStepEventsRequest struct {
	StepID      int
	PageRequest paging.PageRequest
}

type ListExecutionEventsResponse struct {
	Events       []po.ExecutionEvent
	PageResponse paging.PageResponse
}

type Reference struct {
	Title string
	URL   string
}

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

type SkipBlokcedTaskRequest struct {
	StepID int
}

type UnblockTaskRequest struct {
	StepID int
}

type IdentifyUser struct{}
