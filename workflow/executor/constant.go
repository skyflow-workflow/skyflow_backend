package executor

import (
	"github.com/skyflow-workflow/skyflow_backbend/workflow/vo"
)

// MaxStepQueryLimit 最大步骤数查询限制
var MaxStepQueryLimit = 100

// TimeFormat 时间格式化格式
var TimeFormat = "2006-01-02 15:04:05"

// DateFormat 时间格式化格式
var DateFormat = "2006-01-02"

type _ExecutionStatusType string

// ExecutionStatus execution 状态列表
var ExecutionStatus = struct {
	Created    _ExecutionStatusType
	Running    _ExecutionStatusType
	Failed     _ExecutionStatusType
	Success    _ExecutionStatusType
	Aborted    _ExecutionStatusType
	Suspending _ExecutionStatusType
	Blocked    _ExecutionStatusType
}{
	Created:    "Created",    // 刚刚创建
	Running:    "Running",    // 运行中
	Failed:     "Failed",     // 运行失败
	Success:    "Success",    // 运行成功
	Aborted:    "Aborted",    // 任务取消
	Suspending: "Suspending", // 任务暂停
	Blocked:    "Blocked",    // 任务阻塞
}

// ExecutionStatusWeight execution 状态权重
// 状态合并的时候，需要去计算权重最大的状态
var ExecutionStatusWeight = map[_ExecutionStatusType]int{
	ExecutionStatus.Created:    1,
	ExecutionStatus.Success:    2,
	ExecutionStatus.Blocked:    3,
	ExecutionStatus.Suspending: 4,
	ExecutionStatus.Failed:     5,
	ExecutionStatus.Running:    6,
	ExecutionStatus.Aborted:    10,
}

type _StepStatusType string

// StepStatus task 状态列表
var StepStatus = struct {
	Created       _StepStatusType
	WaitInit      _StepStatusType
	Initialize    _StepStatusType
	Wait          _StepStatusType
	Wakeup        _StepStatusType
	Sleep         _StepStatusType
	Running       _StepStatusType
	TaskSubmitted _StepStatusType
	Failed        _StepStatusType
	Success       _StepStatusType
	Cancelled     _StepStatusType
	Skip          _StepStatusType
	// WaitProcess   _TaskStatusType
	Virtual    _StepStatusType
	Suspending _StepStatusType
	Blocked    _StepStatusType
	Unblocking _StepStatusType
}{
	Created:       "Created",       // 刚刚创建的，未初始化状态，
	WaitInit:      "WaitInit",      // 等待初始化 ，对于即将被调度的节点，处于该状态表示上个节点处理完成，这个节点待执行，表示两个状态节点的交接状态。
	Initialize:    "Initialize",    // 已经初始化完成 : 生成Data 字段和 activity 字段
	Sleep:         "Sleep",         // 触发重试后，在需要延时执行的情况下， 状态处于 sleep 状态
	Wait:          "Wait",          // Task等待业务worker调度被执行阶段，此时状态节点准备被调用。
	Wakeup:        "WakeUp",        // Task State 从重试中醒来的阶段，很短，很快会变成 Wait
	Running:       "Running",       // 状态节点运行中
	TaskSubmitted: "TaskSubmitted", //Task节点提交了结果，正常/异常
	Failed:        "Failed",        // 状态节点 运行失败
	Success:       "Success",       //状态节点运行成功
	Cancelled:     "Cancelled",     //状态节点取消运行
	Skip:          "Skip",          //跳过
	// WaitProcess:   "WaitProcess",   // 手动操作后，等待处理
	Virtual:    "Virtual",
	Suspending: "Suspending",
	Blocked:    "Blocked",
	Unblocking: "Unblocking",
}

// ExecutionEventType execution时间类型列表
var ExecutionEventType = struct {
	ExecutionCreated           string
	ExecutionStarted           string
	ExecutionFailed            string
	ExecutionAborted           string
	ExecutionInfoModified      string
	ExecutionSucceeded         string
	ExecutionAcquireLockFailed string
	ExecutionRetry             string
	StateEntered               string
	StateExited                string
	StateSucceed               string
	StateFailed                string
	TaskStateEntered           string
	TaskStateExited            string
	PassStateEntered           string
	PassStateExecuted          string
	PassStateExited            string
	ChoiceStateEntered         string
	ChoiceStateExecuted        string
	ChoiceStateExited          string
	WaitStateEntered           string
	WaitStateExited            string
	SucceedStateEntered        string
	SucceedStateExited         string
	FailStateEntered           string
	FailStateExited            string
	ParallelStateEntered       string
	ParallelStateExited        string
	MapStateEntered            string
	MapStateExited             string
	MapIterationStarted        string
	MapIterationSucceeded      string
	StepGroupSucceed           string
	StepGroupFailed            string
	StateInfoModified          string
	StateManualRetry           string
	StateManualSkip            string
	StateManualFailed          string
	StateAcquireLockFailed     string
	StateBlocked               string
	StateInit                  string
	StateInitFailed            string
	StateExecuteFailed         string
	StateWakeup                string
	TaskTimeout                string
	TaskInitialized            string
	TaskSubmitted              string
	TaskSubmitFailed           string
	TaskSendHeartbeat          string
	TaskSendRefer              string
	TaskPushed                 string
	ChoiceFailed               string
	ActivityScheduled          string
	ResourceScheduled          string
	ParallelStateFailed        string
	ParallelStateSucceed       string
}{
	ExecutionCreated:           "ExecutionCreated",
	ExecutionStarted:           "ExecutionStarted",
	ExecutionFailed:            "ExecutionFailed",
	ExecutionAborted:           "ExecutionAborted",
	ExecutionInfoModified:      "ExecutionInfoModified",
	ExecutionSucceeded:         "ExecutionSucceeded",
	ExecutionAcquireLockFailed: "ExecutionAcquireLockFailed",
	ExecutionRetry:             "ExecutionRetry",

	StateEntered:           "StateEntered",
	StateExited:            "StateExited",
	StateSucceed:           "StateSucceed",
	StateFailed:            "StateFailed",
	TaskStateEntered:       "TaskStateEntered",
	TaskStateExited:        "TaskStateExited",
	PassStateEntered:       "PassStateEntered",
	PassStateExecuted:      "PassStateExecuted",
	PassStateExited:        "PassStateExited",
	ChoiceStateEntered:     "ChoiceStateEntered",
	ChoiceStateExecuted:    "ChoiceStateExecuted",
	ChoiceStateExited:      "ChoiceStateExited",
	WaitStateEntered:       "WaitStateEntered",
	WaitStateExited:        "WaitStateExited",
	SucceedStateEntered:    "SucceedStateEntered",
	SucceedStateExited:     "SucceedStateExited",
	FailStateEntered:       "FailStateEntered",
	FailStateExited:        "FailStateExited",
	ParallelStateEntered:   "ParallelStateEntered",
	ParallelStateExited:    "ParallelStateExited",
	MapStateEntered:        "MapStateEntered",
	MapStateExited:         "MapStateExited",
	MapIterationStarted:    "MapIterationStarted",
	MapIterationSucceeded:  "MapIterationSucceeded",
	StepGroupSucceed:       "StepGroupSucceed",
	StepGroupFailed:        "StepGroupFailed",
	StateInfoModified:      "StateInfoModified",
	StateManualRetry:       "StateManualRetry",
	StateManualSkip:        "StateManualSkip",
	StateManualFailed:      "StateManualFailed",
	StateAcquireLockFailed: "StateAcquireLock",
	StateBlocked:           "StateBlocked",
	StateInit:              "StateInit",
	StateInitFailed:        "StateInitFailed",
	StateExecuteFailed:     "StateExecuteFailed",
	StateWakeup:            "StateWakeup",
	TaskTimeout:            "TaskTimeout",
	TaskSubmitted:          "TaskSubmitted",
	TaskPushed:             "TaskPushed",
	TaskSendHeartbeat:      "TaskSendHeartbeat",
	TaskInitialized:        "TaskInitialized",
	TaskSubmitFailed:       "TaskSubmitFailed",
	TaskSendRefer:          "TaskSendRefer",
	ChoiceFailed:           "ChoiceFailed",
	ActivityScheduled:      "ActivityScheduled",
	ResourceScheduled:      "ResourceScheduled",
	ParallelStateFailed:    "ParallelStateFailed",
	ParallelStateSucceed:   "ParallelStateSucceed",
}

// StateTimerEventType 状态计时器 事件类型
var StateTimerEventType = struct {
	TaskSleep     string
	TaskHeartBeat string
	TaskTimeup    string
	WaitTimeup    string
}{
	TaskSleep:     "TaskSleep",
	TaskHeartBeat: "TaskHeartBeat",
	TaskTimeup:    "TaskTimeup",
	WaitTimeup:    "WaitTimeup",
}

// StateTimerStatus 状态计时器的状态
var StateTimerStatus = struct {
	Created   string
	Running   string
	Process   string
	Closed    string
	Cancelled string
}{

	Created:   "Created",
	Running:   "Running",
	Process:   "Process",
	Closed:    "Closed",
	Cancelled: "Cancelled",
}

// ActivityStatus 活动的状态
var ActivityStatus = struct {
	Enable  string
	Disable string
}{
	Enable:  "Enable",
	Disable: "Disable",
}

// StateMachineStatus statemachine 状态
var StateMachineStatus = struct {
	Enable  string
	Disable string
}{
	Enable:  "Enable",
	Disable: "Disable",
}

// DataModifyField 可以修改的字段
var DataModifyField = struct {
	Input  string
	Status string
}{
	Input:  "Input",
	Status: "Status",
}

// InnerStateError 内部状态错误类型
var InnerStateError = struct {
	StatesALL          string
	StatesTimeout      string
	StatesAbortTimeout string
	StatesTaskFailed   string
}{
	StatesALL:          "States.ALL",
	StatesTimeout:      "States.Timeout",
	StatesAbortTimeout: "State.AbortTimeout",
	StatesTaskFailed:   "States.TaskFailed",
}

// MessageType Message Type
var MessageType = struct {
	ExecutionInit string // Execution 初始化
	// ExecutionStop    string // Execution 停止
	// ExecutionPause   string // Execution 暂停
	ExecutionFailed      string // Execution 失败
	ExecutionTimout      string // Execution 执行超时
	ExecutionAbortTimout string // Execution 终止超时
	ExecutionSucceed     string // Execution 成功
	ExecutionWakeup      string // Execution 唤醒 等待初始化的State
	StateNewTurn         string // State 新到处理
	StateExecute         string // State 执行
	TaskHeartBeatTimeup  string // Task 心跳时间到
	TaskTimeoutTimeup    string // Task 超时时间到
	TaskStateSend        string // Task 发送执行结果
	TaskStateWakeup      string // Task 唤醒执行
	WaitStateWakeup      string // Wait 唤醒执行
	SuspendTimeout       string // Suspend 超时
	StepGroupSucceed     string // Execution Group 执行成功
	StepGroupFailed      string // Execution Group 执行失败
	StepGroupSuspend     string // Execution Group 暂停
	StepGroupBlocked     string // Execution Group
	ParallelStateFinish  string //ParallelState 执行完成
	MapStateFinish       string //MapState 执行完成
	FindNextStep         string // 查找下一个节点
	ReportStepSuspend    string // 查找下一个节点
	ReportStepBlocked    string // 查找下一个节点

	// 人工干预消息
	StateRetry           string
	StateSkip            string
	StateFailed          string
	ManualStateSkip      string
	ManualStateFailed    string
	ManualExecutionStop  string //手动停止任务
	ManualExecutionPause string // 手动暂停任务
}{
	ExecutionInit: "ExecutionInit",
	// ExecutionStop:    "ExecutionStop",
	// ExecutionPause:   "ExecutionPause",
	ExecutionFailed:      "ExecutionFailed",
	ExecutionTimout:      "ExecutionTimout",      // Execution 执行超时
	ExecutionAbortTimout: "ExecutionAbortTimout", // Execution 终止超时
	ExecutionSucceed:     "ExecutionSucceed",
	ExecutionWakeup:      "ExecutionWakeup",

	StateNewTurn:        "StateNewTurn",
	StateExecute:        "StateExecute",
	TaskHeartBeatTimeup: "TaskHeartBeatTimeup",
	TaskTimeoutTimeup:   "TaskTimeoutTimeup",
	TaskStateSend:       "TaskStateSend",
	TaskStateWakeup:     "TaskStateWakeup",
	WaitStateWakeup:     "WaitStateWakeup",
	StepGroupSucceed:    "StepGroupSucceed",
	StepGroupFailed:     "StepGroupFailed",
	StepGroupSuspend:    "StepGroupSuspend",
	StepGroupBlocked:    "StepGroupBlocked",
	ParallelStateFinish: "ParallelStateFinish",
	MapStateFinish:      "MapStateFinish",
	FindNextStep:        "FindNextStep",
	ReportStepSuspend:   "ReportStepSuspend",
	ReportStepBlocked:   "ReportStepBlocked",
	//人工干预消息
	StateRetry:           "StateRetry",
	StateFailed:          "StateFailed",
	StateSkip:            "StateSkip",
	ManualStateSkip:      "ManualStateSkip",
	ManualStateFailed:    "ManualStateFailed",
	ManualExecutionStop:  "ManualExecutionStop",  //手动停止任务
	ManualExecutionPause: "ManualExecutionPause", // 手动暂停任务
}

// MessageStatus 消息状态
var MessageStatus = struct {
	Created        string
	Read           string
	ProcessSucceed string
	ProcessFailed  string
	Ignore         string
}{
	Created:        "Created",
	Read:           "Read",
	ProcessSucceed: "ProcessSucceed",
	ProcessFailed:  "ProcessFailed",
	Ignore:         "Ignore",
}

// EmbedAccessType 内嵌访问类型
var EmbedAccessType = struct {
	ExecutionDetail   string
	ExecutionClassify string
}{
	ExecutionDetail:   "ExecutionDetail",
	ExecutionClassify: "ExecutionClassify",
}

// ExecutionEventStateName  StateName for Execution event
var ExecutionEventStateName = struct {
	Start string
	End   string
}{
	Start: "START",
	End:   "END",
}

// TimeSeriesNodeType  时序Node 的类型列表
var TimeSeriesNodeType = struct {
	Initialize string
	Delay      string
	State      string
}{
	Initialize: "Initialize", //初始化
	Delay:      "Delay",      // 执行延时
	State:      "State",      // 状态执行
}

// StateEventStatusCheckMap  事件的前置状态检查字典
// 不同的事件对状态有要求，通过检查状态， 确认消息是否可消费
var StateEventEventCheckStatus = map[string][]string{
	//normal message
	MessageType.StateNewTurn:        {string(StepStatus.WaitInit)},
	MessageType.FindNextStep:        {string(StepStatus.Skip), string(StepStatus.Success)},
	MessageType.StateExecute:        {string(StepStatus.Initialize), string(StepStatus.Wakeup), string(StepStatus.Unblocking)},
	MessageType.TaskStateSend:       {string(StepStatus.Running)},
	MessageType.TaskStateWakeup:     {string(StepStatus.Sleep)},
	MessageType.WaitStateWakeup:     {string(StepStatus.Running)},
	MessageType.ParallelStateFinish: {string(StepStatus.Running)},
	MessageType.MapStateFinish:      {string(StepStatus.Running)},
	MessageType.StepGroupSucceed:    {string(ExecutionStatus.Running)},
	MessageType.StepGroupFailed:     {string(ExecutionStatus.Running)},
	MessageType.StepGroupSuspend:    {string(ExecutionStatus.Running)},
	MessageType.TaskTimeoutTimeup:   {string(StepStatus.Running)},
	// manual message
	// MessageType.StateRetry: {string(TaskStatus.Failed), string(TaskStatus.WaitProcess)},
	// MessageType.StateSkip:  {string(TaskStatus.Failed), string(TaskStatus.WaitProcess)},
	MessageType.StateFailed: {string(StepStatus.Initialize), string(StepStatus.Running),
		string(StepStatus.Sleep), string(StepStatus.Wait)},
}

// StateEventStatusCheckMap  事件的前置状态检查字典
// 不同的事件对状态有要求，通过检查状态， 确认消息是否可消费
var ExecutionEventCheckStatus = map[string][]string{
	MessageType.ExecutionInit:        {string(ExecutionStatus.Created)},
	MessageType.ExecutionFailed:      {string(ExecutionStatus.Running)},
	MessageType.ExecutionAbortTimout: {string(ExecutionStatus.Running)},
	MessageType.ExecutionSucceed:     {string(ExecutionStatus.Running)},
	MessageType.ExecutionTimout:      {string(ExecutionStatus.Running)},
}

// SchedularRequestInfo schedular的默认调度请求信息
var SchedularRequestInfo = vo.RequestInfo{
	RequestType:   "Schedular",
	RemoteAddress: "",
}
