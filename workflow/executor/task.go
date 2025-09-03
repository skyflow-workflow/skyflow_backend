package executor

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/mmtbak/microlibrary/rdb"
	"github.com/skyflow-workflow/skyflow_backbend/pkg/toolkit"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/po"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/repository/queue"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/vo"
	"gorm.io/gorm/clause"

	"github.com/skyflow-workflow/skyflow_backbend/workflow/parser/states"
)

// Task Execution Task State
type Task struct {
	*ExecutionStep // 继承 ExecutionStep
	TaskState      *states.Task
}

// TaskRunningStatus  Task运营状态
var TaskRunningStatus = struct {
	InitData      string
	WaitResult    string
	SendSuccess   string
	SendFailed    string
	SendHeartBeat string
	Timeout       string
	Finish        string
}{
	InitData:      "InitData",
	WaitResult:    "WaitResult",
	SendSuccess:   "SendSuccess",
	SendFailed:    "SendFailed",
	SendHeartBeat: "SendHeartBeat",
	Timeout:       "Timeout",
	Finish:        "Finish",
}

// TaskInnerData Task State Inner Data Struct
type TaskInnerData struct {
	Retry            []int  `json:"retry"`
	Status           string `json:"status"`
	HeartbeatCount   int    `json:"heartbeat_count"`
	HeartbeatSeconds int    `json:"heartbeat_seconds"`
	TimeoutSeconds   int    `json:"timeout_seconds"`
	Error            string `json:"error"`
	TaskData         string `json:"task_data"`
}

// String serialize task data
func (t *TaskInnerData) String() (string, error) {
	return toolkit.EncodeToString(t)
}

// ExceptionData State异常数据结构
type TaskExceptionData struct {
	Cause      string      `json:"cause"`       // Detail of the failure
	Error      string      `json:"error"`       // Error Code of the failure
	ErrorMatch []string    `json:"error_match"` // Error Code match list
	Output     interface{} `json:"output"`      // Output of state
	Extra      string      `json:"extra"`       // Extra information about exception
}

// LoadTaskData load task data from string
func LoadTaskData(data string) (*TaskInnerData, error) {
	var taskdata TaskInnerData
	err := json.Unmarshal([]byte(data), &taskdata)
	return &taskdata, err
}

// TaskNextState Task Next State
type TaskNextState struct {
	Name         string
	Output       interface{}
	OutputString string
	Time         time.Time
}

// TaskSendFailureResult
type TaskSendFailureResult struct {
	Error string
	Cause string
}

// NewTaskFromID NewTaskFromID
func NewTaskFromID(id int, executor *Executor) (*Task, error) {

	var err error

	dbStep, err := executor.ExecutionService.QueryStepByID(id, []string{}, nil)
	if err != nil {
		return nil, err
	}
	state, err := NewTaskFromData(dbStep, executor)
	return state, err
}

func NewTaskFromData(dbStep *po.Step, executor *Executor) (*Task, error) {

	baseStep, err := NewExecutionStep(dbStep, executor)
	if err != nil {
		return nil, err
	}
	state, ok := baseStep.State.(*states.Task)
	if !ok {
		return nil, fmt.Errorf("step state is not Task type")
	}
	task := &Task{
		ExecutionStep: baseStep,
		TaskState:     state,
	}
	return task, err

}

// GetInnerData GetInnerData
func (t *Task) GetInnerData() (*TaskInnerData, error) {

	var taskData TaskInnerData
	var err error

	if t.Data.Data == "" {
		taskData = TaskInnerData{
			Retry:          make([]int, len(t.TaskState.TaskBody.Retry)),
			Status:         TaskRunningStatus.InitData,
			HeartbeatCount: 0,
		}
	} else {
		err = toolkit.Decode([]byte(t.Data.Data), &taskData)
		if err != nil {
			return nil, err
		}
		if taskData.Retry == nil {
			taskData.Retry = make([]int, len(t.TaskState.TaskBody.Retry))
		}
		if taskData.Status == "" {
			taskData.Status = TaskRunningStatus.InitData
		}
	}
	taskData.HeartbeatSeconds = int(t.TaskState.HeartbeatSeconds)
	taskData.TimeoutSeconds = int(t.TaskState.TimeoutSeconds)

	return &taskData, nil
}

// Run Run
func (t *Task) Run(message queue.InnerMessageBody) error {
	var err error
	var updateTask po.Step

	var dbStep = t.Data
	starttime := time.Now()

	var executeMsg = DefaultStepExecuteMessage
	// 兼容历史消息
	if message.Data != "" {
		err = json.Unmarshal([]byte(message.Data), &executeMsg)
		if err != nil {
			return err
		}
	}

	if executeMsg.Block && t.TaskState.Block {
		// 如果是阻塞的， 则不执行
		slog.Info(fmt.Sprintf("task step '%d' is blocked", dbStep.ID))

		tx, maker := t.Executor.ExecutionService.MetaDB.NewTxMaker(nil)
		defer maker.Close(&err)

		dbStepGroupQueryCond := po.StepGroup{
			ExecutionID: dbStep.ExecutionID,
			SubGroupID:  dbStep.GroupID,
		}
		var dbStepGroup po.StepGroup
		err = tx.Where(dbStepGroupQueryCond).Take(&dbStepGroup).Error
		if err != nil {
			return err
		}

		updateTask = po.Step{
			Status: string(StepStatus.Blocked),
		}

		err = tx.Where(po.Step{ID: dbStep.ID}).Updates(&updateTask).Error
		if err != nil {
			tx.Rollback()
			return err
		}

		tx.Commit()

		// message queue send create message
		message := NewStepMessage(dbStep.ExecutionID, MessageType.StepGroupBlocked, dbStepGroup.StepID, nil)
		err = t.Executor.SendInnerMessage(message, nil)
		if err != nil {
			return err
		}

		finishtime := time.Now()
		// var events []ExecutionEvent
		event1 := vo.ExecutionEvent{
			ExecutionID: dbStep.ExecutionID,
			StartTime:   starttime,
			FinishTime:  finishtime,
			StepID:      dbStep.ID,
			StepName:    dbStep.Name,
			Data:        EventContent_TaskBlocked{},
		}
		t.Executor.SendExecutionEvents(event1)
		return nil

	}

	// 判断 是否超过 MaxExecute, 如果超过， 则拒绝执行
	if dbStep.ExecuteCount >= int(t.State.GetBaseState().MaxExecuteTimes) {
		err = fmt.Errorf("task step execute count reach 'MaxExecuteTimes' argument")
		return err
	}

	inputdata, err := t.GetInput()
	if err != nil {
		return err
	}
	inputDatabyte, err := json.Marshal(inputdata)
	if err != nil {
		return err
	}
	inputDataStr := string(inputDatabyte)

	taskInnerData, err := t.GetInnerData()
	if err != nil {
		return err
	}

	taskInnerData.HeartbeatCount += 1

	taskdatastr, err := taskInnerData.String()
	if err != nil {
		return err
	}

	var events = []vo.ExecutionEvent{}
	resourceUri, err := states.ParseResource(t.TaskState.Resource)
	if err != nil {
		slog.Error(fmt.Sprintf("parse resource '%s' error: %s", t.TaskState.Resource, err.Error()))
		return err
	}
	if resourceUri.ResourceType == string(states.ResourceType.Activity) {
		// 如果是Activity类型，则将待执行的task 写入TaskToken, 等待worker GetActivityTask & SendTaskSuccess
		// TaskToken表是完整的所有待执行的task, ActivityTask只包含还未被抓取的Task。
		// 生成的Token 会先写入 TaskToken和ActivityTask,
		// worker GetActivityTask后， 会删除ActivityTask并给用户返回 Token .
		// worker SendTaskSuccess后，会删除 TaskToken中的记录。
		// 所以 TaskToken 中的UUID 在ActivityTask 中一定是合法的

		// 开始事务
		tx, maker := t.Executor.ExecutionService.MetaDB.NewTxMaker(nil)
		defer maker.Close(&err)

		updateTask = po.Step{
			Status:       string(StepStatus.Wait),
			Data:         taskdatastr,
			ExecuteCount: dbStep.ExecuteCount + 1, // 执行次数 + 1
			Resource:     t.TaskState.Resource,
		}

		err = tx.Where(po.Step{ID: dbStep.ID}).Updates(&updateTask).Error
		if err != nil {
			tx.Rollback()
			return err
		}
		uuidStr, err := toolkit.CreateUUID()
		if err != nil {
			slog.Error(err.Error())
			tx.Rollback()
			return err
		}
		token := po.TaskToken{
			StepID: dbStep.ID,
			Token:  uuidStr,
		}

		err = tx.Create(&token).Error
		if err != nil {
			tx.Rollback()
			return err
		}

		// 创建 activitytask data
		atd := ActivityTaskData{
			TimeoutSeconds:   taskInnerData.TimeoutSeconds,
			HeartbeatSeconds: taskInnerData.HeartbeatSeconds,
			HeartbeatCount:   taskInnerData.HeartbeatCount,
		}
		atdStr, err := toolkit.EncodeToString(atd)
		if err != nil {
			tx.Rollback()
			return err
		}
		activitytask := po.ActivityTask{
			ExecutionID: dbStep.ExecutionID,
			StepID:      dbStep.ID,
			Resource:    resourceUri.Resource,
			Input:       inputDataStr,
			Data:        atdStr,
			Token:       uuidStr,
		}
		err = tx.Create(&activitytask).Error
		if err != nil {
			tx.Rollback()
			return err
		}
		tx.Commit()

		finishtime := time.Now()
		event1 := vo.ExecutionEvent{
			ExecutionID: dbStep.ExecutionID,
			StartTime:   starttime,
			FinishTime:  finishtime,
			StepID:      dbStep.ID,
			StepName:    dbStep.Name,
			Data: EventContent_TaskInitialized{
				TimeoutSeconds:   atd.TimeoutSeconds,
				HeartBeatSeconds: atd.HeartbeatSeconds,
				Input:            inputdata,
				Resource:         resourceUri.Resource,
			},
		}
		events = append(events, event1)
		t.Executor.SendExecutionEvents(events...)

	} else {
		err = fmt.Errorf("resource type '%s' not support", resourceUri.ResourceType)
		return err
	}
	return nil
}

// ProcessEvent ProcessEvent
func (t *Task) ProcessEvent(message queue.InnerMessageBody) error {

	var err error
	// TaskHeartBeatTimeup  string // Task 心跳时间到
	// TaskTimeoutTimeup    string // Task 超时时间到
	// TaskStateSend        string // Task 发送执行结果
	switch message.Type {
	case MessageType.TaskHeartBeatTimeup:
		err = t.ProcessTaskHeartbeatTimeout(message)
	case MessageType.TaskStateWakeup:
		err = t.ProcessTaskStateWakeup(message)
	case MessageType.TaskTimeoutTimeup:
		err = t.ProcessTaskTimeout(message)
	case MessageType.TaskStateSend:
		err = t.ProcessTaskStateSendAfter()
	default:
		slog.Error(fmt.Sprintf("unknown message type: %s", message.Type))
		err = nil
	}
	if err != nil {
		return err
	}
	return err
}

// GetNextState Get Task Next State
func (t *Task) GetNextState() (*NextStep, error) {

	var nextstep NextStep
	var err error

	var exceptiondata = TaskExceptionData{}
	err = json.Unmarshal([]byte(t.Data.Exception), &exceptiondata)
	if err != nil {
		return nil, err
	}

	// input
	var input interface{}
	err = json.Unmarshal([]byte(t.Data.Input), &input)
	if err != nil {
		return nil, err
	}

	taskInnerData, err := t.GetInnerData()
	if err != nil {
		return nil, err
	}

	var taskSendData = states.TaskSendData{}
	// if worker call SendTaskSuccess
	if taskInnerData.Status == TaskRunningStatus.SendSuccess {
		taskSendData = states.TaskSendData{
			Success: true,
			Output:  exceptiondata.Output,
			Retry:   taskInnerData.Retry,
			Errors:  []string{},
		}
	} else {
		// if worker call SendTaskFailure
		erroroutput := TaskSendFailureResult{
			Error: exceptiondata.Error,
			Cause: exceptiondata.Cause,
		}

		taskSendData = states.TaskSendData{
			Success: false,
			Output:  erroroutput,
			Retry:   taskInnerData.Retry,
			Errors:  exceptiondata.ErrorMatch, //异常错误匹配多个错误类型
		}
	}
	ns, err := t.TaskState.GetNextState(input, taskSendData)
	// 执行失败
	if err != nil {
		return nil, err
	}

	nextstep = NextStep{
		NextState: *ns,
		GroupID:   t.Data.GroupID,
	}
	return &nextstep, err
}

// GetInput GetInput
func (t *Task) GetInput() (interface{}, error) {

	var inputdata interface{}
	var err error
	err = toolkit.Decode([]byte(t.Data.Input), &inputdata)
	if err != nil {
		return nil, err
	}
	newinputdata, err := t.TaskState.GetParametersInput(inputdata)
	return newinputdata, err
}

// GetTaskTimeout 获得超时时间点
// 为未来预留API结构，如果以后支持 TimeoutSecondsPath 和 HeartbeatSecondsPath 类似语法
func (t *Task) GetTaskTimeout() (states.TaskTimeout, error) {
	result, err := t.TaskState.GetTaskTimeout()
	return result, err
}

// ProcessTaskTimeout  处理Task 超时信息
func (t *Task) ProcessTaskTimeout(message queue.InnerMessageBody) error {

	var err error
	var dbStep *po.Step

	timeoutMsg := TaskTimeoutMessage{}
	err = json.Unmarshal([]byte(message.Data), &timeoutMsg)
	if err != nil {
		return err
	}

	dbStep, err = t.Executor.ExecutionService.QueryStepByID(
		message.StepID, append(StepFields.L1, "execute_count"), nil)
	if err != nil {
		return err
	}

	/*
		当前节点当前不是Running ,说明节点已经被终止了或者执行成功了, 超时事件不再有效
	*/
	if dbStep.Status != string(StepStatus.Running) {
		return nil
	}
	/*
		如果消息计数器与内部计数器不一致，说明此次计时器消息并不是当前正在执行的。消息的时效性已经失效， 忽略消息。
	*/
	if timeoutMsg.ExecuteCount != dbStep.ExecuteCount {
		return nil
	}
	// 拿到Task token
	// 发送执行失败的消息
	errorname := StandardErrorNames.StatesTimeout
	err = t.SendTaskFailure(context.Background(), errorname, "task timeout extended", nil)
	return err
}

func (t *Task) ProcessTaskHeartbeatTimeout(message queue.InnerMessageBody) error {

	var err error

	var dbStep *po.Step
	timeoutMsg := TaskTimeoutMessage{}
	err = json.Unmarshal([]byte(message.Data), &timeoutMsg)
	if err != nil {
		return err
	}

	tx, maker := t.Executor.ExecutionService.MetaDB.NewTxMaker(nil)
	defer maker.Close(&err)
	txf := rdb.ForUpdate(tx)
	dbStep, err = t.Executor.ExecutionService.QueryStepByID(
		dbStep.ID, append(StepFields.L1, "execute_count", "data"), txf)
	if err != nil {
		return err
	}
	if dbStep.Status != string(StepStatus.Running) {
		return nil
	}
	taskdata, err := LoadTaskData(dbStep.Data)
	if err != nil {
		return err
	}
	// 如果不是当前次的定时器消息，忽略
	if !(timeoutMsg.HeartbeatCount == taskdata.HeartbeatCount && timeoutMsg.ExecuteCount == dbStep.ExecuteCount) {
		return nil
	}

	// 发送执行失败的消息
	errorname := StandardErrorNames.StatesHearbeatTimeout
	err = t.SendTaskFailure(context.Background(), errorname, "", tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return err

}

// ProcessTaskStateSendAfter 处理task 执行后的事件
// nolint: funlen
func (t *Task) ProcessTaskStateSendAfter() error {
	var err error
	var dbStep *po.Step
	starttime := time.Now()

	nextstate, err := t.GetNextState()
	// Execution Failed
	if err != nil {
		return err
	}

	outputstr, err := toolkit.ToString(nextstate.Output)
	if err != nil {
		return err
	}

	// 发现是重试， 还在自己原来的节点
	if nextstate.Retry {
		now := time.Now()

		taskInnerData, err := t.GetInnerData()
		if err != nil {
			return err
		}
		innerdatastr, err := toolkit.ToString(taskInnerData)
		if err != nil {
			return err
		}

		// 开始事务
		tx, maker := t.Executor.ExecutionService.MetaDB.NewTxMaker(nil)
		defer maker.Close(&err)

		updatetask := po.Step{
			Data:   innerdatastr,
			Status: string(StepStatus.Sleep),
		}

		err = tx.Where(po.Step{ID: dbStep.ID}).Updates(&updatetask).Error
		if err != nil {
			return err
		}
		tx.Commit()
		retrytime := now.Add(nextstate.Delay)

		finishtime := time.Now()
		event1 := vo.ExecutionEvent{
			ExecutionID: dbStep.ExecutionID,
			StartTime:   starttime,
			FinishTime:  finishtime,
			StepID:      dbStep.ID,
			StepName:    dbStep.Name,
			Data: EventContent_TaskStateRetry{
				RetryIndex: nextstate.RetryIndex,
				RetryTime:  retrytime,
			},
		}
		t.Executor.SendExecutionEvents(event1)

		// 发送重试消息，带上了 当前执行次数
		mdec := StepWakeupMessage{
			ExecuteCount: dbStep.ExecuteCount,
		}

		// message queue send create message
		// send message
		message := NewStepMessage(dbStep.ExecutionID, MessageType.TaskStateWakeup, dbStep.ID, mdec)
		err = t.Executor.SendInnerMessage(message, &retrytime)

		if err != nil {
			return err
		}
		return nil

	} else {
		// 准备开始下一个节点
		// 当前节点更新为 Success
		finishtime := time.Now()
		tx, maker := t.Executor.ExecutionService.MetaDB.NewTxMaker(nil)
		defer maker.Close(&err)

		updatetask := po.Step{
			FinishTime: &finishtime,
			Status:     string(StepStatus.Success),
			Output:     outputstr,
		}

		err = tx.Where(po.Step{ID: dbStep.ID}).Updates(&updatetask).Error
		if err != nil {
			return err
		}
		tx.Commit()

		event1 := vo.ExecutionEvent{
			ExecutionID: dbStep.ExecutionID,
			StartTime:   starttime,
			FinishTime:  finishtime,
			StepID:      dbStep.ID,
			StepName:    dbStep.Name,
			Data: EventContent_TaskStateExited{
				Output: nextstate.Output,
			},
		}
		t.Executor.SendExecutionEvents(event1)

		fns := FindNextStep{
			Name:    nextstate.Name,
			GroupID: dbStep.GroupID,
		}
		// message queue send create message
		message := NewStepMessage(dbStep.ExecutionID, MessageType.FindNextState, dbStep.ID, fns)
		err = t.Executor.SendInnerMessage(message, nil)

		if err != nil {
			slog.Error(err.Error())
			return err
		}
		return nil

	}
}

// ProcessTaskStateWakeup Task step wakeup when retry sleep is finished
// nolint: funlen
func (t *Task) ProcessTaskStateWakeup(message queue.InnerMessageBody) error {

	var err error
	var dbStep *po.Step

	var stepWakeupMsg = StepWakeupMessage{}

	// 兼容
	if message.Data != "" {
		err = json.Unmarshal([]byte(message.Data), &stepWakeupMsg)
		if err != nil {
			return err
		}
	}

	tx, maker := t.Executor.ExecutionService.MetaDB.NewTxMaker(nil)
	defer maker.Close(&err)

	dbStep, err = t.Executor.ExecutionService.QueryStepByID(
		message.StepID, []string{"id", "execution_id", "status", "execute_count"}, tx)

	if err != nil {
		return err
	}
	// 状态不对, 忽略消息
	if dbStep.Status != string(StepStatus.Sleep) {
		return nil
	}
	// 次数不等，不是当前的重试的消息，忽略消息
	if stepWakeupMsg.ExecuteCount != 0 && dbStep.ExecuteCount != stepWakeupMsg.ExecuteCount {
		return nil
	}
	updatestate := po.Step{
		Status: string(StepStatus.Wakeup),
	}

	err = tx.Where(po.Step{ID: dbStep.ID}).Updates(&updatestate).Error

	if err != nil {
		return err
	}
	tx.Commit()

	// message queue send create message
	newmsg := NewStepMessage(dbStep.ExecutionID, MessageType.StateExecute, dbStep.ID, nil)
	err = t.Executor.SendInnerMessage(newmsg, nil)
	if err != nil {
		return err
	}

	return nil

}

// ProcessRuntimeError ProcessRuntimeError
func (t *Task) ProcessRuntimeError(catcherr error) error {

	return nil
}

// SendTaskFailure send failure message for task step
// nolint: funlen
func (t *Task) SendTaskFailure(ctx context.Context, errorname string, cause string, session rdb.Tx) error {

	var err error
	var dbStep *po.Step
	starttime := time.Now()

	reqinfo := vo.GetRequestInfo(ctx)

	taskInnerData, err := t.GetInnerData()
	if err != nil {
		return err
	}
	taskInnerData.Status = TaskRunningStatus.SendFailed

	taskInnerDataStr, err := toolkit.ToString(taskInnerData)
	if err != nil {
		return err
	}

	// 计算exception data
	errornames := ExtendErrorNames(errorname)
	ted := TaskExceptionData{
		Error:      errorname,
		Cause:      cause,
		ErrorMatch: errornames,
		Extra:      reqinfo.String(),
	}

	tedStr, err := toolkit.ToString(ted)
	if err != nil {
		return err
	}
	if err != nil {
		slog.Error(fmt.Sprintf("SendTaskFailure Task ID '%d'  err : '%s' ", t.Data.ID, err.Error()))
		return err
	}
	slog.Info(fmt.Sprintf("SendTaskFailure:  Task ID '%d'  error : %s  cause :  %s ", t.Data.ID, errorname, cause))

	tx, maker := t.Executor.ExecutionService.MetaDB.NewTxMaker(session)
	defer maker.Close(&err)

	// 加锁查看当前step状态
	txf := rdb.ForUpdate(tx)
	dbStep, err = t.Executor.ExecutionService.QueryStepByID(t.Data.ID, StepFields.L1, txf)
	if err != nil {
		return err
	}
	// 说明当前状态已经不是Running, 忽略消息
	if dbStep.Status != string(StepStatus.Running) {
		return fmt.Errorf("%w: current step status '%s'", vo.ErrorStepStatus, dbStep.Status)
	}

	tasktoken := po.TaskToken{
		StepID: t.Data.ID,
	}

	// 查找token匹配到的StepID
	err = tx.Model(new(po.TaskToken)).Where(tasktoken).Take(&tasktoken).Error
	// 如果查找错误，或者没找到
	if err != nil {
		return err
	}

	updatetask := po.Step{
		Data:      taskInnerDataStr,
		Exception: tedStr,
		// 更新状态到TaskSubmitted , 防止其他流程冲突，
		// 不能更新，这时候还是Running，TaskStatus 没有Submitted 状态，只能记录在Inner 状态中
		// Status: string(TaskStatus.TaskSubmitted),
	}

	err = tx.Where(po.Step{ID: tasktoken.StepID}).Updates(&updatetask).Error
	if err != nil {
		return err
	}

	err = tx.Where(po.TaskToken{ID: tasktoken.ID}).Delete(new(po.TaskToken)).Error
	if err != nil {
		return err
	}
	tx.Commit()

	finishtime := time.Now()

	event := vo.ExecutionEvent{
		ExecutionID: dbStep.ExecutionID,
		StepID:      dbStep.ID,
		StepName:    dbStep.Name,
		StartTime:   starttime,
		FinishTime:  finishtime,
		Data: EventContent_TaskSubmitFailed{
			Error: ted.Error,
			Cause: ted.Cause,
		},
	}
	t.Executor.SendExecutionEvents(event)

	// message queue send create message
	message := NewStepMessage(dbStep.ExecutionID, MessageType.TaskStateSend, dbStep.ID, nil)
	err = t.Executor.SendInnerMessage(message, nil)
	if err != nil {
		return err
	}
	return nil
}

// SendTaskSuccess send task execute result to step
/* @token : task token
 * @output: worker execute output
 */
// nolint: funlen
func (t *Task) SendTaskSuccess(ctx context.Context, output interface{}) error {

	var err error
	var dbStep *po.Step
	reqinfo := vo.GetRequestInfo(ctx)
	starttime := time.Now()

	taskInnerData, err := t.GetInnerData()
	if err != nil {
		return err
	}

	taskInnerData.Status = TaskRunningStatus.SendSuccess
	taskInnerStr, err := toolkit.ToString(taskInnerData)
	if err != nil {
		return err
	}
	// 计算exception data
	ted := TaskExceptionData{
		Output: output,
		Extra:  reqinfo.String(),
	}
	tedStr, err := toolkit.ToString(ted)
	if err != nil {
		return err
	}
	// 开始事务
	tx, maker := t.Executor.ExecutionService.MetaDB.NewTxMaker(nil)
	defer maker.Close(&err)

	// 加锁
	txf := rdb.ForUpdate(tx)
	dbStep, err = t.Executor.ExecutionService.QueryStepByID(t.Data.ID, StepFields.L1, txf)
	if err != nil {
		return err
	}

	// 状态判断
	if dbStep.Status != string(StepStatus.Running) {
		return fmt.Errorf("%w: current step status '%s'", vo.ErrorStepStatus, dbStep.Status)
	}

	dbTaskToken := po.TaskToken{}

	// 查找token匹配到的StepID
	err = tx.Where(po.TaskToken{
		StepID: dbStep.ID,
	}).Take(&dbTaskToken).Error
	// 如果查找错误，或者没找到
	if err != nil {
		return err
	}

	finishtime := time.Now()
	updatetask := po.Step{
		FinishTime: &finishtime,
		Data:       taskInnerStr,
		Exception:  tedStr,
	}

	err = tx.Where(po.Step{ID: dbStep.ID}).Updates(&updatetask).Error

	if err != nil {
		slog.Error(fmt.Sprintf("Failed To Update Task, StepID: [ %d ], Error: %s", dbStep.ID, err.Error()))
		return err
	}
	err = tx.Where(po.TaskToken{ID: dbTaskToken.ID}).Delete(new(po.TaskToken)).Error

	if err != nil {
		slog.Error(fmt.Sprintf("Failed To Delete Token, StepID: [ %d ], Error: %s", dbTaskToken.StepID, err.Error()))
		return err
	}
	tx.Commit()

	event := vo.ExecutionEvent{
		ExecutionID: dbStep.ExecutionID,
		StepID:      dbStep.ID,
		StepName:    dbStep.Name,
		StartTime:   starttime,
		FinishTime:  time.Now(),
		Data: EventContent_TaskSubmitted{
			Result: output,
		},
	}
	t.Executor.SendExecutionEvents(event)

	// message queue send create message
	message := NewStepMessage(dbStep.ExecutionID, MessageType.TaskStateSend, dbStep.ID, nil)
	err = t.Executor.SendInnerMessage(message, nil)

	if err != nil {
		slog.Error(
			fmt.Sprintf("Send Inner Message Failed: ExecutionID %d, StepID %d, Priority %s, Type %s ErrorDetail: %s",
				message.ExecutionID, message.StepID, message.Class, message.Type, err.Error()))
		return err
	}
	return nil
}

// SendTaskRefer send task refer information
func (t *Task) SendTaskReference(ctx context.Context, title string, url string) error {

	var err error
	var dbStep *po.Step
	step_id := t.Data.ID
	starttime := time.Now()
	reqinfo := vo.GetRequestInfo(ctx)

	refer := vo.Reference{

		Title: title,
		URL:   url,
	}
	referstr, err := toolkit.ToString(refer)
	if err != nil {
		return err
	}

	// 开始事务
	tx, maker := t.Executor.ExecutionService.MetaDB.NewTxMaker(nil)
	defer maker.Close(&err)

	dbStep, err = t.Executor.ExecutionService.QueryStepByID(step_id, StepFields.L1, tx)
	if err != nil {
		return err
	}

	updateUserData := po.UserStepData{
		ID:        step_id,
		Reference: referstr,
	}
	// use insert or update
	err = tx.Model(new(po.UserStepData)).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "step_id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"reference": referstr}),
	}).Create(
		&updateUserData,
	).Error

	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	finishtime := time.Now()
	event := vo.ExecutionEvent{
		ExecutionID: dbStep.ExecutionID,
		StepID:      dbStep.ID,
		StepName:    dbStep.Name,
		StartTime:   starttime,
		FinishTime:  finishtime,
		Data: EventContent_TaskSendReference{
			Reference:   referstr,
			RequestInfo: reqinfo,
		},
	}
	t.Executor.SendExecutionEvents(event)

	return nil
}

// SendTaskHeartbeat send heartbeat of a task step
func (t *Task) SendTaskHeartbeat(ctx context.Context, message string, session rdb.Tx) error {

	var err error
	step_id := t.Data.ID
	starttime := time.Now()
	requestinfo := vo.GetRequestInfo(ctx)

	// Todo : 重置heartbeattimeoutsecond 的时间
	// 开始事务
	tx, maker := t.Executor.ExecutionService.MetaDB.NewTxMaker(session)
	defer maker.Close(&err)

	txf := rdb.ForUpdate(tx)
	dbStep, err := t.Executor.ExecutionService.QueryStepByID(step_id, append(StepFields.L1, "execute_count", "data"), txf)
	if err != nil {
		return err
	}
	innerdata, err := LoadTaskData(dbStep.Data)
	if err != nil {
		return err
	}
	// HeartbeatCount +1
	innerdata.HeartbeatCount++
	innerdatastr, err := innerdata.String()
	if err != nil {
		return err
	}
	err = tx.Where(po.Step{ID: step_id}).Updates(po.Step{Data: innerdatastr}).Error
	if err != nil {
		return err
	}
	tx.Commit()

	// 发送超时消息
	mec := TaskTimeoutMessage{
		HeartbeatCount: innerdata.HeartbeatCount,
		ExecuteCount:   dbStep.ExecuteCount,
	}

	// 如果设置有有效的心跳时间的话，发送心跳消息,发送新的心跳计数器
	if innerdata.HeartbeatSeconds > 0 {
		// send timing message 发送定时消息,
		// 发送超时事件， 准备做超时处理。
		nextTimeoutTime := time.Now().Add(time.Second * time.Duration(innerdata.HeartbeatSeconds))
		tasktimeoutMsg := NewStepMessage(dbStep.ExecutionID,
			MessageType.TaskHeartBeatTimeup, dbStep.ID, mec)
		err = t.Executor.SendInnerMessage(tasktimeoutMsg, &nextTimeoutTime)
		if err != nil {
			return err
		}
	}

	finishtime := time.Now()
	event := vo.ExecutionEvent{
		ExecutionID: dbStep.ExecutionID,
		StepID:      dbStep.ID,
		StepName:    dbStep.Name,
		StartTime:   starttime,
		FinishTime:  finishtime,
		Data: EventContent_TaskSendHeartbeat{
			Message:     message,
			RequestInfo: requestinfo,
		},
	}
	t.Executor.SendExecutionEvents(event)

	return nil
}
