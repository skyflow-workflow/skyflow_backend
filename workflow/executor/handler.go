// file: workflow/executor/handler.go
// for execution api realize
package executor

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"time"

	"github.com/mmtbak/microlibrary/rdb"
	"github.com/skyflow-workflow/skyflow_backbend/pkg/toolkit"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/parser"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/parser/states"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/po"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/repository/queue"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/vo"
	"gorm.io/gorm"
)

// StepErrorProcess state 执行报错时候的处理方式
func (svc *executionService) StepErrorProcess(catcherr error, dbStep *po.Step, msg queue.InnerMessage) error {
	var err error
	var tx rdb.Tx
	var maker *rdb.TxMaker

	var updatetask po.Step
	starttime := time.Now()

	// 如果是状态异常，忽略消息
	if errors.Is(catcherr, vo.ErrorStepStatus) {
		slog.Error(fmt.Sprintln(err, msg))
		return nil
	}

	// 为什么要在这里更新步骤状态？
	// 因为有些步骤失败的时候， 并没有进行到 更新自身状态的阶段， 比如 初始化节点，就报错了，最终会catch 到这里处理异常。

	//
	// 如果是最外层的步骤
	// 更新节点自身状态为Failed，并调用ExecutionFailed

	if dbStep.GroupID == states.StartGroupID {
		tx, maker = svc.MetaDB.NewTxMaker(nil)
		defer maker.Close(&err)

		updatetask = po.Step{
			Status: string(StepStatus.Failed),
		}

		err = tx.Where(po.Step{ID: dbStep.ID}).Updates(&updatetask).Error

		if err != nil {
			err2 := fmt.Errorf("update state '%d ' Status [ Failed ] error : %w", dbStep.ID, err)
			slog.Error(err2.Error())
			return err2
		}
		tx.Commit()

		err = svc.ExecutionErrorProcess(catcherr, dbStep.ExecutionID)
		return err
	}

	// 如果是内部子流程，需要处理根据状态。处理StepGroup/ Parallel/Map 等步骤判断，最后发送失败。

	var message queue.InnerMessageBody
	var events []vo.ExecutionEvent
	var stepGroupId int

	// 更新step 状态
	tx, maker = svc.MetaDB.NewTxMaker(nil)
	defer maker.Close(&err)

	// 如果是stepgroup 步骤执行失败
	if dbStep.Type == string(states.StateTypes.StateGroup) {
		stepGroupId = dbStep.GroupID
	} else {
		// 更新
		updatetask = po.Step{
			Status: string(StepStatus.Failed),
		}

		err = tx.Where(po.Step{ID: dbStep.ID}).Updates(&updatetask).Error

		dbStepGroup, err := svc.QueryStepGroupBySubGroupID(dbStep.ExecutionID, dbStep.GroupID, tx)
		if err != nil {
			return err
		}
		stepGroupId = dbStepGroup.StepID
	}
	tx.Commit()
	// 准备发送事件
	var event1 vo.ExecutionEvent
	finishtime := time.Now()

	event1 = vo.ExecutionEvent{
		ExecutionID: dbStep.ExecutionID,
		StartTime:   starttime,
		FinishTime:  finishtime,
		StepID:      dbStep.ID,
		StepName:    dbStep.Name,
		Data: EventContent_StateFailed{
			Error: StandardErrorNames.StatesRuntime,
			Cause: catcherr.Error(),
		},
	}
	events = append(events, event1)

	svc.SendExecutionEvents(events...)
	// message queue send create message
	message = NewStepMessage(dbStep.ExecutionID, MessageType.StepGroupFailed, stepGroupId, nil)
	err = svc.SendInnerMessage(message, nil)
	if err != nil {
		return err
	}
	return nil
}

// ExecutionErrorProcess  ExecutionErrorProcess
func (svc *executionService) ExecutionErrorProcess(catcherr error, execution_id int) error {

	var err error
	starttime := time.Now()

	exe, err := svc.NewExecutionFromID(execution_id, ExecutionFields.L1, nil)
	if err != nil {
		return err
	}

	err = exe.ChangeExecutionStatus(ExecutionStatus.Failed, nil)
	if err != nil {
		return err
	}
	finishtime := time.Now()

	event1 := vo.ExecutionEvent{
		ExecutionID: execution_id,
		StartTime:   starttime,
		FinishTime:  finishtime,
		Data:        EventContent_ExecutionFailed{Error: "", Cause: catcherr.Error()},
	}
	svc.SendExecutionEvents(event1)
	return nil
}

// StartExecution 创建一个Execution
func (svc *executionService) StartExecution(req vo.StartExecutionRequest) (po.Execution, error) {
	var err error
	var dbNull po.Execution

	sm, err := parser.ParseStateMachine(req.StateMachineDefinition)
	if err != nil {
		return dbNull, err
	}

	// 准备插入Execution
	tx, maker := svc.MetaDB.NewTxMaker(nil)
	defer maker.Close(&err)

	tx.Begin()
	resp, err := svc._StartExecution(req, sm, tx)
	if err != nil {
		return dbNull, err
	}
	tx.Commit()

	svc.SendExecutionEvents(resp.Events...)
	for _, msg := range resp.Messages {
		err = svc.SendInnerMessage(msg, nil)
		if err != nil {
			return dbNull, err
		}
	}
	return *resp.Data, nil
}

// _StartExecution 创建一个Execution
func (svc *executionService) _StartExecution(req vo.StartExecutionRequest, statemachine *states.StateMachine, tx rdb.Tx,
) (resp StartExecutionResponse, err error) {

	var uuids string

	_, err = states.ToMap(req.Input)
	if err != nil {
		err = fmt.Errorf("%w: input error: %s", vo.ErrorParameterInvalid, err.Error())
		return
	}

	header := statemachine.StateMachineHeader
	headerstr, err := states.ToString(header)
	if err != nil {
		return
	}

	starttime := time.Now()
	uuids = req.ExecutionUUID
	// 处理UUID， 确认UUID合法
	if uuids != "" {
		checkexecution := po.Execution{}
		err = tx.Select(ExecutionFields.L1).Where("uuid = ?", req.ExecutionUUID).Take(&checkexecution).Error
		// 能查到
		if err == nil {
			err = vo.ErrorExecutionUUIDExisted
			return
		}
		// 如果是其他错误，返回
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return
		}
	} else {
		// 自动创建UUID
		checkpass := false
		// 最大尝试次数
		var maxtryuuid = 3
		for i := 0; i < maxtryuuid; i++ {
			uuids, err = toolkit.CreateUUID()
			if err != nil {
				continue
			}
			_, err = svc.QueryExecutionByUUID(uuids, ExecutionFields.L1, tx)
			// 存在
			if err == nil {
				continue
			}
			// 如果是其他错误，返回
			if !errors.Is(err, vo.ErrorExecutionNotFound) {
				return
			}
			// 发现可以用
			checkpass = true
			break
		}
		if !checkpass {
			err = ErrorCreateExecutionUUIDFailed
			return
		}
	}
	exedatastr, err := states.ToString(req.Data)
	if err != nil {
		return
	}

	dbexecution := po.Execution{
		URI:             req.StateMachineURI,
		UUID:            uuids,
		MaxExecuteIndex: 0,
		// trim一下
		Definition:   strings.TrimSpace(req.StateMachineDefinition),
		Title:        req.Title,
		Input:        req.Input,
		Output:       "",
		ExecuteCount: 0,
		Data:         exedatastr,
		Status:       string(ExecutionStatus.Created),
		Header:       headerstr,
	}
	// 创建时，不更新时间相关字段，不然会报错
	err = tx.Create(&dbexecution).Error
	if err != nil {
		return
	}
	//  insert shade execution info
	dbexecutionshade := po.ExecutionShade{
		ID:   dbexecution.ID,
		UUID: dbexecution.UUID,
	}
	err = tx.Create(&dbexecutionshade).Error
	if err != nil {
		return
	}

	finishtime := time.Now()
	event := vo.ExecutionEvent{
		ExecutionID: dbexecution.ID,
		StartTime:   starttime,
		FinishTime:  finishtime,
		Data:        EventContent_ExecutionCreated{Input: req.Input},
	}

	// message queue send create message
	message := NewExecutionMessage(dbexecution.ID, MessageType.ExecutionInit, nil)

	resp.Data = dbexecution
	resp.Events = append(resp.Events, event)
	resp.Messages = append(resp.Messages, message)

	return
}

// StopExecution stop certain execution
func (svc *executionService) StopExecution(req vo.StopExecutionRequest) error {

	var err error

	starttime := time.Now()
	// 先尝试初始化 execution
	exestate, err := svc.NewExecutionFromUUID(req.ExecutionUUID, ExecutionFields.L1, nil)
	if err != nil {
		return err
	}
	// 更新Execution 状态
	err = exestate.StopExecution(req.Error, req.Cause)
	if err != nil {
		return err
	}

	// 记录事件
	finishtime := time.Now()
	event := vo.ExecutionEvent{
		ExecutionID: exestate.Data.ID,
		StartTime:   starttime,
		FinishTime:  finishtime,
		Data: EventContent_ExecutionAbort{
			Error: req.Error,
			Cause: req.Cause,
		},
	}

	svc.SendExecutionEvents(event)
	return nil

}

// RestartExecution 重新创建一个Execution
func (svc *executionService) RestartExecution(req vo.RestartExecutionRequest) (*po.Execution, error) {

	var err error
	var dbNull *po.Execution
	var dbExecution *po.Execution

	dbExecution, err = svc.QueryExecutionByUUID(req.ExecutionUUID, append(ExecutionFields.L3, "title", "input"), nil)
	if err != nil {
		return dbNull, err
	}

	wf, err := parser.ParseStateMachine(dbExecution.Definition)
	if err != nil {
		return dbNull, err
	}

	// 加锁
	lock := svc.LockService.LockExecution(dbExecution.ID)
	err = lock.Lock()
	if err != nil {
		return dbNull, err
	}
	defer lock.Unlock()

	// 开启事务
	tx, maker := svc.MetaDB.NewTxMaker(nil)
	defer maker.Close(&err)

	// 查询Execution
	dbExecution, err = svc.QueryExecutionByUUID(req.UUID, ExecutionFields.L1, tx)
	if err != nil {
		return dbNull, err
	}

	//避免重复stop
	if req.Abort == true && dbExecution.Status == string(ExecutionStatus.Aborted) {
		err = fmt.Errorf("execution has been stopped ")
		return dbNull, err
	}

	// 创建Execution Object
	exe, err := NewExecutionFromData(&dbExecution, svc)
	if err != nil {
		return dbNull, err
	}

	var respcreate StartExecutionResponse
	var respstop StopExecutionResponse

	if req.Abort == true {
		// 使用内部_StopExecution 方法执行
		respstop, err = exe._StopExecution(req.Error, req.Cause, tx)
		if err != nil {
			return dbNull, err
		}
	}

	// 准备创建新任务
	reqcreate := vo.StartExecutionRequest{
		//默认UUID 是 “”， 自动产生UUID
		StateMachineURI:        "",
		ExecutionUUID:          dbExecution.URI,
		Title:                  dbExecution.Title,
		StateMachineDefinition: dbExecution.Definition,
		Input:                  dbExecution.Input,
	}
	respcreate, err = svc._StartExecution(reqcreate, wf, tx)
	if err != nil {
		return dbNull, err
	}

	tx.Commit()

	svc.SendExecutionEvents(respcreate.Events...)
	svc.SendExecutionEvents(respstop.Events...)

	for _, msg := range respcreate.Messages {
		err = svc.SendInnerMessage(msg, nil)
		if err != nil {
			return dbNull, err
		}
	}
	// 清理过期的消息
	err = svc.InnerQueue.CleanExecutionMessage(dbExecution.ID)
	if err != nil {
		return dbNull, err
	}

	return respcreate.Data, nil

}

// ProcessFindNextStep  find next state
// NOCC:golint/fnsize("设计如此")
func (svc *executionService) ProcessFindNextStep(message queue.InnerMessageBody) error {

	var dbStep po.Step
	var err error
	step_id := message.StepID

	dbStep, err = svc.QueryStepByID(step_id, []string{"id", "name", "group_id", "output", "execution_id"}, nil)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	var fns = FindNextStep{}
	err = json.Unmarshal([]byte(message.Data), &fns)
	if err != nil {
		return err
	}

	if fns.Name == "" {
		// 如果是 End 节点 ，发送结束消息
		// message q
		// 如果是 最外层的Group, 说明任务准备结束, 保存output为Execution 的output
		if dbStep.GroupID == states.StartGroupID {

			exe, err := svc.NewExecutionFromID(dbStep.ExecutionID, []string{"id", "status"}, nil)
			if err != nil {
				return err
			}
			err = exe.ProcessSucceed(dbStep.Output)
			return err

		} else {

			// 如果是内部子流程
			// 更新 StepGroup Group output的值
			// 开始事务
			tx, maker := svc.MetaDB.NewTxMaker(nil)
			defer maker.Close(&err)

			// 查询stepgroup
			var dbStepGroup = po.StepGroup{
				ExecutionID: dbStep.ExecutionID,
				SubGroupID:  dbStep.GroupID,
			}

			err = tx.Where(dbStepGroup).Take(&dbStepGroup).Error
			if err != nil {
				return err
			}
			// 更新stepgroup lastat
			updateStepGroup := po.StepGroup{
				LastAt: dbStep.Name,
			}
			err = tx.Where(po.StepGroup{ID: dbStepGroup.ID}).Updates(&updateStepGroup).Error
			if err != nil {
				return err
			}
			// 更新 Group所关联的step 的output
			updatestep := po.Step{
				Output: dbStep.Output,
			}
			err = tx.Where(po.Step{ID: dbStepGroup.StepID}).Updates(&updatestep).Error
			if err != nil {
				return err
			}
			tx.Commit()

			newmessage := NewStepMessage(dbStep.ExecutionID, MessageType.StepGroupSucceed, dbStepGroup.StepID, nil)
			err = svc.SendInnerMessage(newmessage, time.Now())
			return err
		}
	}
	// 非End节点
	var dbnextstate = po.Step{
		Name:        fns.Name,
		GroupID:     fns.GroupID,
		ExecutionID: dbStep.ExecutionID,
	}

	// 查找next 是否存在

	// 开始事务
	tx, maker := svc.MetaDB.NewTxMaker(nil)
	defer maker.Close(&err)

	err = tx.Where(dbnextstate).Take(&dbnextstate).Error
	if err != nil {
		return err
	}

	nextupdatestate := po.Step{
		Status: string(StepStatus.WaitInit),
		Input:  dbStep.Output,
	}

	err = tx.Where(po.Step{ID: dbnextstate.ID}).Updates(&nextupdatestate).Error
	if err != nil {
		return err
	}
	tx.Commit()
	// Next State Start Message
	nextmessage := NewStepMessage(dbStep.ExecutionID, MessageType.StateNewTurn, dbnextstate.ID, nil)
	err = svc.SendInnerMessage(nextmessage, time.Now())
	if err != nil {
		return err
	}
	return nil
}

// ProcessReportStepSuspend  report step failed
// 报告有个步骤失败， 处理这个异常
// NOCC:golint/fnsize("设计如此")
func (svc *executionService) ProcessReportStepSuspend(message queue.InnerMessage) error {

	var dbstep po.Step
	var err error
	step_id := message.StepID

	dbstep, err = svc.QueryStepByID(step_id, StepFields.L2, nil)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	// 如果是 End 节点 ，发送结束消息
	// message q
	// 如果是 最外层的Group, 说明任务准备结束, 保存output为Execution 的output
	if dbstep.GroupID == states.StartGroupID {

		exe, err := svc.NewExecutionFromID(dbstep.ExecutionID, ExecutionFields.L1, nil)
		if err != nil {
			return err
		}
		err = exe.ProcessExecutionSuspend(message)
		return err
	}

	// 如果是内部子流程
	// 更新 StepGroup Group output的值
	// 开始事务

	tx := svc.MetaDB.Client.DB()
	// 查询stepgroup
	var dbstepgroup = po.StepGroup{
		ExecutionID: dbstep.ExecutionID,
		SubGroupID:  dbstep.GroupID,
	}

	err = tx.Where(dbstepgroup).Take(&dbstepgroup).Error
	if err != nil {
		return err
	}

	newmessage := NewStepMessage(dbstep.ExecutionID, MessageType.StepGroupSuspend, dbstepgroup.StepID, nil)
	err = svc.SendInnerMessage(newmessage, time.Now())
	if err != nil {
		return err
	}
	return err
}

// ProcessReportStepBlocked  report step blocked
// 报告有个步骤失败， 处理这个异常
// NOCC:golint/fnsize("设计如此")
func (svc *executionService) ProcessReportStepBlocked(message queue.InnerMessage) error {

	var dbstep po.Step
	var err error
	step_id := message.StepID

	dbstep, err = svc.QueryStepByID(step_id, StepFields.L2, nil)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	// 如果是 End 节点 ，发送结束消息
	// message q
	// 如果是 最外层的Group, 说明任务准备结束, 保存output为Execution 的output
	if dbstep.GroupID == states.StartGroupID {

		exe, err := svc.NewExecutionFromID(dbstep.ExecutionID, ExecutionFields.L1, nil)
		if err != nil {
			return err
		}
		err = exe.ProcessExecutionBlocked(message)
		return err

	}

	// 如果是内部子流程
	// 更新 StepGroup Group output的值
	// 开始事务

	tx, maker := svc.MetaDB.NewTxMaker(nil)
	defer maker.Close(&err)

	// 查询stepgroup
	var dbstepgroup = po.StepGroup{
		ExecutionID: dbstep.ExecutionID,
		SubGroupID:  dbstep.GroupID,
	}

	err = tx.Where(dbstepgroup).Take(&dbstepgroup).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	newmessage := NewStepMessage(dbstep.ExecutionID, MessageType.StepGroupBlocked, dbstepgroup.StepID, nil)
	err = svc.SendInnerMessage(newmessage, time.Now())
	if err != nil {
		return err
	}
	return err
}

// SendStepSkip 用户操作发送跳过当前处于失败状态的步骤
// @step_id 要跳过的步骤
// @next_step_name 要跳到的下一个步骤名
// @next_input 下一个步骤的输入
// NOCC:golint/fnsize("设计如此")
func (svc *executionService) SendStepSkip(ctx context.Context, req vo.SendStepSkipRequest) error {

	// sendstepskip 是跳过当前处于失败状态的步骤
	var err error
	var dbexecution po.Execution
	var dbstep = po.Step{}

	requestinfo := vo.GetRequestInfo(ctx)

	dbstep, err = svc.QueryStepByID(req.StepID, StepFields.L1, nil)
	if err != nil {
		return err
	}

	// 加锁
	lock := svc.LockService.LockExecution(dbstep.ExecutionID)
	err = lock.Lock()
	if err != nil {
		return err
	}
	defer lock.Unlock()

	// 开始事务
	tx, maker := svc.MetaDB.NewTxMaker(nil)
	defer maker.Close(&err)

	starttime := time.Now()

	// L3 包含 definition
	dbstep, err = svc.QueryStepByID(req.StepID, StepFields.L3, tx)
	if err != nil {
		return err
	}
	dbexecution, err = svc.QueryExecutionByID(dbstep.ExecutionID, ExecutionFields.L1, tx)
	if err != nil {
		return err
	}

	// Step是 Failed 状态的才能跳过
	if dbstep.Status != string(StepStatus.Failed) {
		err = fmt.Errorf("%w: current step status '%s'", vo.ErrorStepStatus, dbstep.Status)
		return err
	}

	if !states.IsExecutableStateType(dbstep.Type) {
		err = fmt.Errorf("%w: current step type '%s' is not executable", vo.ErrorParamterInvalid, dbstep.Type)
		return err
	}

	// Execution 是Failed /Running 状态
	if !slices.Contains([]string{string(ExecutionStatus.Failed), string(ExecutionStatus.Running)}, dbexecution.Status) {
		err = fmt.Errorf("%w: current execution status '%s'", vo.ErrorExecutionStatus, dbexecution.Status)
		return err
	}

	sis, err := parser.ParserStep(dbstep.Definition, dbexecution.FlowType)
	if err != nil {
		return err
	}
	bone := sis.GetNode().GetBone()

	// 判断是否是pipeline
	if dbexecution.FlowType == flow.WorkflowType.Pipeline && req.NextStepName == "" && len(bone.Next) != 0 {
		req.NextStepName = bone.Next[0]
	}

	// 判断是否是最后一个节点
	//
	// 最后一个节点,next 必需为空

	if bone.End && req.NextStepName != "" {

		err = fmt.Errorf("%w: step [ %s ] is End step, has no next step '%s'", vo.ErrorParamterInvalid, dbstep.Name, req.NextStepName)
		return err
	} else if !bone.End && !slices.Contains(bone.Next, req.NextStepName) {
		//当前节点不是最后一个节点， next 必需在 bone 的next列表中
		err = fmt.Errorf(" %w : step [ %s ] has no next step  [ %s ]  ", vo.ErrorParamterInvalid, dbstep.Name, req.NextStepName)
		return err
	}

	// 先修改公共属性
	finishtime := time.Now()
	updatestate := po.Step{
		Status:     string(StepStatus.Skip),
		Output:     req.Output,
		FinishTime: &finishtime,
	}

	// 更新step 表
	err = tx.Where(po.Step{ID: dbstep.ID}).Updates(&updatestate).Error
	if err != nil {
		return err
	}
	// 清理过期的token
	err = svc.CleanStepToken(dbstep.ID, tx)
	if err != nil {
		return err
	}

	// 递归向上修改group、parallel/map 状态
	err = svc.ChangeMasterGroupStatus(dbexecution.ID, dbstep.GroupID, ExecutionStatus.Running, tx)
	if err != nil {
		return err
	}

	// execution status Running
	updateexecution := po.Execution{
		Status: string(ExecutionStatus.Running),
	}
	err = tx.Where(po.Execution{ID: dbstep.ExecutionID}).Updates(&updateexecution).Error
	// 更新execution 表
	if err != nil {
		return err
	}

	//更新DB完成
	tx.Commit()

	event1 := vo.ExecutionEvent{
		ExecutionID: dbstep.ExecutionID,
		StepID:      dbstep.ID,
		StepName:    dbstep.Name,
		StartTime:   starttime,
		FinishTime:  finishtime,

		Data: EventContent_StepSkip{
			NextStepName: req.NextStepName,
			Output:       req.Output,
			RequestInfo:  requestinfo,
		},
	}

	event2 := vo.ExecutionEvent{
		ExecutionID: dbstep.ExecutionID,
		StartTime:   starttime,
		FinishTime:  finishtime,
		StepID:      dbstep.ID,
		StepName:    dbstep.Name,
		Data: EventContent_ExecutionInfoModified{
			Field:  DataModifyField.Status,
			Before: dbexecution.Status,
			After:  string(ExecutionStatus.Running),
		},
	}

	svc.SendExecutionEvents(event1, event2)

	fns := FindNextStep{
		Name:    req.NextStepName,
		GroupID: dbstep.GroupID,
	}

	// 写入message queue ,处理人工干预消息
	message := NewStepMessage(dbstep.ExecutionID, MessageType.FindNextStep, req.StepID, fns)
	err = svc.SendInnerMessage(message, time.Now())
	if err != nil {
		return err
	}
	return nil
}

// SendStepFailed 手动强制步骤失败
// step_id int 步骤it
func (svc *executionService) SendStepFailed(step_id int) error {

	var err error
	var dbstep *po.Step
	var dbExecution *po.Execution

	// 初始化查询信息
	dbstep, err = svc.QueryStepByID(step_id, StepFields.L1, nil)
	if err != nil {
		return err
	}

	// 加锁
	lock := svc.LockService.LockExecution(dbstep.ExecutionID)
	err = lock.Lock()
	if err != nil {
		return err
	}
	defer lock.Unlock()

	// 开始事务
	tx, maker := svc.MetaDB.NewTxMaker(nil)
	defer maker.Close(&err)

	starttime := time.Now()

	// L3 包含 definition
	dbstep, err = svc.QueryStepByID(step_id, StepFields.L1, tx)
	if err != nil {
		return err
	}
	dbExecution, err = svc.QueryExecutionByID(dbstep.ExecutionID, ExecutionFields.L1, tx)
	if err != nil {
		return err
	}

	// Step是 Failed 说明已经是Failed了。 不能再失败
	if dbstep.Status == string(StepStatus.Failed) {
		err = fmt.Errorf("%w: current step status has been '%s'", vo.ErrorStepStatus, dbstep.Status)
		return err
	}

	// Execution 是Failed /Running 状态
	if !slices.Contains([]string{string(ExecutionStatus.Running)}, dbExecution.Status) {
		err = fmt.Errorf("%w: current execution status '%s' is not running", vo.ErrorExecutionStatus, dbExecution.Status)
		return err
	}

	// 清理过期的token
	err = svc.CleanStepToken(dbstep.ID, tx)
	if err != nil {
		return err
	}
	// 修改步骤状态
	err = svc.ChangeStepStatus(dbstep.ID, StepStatus.Failed, tx)
	if err != nil {
		return err
	}
	// 修改父节点组的状态
	err = svc.ChangeStepGroupStatus(dbstep.ID, ExecutionStatus.Failed, tx)
	if err != nil {
		return err
	}
	tx.Commit()

	finishtime := time.Now()
	event := vo.ExecutionEvent{
		ExecutionID: dbstep.ExecutionID,
		StepID:      dbstep.ID,
		StepName:    dbstep.Name,
		StartTime:   starttime,
		FinishTime:  finishtime,
		Data:        EventContent_SendStepFailed{},
	}

	svc.SendExecutionEvents(event)
	return nil
}

// RedoStep 重做一个步骤
func (svc *executionService) RedoStep(step_id int) error {

	var err error
	var dbStep *po.Step
	var dbExecution *po.Execution

	// 初始化查询信息
	dbStep, err = svc.QueryStepByID(step_id, StepFields.L1, nil)
	if err != nil {
		return err
	}

	// 加锁
	lock := svc.LockService.LockExecution(dbStep.ExecutionID)
	err = lock.Lock()
	if err != nil {
		return err
	}
	defer lock.Unlock()

	// 开始事务
	tx, maker := svc.MetaDB.NewTxMaker(nil)
	defer maker.Close(&err)

	starttime := time.Now()

	// L3 包含 definition
	dbStep, err = svc.QueryStepByID(step_id, StepFields.L2, tx)
	if err != nil {
		return err
	}
	dbExecution, err = svc.QueryExecutionByID(dbStep.ExecutionID, ExecutionFields.L1, tx)
	if err != nil {
		return err
	}

	// 去掉条件
	// // Step是 Failed 说明已经是Failed了。 不能再失败
	// if dbstep.Status == string(StepStatus.Failed) {
	// 	err = fmt.Errorf("%w: current step status has been '%s'", vo.ErrorStepStatus, dbstep.Status)
	// 	return err
	// }

	// Execution 是Failed /Running 状态
	if !slices.Contains([]string{string(ExecutionStatus.Running), string(ExecutionStatus.Failed)}, dbExecution.Status) {
		err = fmt.Errorf("%w: current execution status '%s' is invalid", vo.ErrorExecutionStatus, dbExecution.Status)
		return err
	}

	// 清理过期的token
	err = svc.CleanStepToken(dbStep.ID, tx)
	if err != nil {
		return err
	}
	// 修改步骤状态
	err = svc.ChangeStepStatus(dbStep.ID, StepStatus.WaitInit, tx)
	if err != nil {
		return err
	}

	// 	递归修改组状态
	err = svc.ChangeMasterGroupStatus(dbStep.ExecutionID, dbStep.GroupID, ExecutionStatus.Running, tx)
	if err != nil {
		return err
	}

	// 更新execution 状态
	err = svc.ChangeExecutionStatus(dbStep.ExecutionID, ExecutionStatus.Running, tx)
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
		Data: EventContent_RedoStep{
			ExecuteCount: dbStep.ExecuteCount,
		},
	}

	svc.SendExecutionEvents(event)

	//  写入message queue ,发送手动处理消息
	message := NewStepMessage(dbStep.ExecutionID, MessageType.StateNewTurn, dbStep.ID, nil)
	err = svc.SendInnerMessage(message, nil)
	if err != nil {
		return err
	}
	return nil
}

func (svc *executionService) ResumeExecution(execution_id int) error {

	var err error
	var dbexecution po.Execution
	var dbsteps []po.Step

	// 加Execution 锁
	// L3 包含 definition

	dbexecution, err = svc.QueryExecutionByID(execution_id, ExecutionFields.L1, nil)
	if err != nil {
		return err
	}

	lock := svc.LockService.LockExecution(dbexecution.ID)
	err = lock.Lock()
	if err != nil {
		return err
	}
	defer lock.Unlock()

	// Execution 是Failed /Running 状态
	if !slices.Contains([]string{string(ExecutionStatus.Suspending)},
		dbexecution.Status) {
		err = fmt.Errorf("execution [ %s ]  Status Should Not  Be [ %s ]", dbexecution.UUID, dbexecution.Status)
		return err
	}
	// 开始事务
	tx, maker := svc.MetaDB.NewTxMaker(nil)
	defer maker.Close(&err)
	starttime := time.Now()

	suspendstepcond := po.Step{
		ExecutionID: dbexecution.ID,
		Status:      string(ExecutionStatus.Suspending),
		Type:        string(states.StateTypes.Suspend),
	}
	err = tx.Where(suspendstepcond).Select(StepFields.L2).Find(&dbsteps).Error
	if err != nil {
		return err
	}

	var events []ExecutionEvent
	var messages []queue.InnerMessage

	for _, dbstep := range dbsteps {

		exestep, err := svc.NewStepFromID(dbstep.ID, tx)
		if err != nil {
			return err
		}
		ns, err := exestep.GetNextStep(nil)
		if err != nil {
			return err
		}
		output, err := states.ToString(ns.Output)
		if err != nil {
			return err
		}
		// 修改成功
		finishtime := time.Now()
		updatestate := po.Step{
			Status:     string(StepStatus.Success),
			FinishTime: &finishtime,
			Output:     output,
		}

		// 更新step表，suspend 步骤执行成功
		err = tx.Where(po.Step{ID: dbstep.ID}).Updates(&updatestate).Error
		if err != nil {
			return err
		}
		err = svc.ChangeMasterGroupStatus(dbstep.ExecutionID, dbstep.GroupID, ExecutionStatus.Running, tx)
		if err != nil {
			return err
		}

		fns := FindNextStep{
			Name:    ns.Name,
			GroupID: dbstep.GroupID,
		}

		event := vo.ExecutionEvent{
			ExecutionID: dbstep.ExecutionID,
			StepID:      dbstep.ID,
			StepName:    dbstep.Name,
			StartTime:   starttime,
			FinishTime:  finishtime,

			Data: EventContent_SuspendStepResume{},
		}

		events = append(events, event)
		message := NewStepMessage(dbstep.ExecutionID, MessageType.FindNextStep, dbstep.ID, fns)
		messages = append(messages, message)

	}
	// 更新execution 表
	err = svc.ChangeExecutionStatus(dbexecution.ID, ExecutionStatus.Running, tx)
	if err != nil {
		return err
	}

	tx.Commit()
	finishtime := time.Now()
	event2 := vo.ExecutionEvent{
		ExecutionID: dbexecution.ID,
		StartTime:   starttime,
		FinishTime:  finishtime,
		StepID:      0,
		StepName:    "",
		Data: EventContent_ExecutionInfoModified{
			Field:  DataModifyField.Status,
			Before: dbexecution.Status,
			After:  string(ExecutionStatus.Running),
		},
	}

	events = append(events, event2)

	svc.SendExecutionEvents(events...)

	for _, message := range messages {
		err = svc.SendInnerMessage(message, time.Now())
		if err != nil {
			return err
		}
	}
	return nil

}

// ResumeSuspendingStep 继续一个暂停中的步骤
// @step_id 要继续的步骤
// NOCC:golint/fnsize("设计如此")
func (svc *executionService) ResumeSuspendingStep(step_id int) error {

	var err error
	var dbExecution *po.Execution
	var dbstep *po.Step

	// 加Execution 锁
	// L3 包含 definition
	dbstep, err = svc.QueryStepByID(step_id, StepFields.L1, nil)
	if err != nil {
		return err
	}

	lock := svc.LockService.LockExecution(dbstep.ExecutionID)
	err = lock.Lock()
	if err != nil {
		return err
	}
	defer lock.Unlock()

	exestep, err := svc.NewStepFromID(dbstep.ID, nil)
	if err != nil {
		return err
	}

	ns, err := exestep.GetNextStep(nil)
	if err != nil {
		return err
	}
	output, err := states.ToString(ns.Output)
	if err != nil {
		return err
	}

	// 开始事务
	tx, maker := svc.MetaDB.NewTxMaker(nil)
	defer maker.Close(&err)
	starttime := time.Now()

	// L3 包含 definition
	dbstep, err = svc.QueryStepByID(step_id, StepFields.L5, tx)
	if err != nil {
		return err
	}
	dbExecution, err = svc.QueryExecutionByID(dbstep.ExecutionID, ExecutionFields.L1, tx)
	if err != nil {
		return err
	}

	// Step是 Suspend 类型，并且 状态是 Suspending 才行
	if !(dbstep.Status == string(StepStatus.Suspending) && dbstep.Type == states.StateType.Suspend) {
		err = fmt.Errorf("step name  [ %s ]  expect  type [ %s ] status [ %s ] , actually status [ %s ] , type [ %s ] Found",
			dbstep.Name, states.StateType.Suspend, string(StepStatus.Suspending), dbstep.Type, dbstep.Status)
		return err
	}

	// Execution 是Failed /Running 状态
	if !slices.Contains([]string{string(ExecutionStatus.Suspending), string(ExecutionStatus.Running)},
		dbExecution.Status) {
		err = fmt.Errorf("execution [ %s ]  Status Should Not  Be [ %s ]", dbExecution.UUID, dbExecution.Status)
		return err
	}

	// 更新execution 表
	err = svc.ChangeExecutionStatus(dbstep.ExecutionID, ExecutionStatus.Running, tx)
	if err != nil {
		return err
	}

	// 修改成功
	finishtime := time.Now()
	updatestate := po.Step{
		Status:     string(StepStatus.Success),
		FinishTime: &finishtime,
		Output:     output,
	}

	// 更新step表，suspend 步骤执行成功
	err = tx.Where(po.Step{ID: dbstep.ID}).Updates(&updatestate).Error
	if err != nil {
		return err
	}
	err = svc.ChangeMasterGroupStatus(dbstep.ExecutionID, dbstep.GroupID, ExecutionStatus.Running, tx)
	if err != nil {
		return err
	}

	//更新DB完成
	tx.Commit()

	event1 := vo.ExecutionEvent{
		ExecutionID: dbstep.ExecutionID,
		StepID:      dbstep.ID,
		StepName:    dbstep.Name,
		StartTime:   starttime,
		FinishTime:  finishtime,

		Data: EventContent_SuspendStepResume{},
	}

	event2 := vo.ExecutionEvent{
		ExecutionID: dbstep.ExecutionID,
		StartTime:   starttime,
		FinishTime:  finishtime,
		StepID:      dbstep.ID,
		StepName:    dbstep.Name,
		Data: EventContent_ExecutionInfoModified{
			Field:  DataModifyField.Status,
			Before: dbExecution.Status,
			After:  string(ExecutionStatus.Running),
		},
	}

	svc.SendExecutionEvents(event1, event2)

	fns := FindNextStep{
		Name:    ns.Name,
		GroupID: dbstep.GroupID,
	}

	// 写入message queue ,处理人工干预消息
	message := NewStepMessage(dbstep.ExecutionID, MessageType.FindNextStep, step_id, fns)
	err = svc.SendInnerMessage(message, nil)
	if err != nil {
		return err
	}
	return nil
}

// SendStepRetry 重试Step
// NOCC:golint/fnsize("设计如此")
func (svc *executionService) SendStepRetry(step_id int) error {

	var err error
	var dbStep *po.Step
	var dbExecution *po.Execution

	// 开始事务
	tx, maker := svc.MetaDB.NewTxMaker(nil)
	defer maker.Close(&err)
	// 加锁
	txf := rdb.ForUpdate(tx)

	starttime := time.Now()

	// L3 包含 definition
	dbStep, err = svc.QueryStepByID(step_id, StepFields.L2, txf)
	if err != nil {
		return err
	}
	dbExecution, err = svc.QueryExecutionByID(dbStep.ExecutionID, ExecutionFields.L1, tx)
	if err != nil {
		return err
	}

	if !slices.Contains([]string{string(ExecutionStatus.Failed), string(ExecutionStatus.Running)},
		dbExecution.Status) {
		err = fmt.Errorf("%w: current execution status ' %s'", vo.ErrorExecutionStatus, dbExecution.Status)
		return err
	}
	if dbStep.Status != string(StepStatus.Failed) {
		err = fmt.Errorf("%w: current execution status ' %s'", vo.ErrorStepStatus, dbStep.Status)
		return err
	}
	if dbStep.Type == string(states.StateTypes.Parallel) || dbStep.Type == string(states.StateTypes.Map) || dbStep.Type == string(states.StateTypes.StateGroup) {
		err = fmt.Errorf("%w: operation: 'Retry' , step type '%s'", vo.ErrorUnsupportOperationForStep, dbStep.Type)
		return err
	}

	// 清理step token
	err = svc.CleanStepToken(step_id, tx)
	if err != nil {
		return err
	}

	// 更新step 状态
	err = svc.ChangeStepStatus(step_id, StepStatus.WaitInit, tx)
	if err != nil {
		return err
	}

	// 	递归修改组状态
	err = svc.ChangeMasterGroupStatus(dbStep.ExecutionID, dbStep.GroupID, ExecutionStatus.Running, tx)
	if err != nil {
		return err
	}

	// 更新execution 状态
	err = svc.ChangeExecutionStatus(dbStep.ExecutionID, ExecutionStatus.Running, tx)
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
		Data: EventContent_StepRetry{
			ExecuteCount: dbStep.ExecuteCount,
		},
	}

	svc.SendExecutionEvents(event)

	msg := StepExecuteMessage{
		Block: true,
	}
	//  写入message queue ,发送手动处理消息
	message := NewStepMessage(dbStep.ExecutionID, MessageType.StateNewTurn, dbStep.ID, msg)
	err = svc.SendInnerMessage(message, nil)
	if err != nil {
		return err
	}
	return nil
}

// CleanStepToken 清理Step 可能存在的token
// 修改状态
// 1. 修改本身的状态，
// 2. 删除其他的关联信息： activitytask/tasktoken
// 3. 修改子组的状态
func (svc *executionService) CleanStepToken(step_id int, tx rdb.Tx) error {

	var err error

	// 开始事务
	tx, maker := svc.MetaDB.NewTxMaker(nil)
	defer maker.Close(&err)

	//deleted state token
	err = tx.Where(po.TaskToken{StepID: step_id}).Delete(new(po.TaskToken)).Error
	if err != nil {
		return err
	}
	//deleted activity task
	err = tx.Where(po.ActivityTask{StepID: step_id}).Delete(new(po.ActivityTask)).Error
	if err != nil {
		return err
	}

	return nil
}

// ChangeStepStatus 修改Step状态
// 级联修改，修改一个步骤的状态已经相关的所在组的状态， 递归修改一来路径上所有的状态。
// @step_id  步骤ID
// @status 状态
// @session
func (svc *executionService) ChangeStepStatus(step_id int, status _StepStatusType, tx rdb.Tx) error {

	var err error
	// 开始事务
	tx, maker := svc.MetaDB.NewTxMaker(tx)
	defer maker.Close(&err)

	updatestate := po.Step{
		Status: string(status),
	}

	// 更新step 状态
	err = tx.Where(po.Step{ID: step_id}).Updates(&updatestate).Error
	if err != nil {
		return err
	}
	return nil
}

// ChangeMasterGroupStatus 修改父节点组的状态。递归的发现父节点所有组的状态并且返回
// 找到当前组的所有父节点，并且递归上报组的状态
// 输入参数 :
// @execution_id: 当前execution_id
// @group_id :当前所处group_id
// @status: 目标状态 status
func (svc *executionService) ChangeMasterGroupStatus(execution_id int, group_id int, status _ExecutionStatusType, tx rdb.Tx) error {
	// 如果是最外层节点，不需要修改
	if group_id == states.StartGroupID {
		return nil
	}
	// 说明不是最外层节点，需要修改外层节点状态成 Running
	dbstepgroup, err := svc.QueryStepGroupBySubGroupID(execution_id, group_id, tx)
	if err != nil {
		return err
	}
	err = svc.ChangeStepGroupStatus(dbstepgroup.StepID, status, tx)
	if err != nil {
		return err
	}
	return nil
}

// ChangeStepStatus 修改Step状态
// 级联修改，修改一个步骤的状态已经相关的所在组的状态， 递归修改一来路径上所有的状态。
// @step_id  步骤ID
// @status 状态
// @session
func (svc *executionService) ChangeStepGroupStatus(step_id int, status _ExecutionStatusType, tx rdb.Tx) error {

	var err error
	var dbStep *po.Step
	var dbStepGroup *po.StepGroup
	var curstep_id = step_id

	// 开始事务
	tx, maker := svc.MetaDB.NewTxMaker(tx)
	defer maker.Close(&err)

	updatestate := po.Step{
		Status: string(status),
	}

	for {
		dbStep, err = svc.QueryStepByID(curstep_id, StepFields.L2, tx)
		if err != nil {
			return err
		}
		// 如果已经更新了。则说明不需要更新， 说明上层都是Running 预期，可以直接返回
		if dbStep.Status == string(status) {
			return nil
		}
		// 更新step 状态
		err = tx.Where(po.Step{ID: curstep_id}).Updates(&updatestate).Error
		if err != nil {
			return err
		}
		// 如果发现自己在最外层， 就可以返回了
		if dbStep.GroupID == states.StartGroupID {
			break
		}
		// 如果自身是StepGroup， 算出自身的MasterID步骤,更新MasterStep的状态
		if dbStep.Type == states.StateTypes.StateGroup {

			// 找到MasterStepID
			dbStepGroup, err = svc.QueryStepGroupByStepID(dbStep.ID, tx)
			if err != nil {
				return err
			}
			curstep_id = dbStepGroup.MasterStepID
		} else {
			// 找到StepGroupID
			// 如果不是stepgroup ,计算出当前的stepgroup
			stepgroupcond := po.StepGroup{
				ExecutionID: dbStep.ExecutionID,
				SubGroupID:  dbStep.GroupID,
			}
			err = tx.Where(stepgroupcond).Take(&dbStepGroup).Error
			if err != nil {
				return err
			}
			curstep_id = dbStepGroup.StepID
		}

		//如果不是running , 就结束了
		// 如果是running ，递归到更高的组
		if status != ExecutionStatus.Running {
			return nil
		}
	}

	return nil
}

func (svc *executionService) ChangeExecutionStatus(execution_id int, status _ExecutionStatusType, tx rdb.Tx) error {

	var err error
	// 开始事务
	tx, maker := svc.MetaDB.NewTxMaker(tx)
	defer maker.Close(&err)

	updateexecution := po.Execution{
		Status: string(status),
	}
	err = tx.Where(po.Execution{ID: execution_id}).Updates(&updateexecution).Error
	if err != nil {
		return err
	}
	return err
}

// GetActivityTask 检索待执行的Activity
/*
	参数:
		uri : 检查的activity 的 uri
	返回:
		NewActivityTaskStruct 待执行获得的信息。
		error

*/
// NOCC:golint/fnsize("设计如此")
func (svc *executionService) GetActivityTask(ctx context.Context, req vo.GetActivityTaskRequest) (resp vo.GetActivityTaskResponse, err error) {
	var activitytaskid int
	var found = false

	var dbTask *po.Step
	var dbExecution *po.Execution
	var atData ActivityTaskData
	var dbActivityTask po.ActivityTask

	starttime := time.Now()
	requestinfo := vo.GetRequestInfo(ctx)

	cm := svc.cachemap

	// 增加控制session
	tx, maker := svc.MetaDB.NewTxMaker(nil)
	defer maker.Close(&err)

	for i := 0; i < 3; i++ {
		// 最大尝试3次

		activitytaskid = cm.Pop(req.ActivityURI)
		if activitytaskid == 0 {
			continue
		}
		//加锁查询

		// 加排它锁， 处理activity task
		err = rdb.ForUpdate(tx).Where(po.ActivityTask{ID: int64(activitytaskid)}).Take(&dbActivityTask).Error
		// 如果发现不存在，则忽略，尝试下一个
		if rdb.IsErrRecordNotFound(err) {
			continue
		}
		if err != nil {
			return
		}
		err = tx.Where(po.Step{ID: dbActivityTask.StepID}).Select(append(StepFields.L1, "execute_count")).Take(&dbTask).Error
		if err != nil {
			return
		}
		err = tx.Where(po.Execution{ID: dbActivityTask.ExecutionID}).Select(ExecutionFields.L1).Take(&dbExecution).Error
		if err != nil {
			return
		}
		if !(dbTask.Status == string(StepStatus.Wait) &&
			dbExecution.Status == string(ExecutionStatus.Running)) {
			err = tx.Where(po.ActivityTask{ID: dbActivityTask.ID}).Delete(new(po.ActivityTask)).Error
			if err != nil {
				tx.Rollback()
				return
			}
			continue
		}
		// 到这里说明都满足条件了。 就算找到了
		found = true
		break
	}
	// 超过查找次数没找到
	if !found {
		tx.Commit()
		err = fmt.Errorf("%w : %s", vo.ErrorActivityTaskNotFound, req.ActivityURI)
		return
	}
	err = json.Unmarshal([]byte(dbActivityTask.Data), &atData)
	if err != nil {
		return
	}

	updatetask := po.Step{
		Status: string(StepStatus.Running),
	}

	err = tx.Where(po.Step{ID: dbActivityTask.StepID}).Updates(&updatetask).Error

	if err != nil {
		tx.Rollback()
		return
	}

	err = tx.Where(po.ActivityTask{ID: dbActivityTask.ID}).Delete(new(po.ActivityTask)).Error
	if err != nil {
		return
	}
	tx.Commit()

	//Event
	finishtime := time.Now()

	event := vo.ExecutionEvent{
		ExecutionID: dbActivityTask.ExecutionID,
		StepID:      dbActivityTask.StepID,
		StepName:    dbTask.Name,
		StartTime:   starttime,
		FinishTime:  finishtime,
		Data: EventContent_ActivityScheduled{
			Input:       dbActivityTask.Input,
			Resource:    dbActivityTask.Resource,
			Timeout:     atData,
			RequestInfo: requestinfo,
		},
	}

	svc.SendExecutionEvents(event)

	// 发送超时消息
	mec := TaskTimeoutMessage{
		HeartbeatCount: atData.HeartbeatCount,
		ExecuteCount:   dbTask.ExecuteCount,
		TaskToken:      dbActivityTask.Token,
	}
	// send mq
	// 如果设置有有效的超时时间的话，发送超时消息
	if atData.TimeoutSeconds > 0 {
		// send timing message 发送定时消息,
		// 发送超时事件， 准备做超时处理。
		timouttime := time.Now().Add(time.Second * time.Duration(atData.TimeoutSeconds))

		tasktimeoutMsg := NewStepMessage(dbActivityTask.ExecutionID,
			MessageType.TaskTimeoutTimeup, dbTask.ID, mec)
		err = svc.SendInnerMessage(tasktimeoutMsg, &timouttime)
		if err != nil {
			return
		}
	}

	// 如果设置有有效的心跳时间的话，发送心跳消息
	if atData.HeartbeatSeconds > 0 {
		// send timing message 发送定时消息,
		// 发送超时事件， 准备做超时处理。
		timouttime := time.Now().Add(time.Second * time.Duration(atData.HeartbeatSeconds))
		tasktimeoutMsg := NewStepMessage(dbActivityTask.ExecutionID,
			MessageType.TaskHeartBeatTimeup, dbTask.ID, mec)
		err = svc.SendInnerMessage(tasktimeoutMsg, &timouttime)
		if err != nil {
			return
		}
	}

	resp = vo.GetActivityTaskResponse{
		Step:             dbTask,
		Execution:        dbExecution,
		Input:            dbActivityTask.Input,
		TaskToken:        dbActivityTask.Token,
		TimeoutSeconds:   atData.TimeoutSeconds,
		HeartbeatSeconds: atData.HeartbeatSeconds,
	}
	return
}

// SendTaskSuccess send success info of a state
/* @token : task token
 * @output: 执行输出
 */
func (svc *executionService) SendTaskSuccess(ctx context.Context, req vo.SendTaskSuccessRequest) error {
	slog.Info(fmt.Sprintf("SendTaskSuccess token: [ %s ] output length: %d", req.TaskToken, len(req.Output)))
	requestinfo := vo.GetRequestInfo(ctx)
	var outputs interface{}
	outputs, err := toolkit.ToJSON(req.Output)
	if err != nil {
		return err
	}
	step, err := svc.NewTaskFromToken(req.TaskToken, nil)
	if err != nil {
		return err
	}

	err = svc.JudgeExecutionRunningStatus(step.Data.ExecutionID)
	if err != nil {
		return err
	}

	err = step.SendTaskSuccess(outputs, requestinfo)
	return err

}

// SendTaskFailure send failure info of a state
func (svc *executionService) SendTaskFailure(ctx context.Context, req vo.SendTaskFailureRequest) error {

	requestinfo := vo.GetRequestInfo(ctx)

	step, err := svc.NewTaskFromToken(req.TaskToken, nil)
	if err != nil {
		return err
	}
	err = svc.JudgeExecutionRunningStatus(step.Data.ExecutionID)
	if err != nil {
		return err
	}

	err = step.SendTaskFailure(req.Error, req.Cause, requestinfo, nil)
	if err != nil {
		return err
	}

	return err

}

// SendTaskRefer send task refer information
func (svc *executionService) SendTaskReference(ctx context.Context, req vo.SendTaskReferenceRequest) error {

	requestinfo := vo.GetRequestInfo(ctx)

	step, err := svc.NewTaskFromToken(req.TaskToken, nil)
	if err != nil {
		return err
	}

	err = svc.JudgeExecutionRunningStatus(step.Data.ExecutionID)
	if err != nil {
		return err
	}

	err = step.SendTaskReference(req.Title, req.URL, requestinfo)
	if err != nil {
		return err
	}

	return err
}

// SendTaskHeartbeat send heartbeat of a state
func (svc *executionService) SendTaskHeartbeat(ctx context.Context, req vo.SendTaskHeartbeatRequest) error {

	var err error

	step, err := svc.NewTaskFromToken(req.TaskToken, nil)
	if err != nil {
		return err
	}
	err = svc.JudgeExecutionRunningStatus(step.Data.ExecutionID)
	if err != nil {
		return err
	}
	err = step.SendTaskHeartbeat(ctx, req.Message, nil)
	if err != nil {
		return err
	}
	return err
}

// RetryExecution    Retry all failed state in  execution
/* @uuid  execution uuid
 */
// NOCC:golint/fnsize("设计如此")
func (svc *executionService) RetryExecution(execution_id int) error {

	var err error
	starttime := time.Now()

	var dbexecution = po.Execution{}

	// 加锁
	lock := svc.LockService.LockExecution(execution_id)
	err = lock.Lock()
	if err != nil {
		return err
	}
	defer lock.Unlock()

	// 开始事务
	tx, maker := svc.MetaDB.NewTxMaker(nil)
	defer maker.Close(&err)

	dbexecution, err = svc.QueryExecutionByID(execution_id, ExecutionFields.L1, tx)
	if err != nil {
		return err
	}

	if dbexecution.Status != string(ExecutionStatus.Failed) {
		err := fmt.Errorf("%w : current execution status '%s' ", vo.ErrorExecutionStatus, dbexecution.Status)
		return err
	}
	var dbsteps []po.Step
	var condstate = po.Step{
		ExecutionID: dbexecution.ID,
		Status:      string(StepStatus.Failed),
	}

	err = tx.Model(po.Step{}).Select(StepFields.L2).Find(&dbsteps, condstate).Error
	if err != nil {
		return err
	}
	// 需要处理的失败节点列表
	failedstates := []po.Step{}
	for _, dbstep := range dbsteps {
		if dbstep.Type == states.StateType.Parallel || dbstep.Type == states.StateType.Map || dbstep.Type == states.StateType.StateGroup {
			continue
		}
		failedstates = append(failedstates, dbstep)
	}

	if len(failedstates) == 0 {
		return fmt.Errorf("no failed states found")
	}

	var events []ExecutionEvent
	var messages []queue.InnerMessageBody

	event := vo.ExecutionEvent{
		ExecutionID: dbexecution.ID,
		StartTime:   starttime,
		FinishTime:  time.Now(),
		Data:        EventContent_ExecutionRetry{},
	}
	events = append(events, event)

	// 更新execution 状态
	err = svc.ChangeExecutionStatus(dbexecution.ID, ExecutionStatus.Running, tx)
	if err != nil {
		return err
	}

	for _, dbstep := range failedstates {

		// 重试
		// 使用子流程修改状态，更加完善
		err = svc.ChangeStepStatus(dbstep.ID, StepStatus.WaitInit, tx)
		if err != nil {
			return err
		}
		err = svc.CleanStepToken(dbstep.ID, tx)
		if err != nil {
			return err
		}
		// 	递归修改组状态
		err = svc.ChangeMasterGroupStatus(dbstep.ExecutionID, dbstep.GroupID, ExecutionStatus.Running, tx)
		if err != nil {
			return err
		}
		finishtime := time.Now()
		event := vo.ExecutionEvent{
			ExecutionID: dbexecution.ID,
			StepID:      dbstep.ID,
			StepName:    dbstep.Name,
			StartTime:   starttime,
			FinishTime:  finishtime,
			Data: EventContent_StepRetry{
				ExecuteCount: dbstep.ExecuteCount,
			},
		}
		events = append(events, event)

		msg := StepExecuteMessage{
			Block: true,
		}

		message := NewStepMessage(dbstep.ExecutionID, MessageType.StateNewTurn, dbstep.ID, msg)
		messages = append(messages, message)
	}
	tx.Commit()

	svc.SendExecutionEvents(events...)

	// 写入message queue ,发送手动处理消息
	for _, message := range messages {
		err = svc.SendInnerMessage(message, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

// StoreTaskData 存储任务数据
func (svc *executionService) StoreTaskData(ctx context.Context, req vo.StoreTaskDataRequest) error {

	var err error

	starttime := time.Now()
	requestinfo := vo.GetRequestInfo(ctx)

	// 先查询step_id
	step_id, err := svc.QueryStepIDByTaskToken(req.TaskToken, nil)
	if err != nil {
		return err
	}

	// 开始事务
	tx, maker := svc.MetaDB.NewTxMaker(nil)
	defer maker.Close(&err)
	// 必须枷锁查询，保证数据的一致性
	txf := rdb.ForUpdate(tx)
	// 枷锁查询 step_id
	dbstep, err := svc.QueryStepByID(int(step_id), append(StepFields.L1, "data"), txf)
	if err != nil {
		return err
	}
	taskdata, err := LoadTaskData(dbstep.Data)
	if err != nil {
		return err
	}
	taskdata.TaskData = req.Data

	newdatastr, err := taskdata.String()
	if err != nil {
		return err
	}
	err = tx.Where(po.Step{ID: int(step_id)}).Updates(po.Step{Data: newdatastr}).Error
	if err != nil {
		return err
	}
	// Comment 立刻提交事务，关闭事务
	tx.Commit()

	finishtime := time.Now()
	// 发送事件
	event1 := vo.ExecutionEvent{
		ExecutionID: dbstep.ExecutionID,
		StepID:      dbstep.ID,
		StepName:    dbstep.Name,
		StartTime:   starttime,
		FinishTime:  finishtime,
		Data: EventContent_StoreTaskData{
			Data:        req.Data,
			RequestInfo: requestinfo,
		},
	}
	svc.SendExecutionEvents(event1)
	return nil
}

// LoadTaskData 加载任务数据
func (svc *executionService) LoadTaskData(ctx context.Context, step_id int) (data string, err error) {

	// 查询不开启事务
	dbstep, err := svc.QueryStepByID(int(step_id), []string{"id", "data"}, nil)
	if err != nil {
		return "", err
	}

	taskdata, err := LoadTaskData(dbstep.Data)
	if err != nil {
		return "", err
	}
	return taskdata.TaskData, nil
}

// SkipBlokcedTask 跳过阻塞的任务
func (svc *executionService) SkipBlockedTask(ctx context.Context, req vo.SkipBlokcedTaskRequest) error {

	var err error
	var dbStep *po.Step
	var dbNextStep *po.Step
	var dbExe *po.Execution
	starttime := time.Now()
	requestinfo := vo.GetRequestInfo(ctx)
	// 查询不开启事务
	dbStep, err = svc.QueryStepByID(req.StepID, StepFields.L1, nil)
	if err != nil {
		return err
	}

	tx, maker := svc.MetaDB.NewTxMaker(nil)
	defer maker.Close(&err)
	txf := rdb.ForUpdate(tx)
	dbExe, err = svc.QueryExecutionByID(dbStep.ExecutionID, ExecutionFields.L1, txf)
	if err != nil {
		return err
	}
	dbStep, err = svc.QueryStepByID(req.StepID, append(StepFields.L1, "group_id", "group_index"), nil)
	if err != nil {
		return err
	}
	if !(dbExe.Status == string(ExecutionStatus.Blocked) && dbStep.Status == string(StepStatus.Blocked)) {
		err = fmt.Errorf("%w: current execution status '%s' step status '%s'", vo.ErrorExecutionStatus, dbExe.Status, dbStep.Status)
		return err
	}

	nextstepquery := po.Step{
		ExecutionID: dbStep.ExecutionID,
		GroupID:     dbStep.GroupID,
		GroupIndex:  dbStep.GroupIndex + 1,
	}
	err = tx.Where(nextstepquery).Select(StepFields.L1).Take(&dbNextStep).Error
	if err != nil && !rdb.IsErrRecordNotFound(err) {
		return err
	}

	now := time.Now()
	updatestep := po.Step{
		Status:     string(StepStatus.Skip),
		Output:     "{}",
		FinishTime: &now,
	}
	err = tx.Where(po.Step{ID: dbStep.ID}).Updates(&updatestep).Error
	if err != nil {
		return err
	}
	err = svc.ChangeMasterGroupStatus(dbStep.ExecutionID, dbStep.GroupID, ExecutionStatus.Running, tx)
	if err != nil {
		return err
	}

	err = tx.Where(po.Execution{ID: dbExe.ID}).Updates(po.Execution{Status: string(ExecutionStatus.Running)}).Error
	if err != nil {
		return err
	}
	tx.Commit()

	fns := FindNextStep{
		Name:    dbNextStep.Name,
		GroupID: dbStep.GroupID,
	}

	// 写入message queue ,处理人工干预消息
	message := NewStepMessage(dbStep.ExecutionID, MessageType.FindNextState, req.StepID, fns)
	err = svc.SendInnerMessage(message, nil)
	if err != nil {
		return err
	}

	finishtime := time.Now()

	event1 := vo.ExecutionEvent{
		ExecutionID: dbStep.ExecutionID,
		StepID:      dbStep.ID,
		StepName:    dbStep.Name,
		StartTime:   starttime,
		FinishTime:  finishtime,

		Data: EventContent_SkipBlockedTask{
			NextStepName: "",
			Output:       "{}",
			RequestInfo:  requestinfo,
		},
	}

	event2 := vo.ExecutionEvent{
		ExecutionID: dbStep.ExecutionID,
		StartTime:   starttime,
		FinishTime:  finishtime,
		StepID:      dbStep.ID,
		StepName:    dbStep.Name,
		Data: EventContent_ExecutionInfoModified{
			Field:  DataModifyField.Status,
			Before: dbExe.Status,
			After:  string(ExecutionStatus.Running),
		},
	}

	svc.SendExecutionEvents(event1, event2)

	return nil
}

// UnblockTask 解除阻塞的任务
func (svc *executionService) UnblockTask(ctx context.Context, req vo.UnblockTaskRequest) error {

	var err error
	var dbstep *po.Step

	var dbexe *po.Execution
	starttime := time.Now()
	requestinfo := vo.GetRequestInfo(ctx)
	// 查询不开启事务
	dbstep, err = svc.QueryStepByID(req.StepID, StepFields.L1, nil)
	if err != nil {
		return err
	}

	tx, maker := svc.MetaDB.NewTxMaker(nil)
	defer maker.Close(&err)
	txf := rdb.ForUpdate(tx)
	dbexe, err = svc.QueryExecutionByID(dbstep.ExecutionID, ExecutionFields.L1, txf)
	if err != nil {
		return err
	}
	dbstep, err = svc.QueryStepByID(req.StepID, append(StepFields.L1, "group_id"), nil)
	if err != nil {
		return err
	}
	if !(dbexe.Status == string(ExecutionStatus.Blocked) && dbstep.Status == string(StepStatus.Blocked)) {
		err = fmt.Errorf("%w: current execution status '%s' step status '%s'", vo.ErrorExecutionStatus, dbexe.Status, dbstep.Status)
		return err
	}

	updatestep := po.Step{
		Status: string(StepStatus.Unblocking),
	}
	err = tx.Where(po.Step{ID: dbstep.ID}).Updates(&updatestep).Error
	if err != nil {
		return err
	}

	err = svc.ChangeMasterGroupStatus(dbstep.ExecutionID, dbstep.GroupID, ExecutionStatus.Running, tx)
	if err != nil {
		return err
	}

	err = tx.Where(po.Execution{ID: dbexe.ID}).Updates(po.Execution{Status: string(ExecutionStatus.Running)}).Error
	if err != nil {
		return err
	}
	tx.Commit()

	msgdata := StepExecuteMessage{
		Block: true,
	}
	msg := NewStepMessage(dbstep.ExecutionID, MessageType.StateExecute, dbstep.ID, msgdata)
	err = svc.SendInnerMessage(msg, nil)
	if err != nil {
		return err
	}

	finishtime := time.Now()
	event1 := vo.ExecutionEvent{
		ExecutionID: dbstep.ExecutionID,
		StepID:      dbstep.ID,
		StepName:    dbstep.Name,
		StartTime:   starttime,
		FinishTime:  finishtime,
		Data: EventContent_UnblockTask{
			RequestInfo: requestinfo,
		},
	}

	svc.SendExecutionEvents(event1)

	return nil
}

// UnblockExecution 解除阻塞的任务
func (svc *executionService) UnblockExecution(ctx context.Context, execution_id int) error {

	var err error
	starttime := time.Now()
	requestinfo := vo.GetRequestInfo(ctx)
	var dbexecution = &po.Execution{}

	// 加锁
	lock := svc.LockService.LockExecution(execution_id)
	err = lock.Lock()
	if err != nil {
		return err
	}
	defer lock.Unlock()

	// 开始事务
	tx, maker := svc.MetaDB.NewTxMaker(nil)
	defer maker.Close(&err)

	dbexecution, err = svc.QueryExecutionByID(execution_id, ExecutionFields.L1, tx)
	if err != nil {
		return err
	}

	if dbexecution.Status != string(ExecutionStatus.Blocked) {
		err := fmt.Errorf("%w : current execution status '%s' ", vo.ErrorExecutionStatus, dbexecution.Status)
		return err
	}
	var dbsteps []po.Step
	var condstate = po.Step{
		ExecutionID: dbexecution.ID,
		Status:      string(StepStatus.Blocked),
	}

	err = tx.Model(po.Step{}).Select(StepFields.L2).Find(&dbsteps, condstate).Error
	if err != nil {
		return err
	}
	// 需要处理的失败节点列表
	waitingtates := []po.Step{}
	for _, dbstep := range dbsteps {
		if slices.Contains(
			[]string{string(states.StateTypes.Parallel), string(states.StateTypes.Map), string(states.StateTypes.StateGroup)},
			dbstep.Type) {
			continue
		}

		waitingtates = append(waitingtates, dbstep)
	}

	if len(waitingtates) == 0 {
		return fmt.Errorf("no blocked states found")
	}

	var events []ExecutionEvent
	var messages []queue.InnerMessage

	event := vo.ExecutionEvent{
		ExecutionID: dbexecution.ID,
		StartTime:   starttime,
		FinishTime:  time.Now(),
		Data:        EventContent_UnblockExecution{},
	}
	events = append(events, event)

	// 更新execution 状态
	err = svc.ChangeExecutionStatus(dbexecution.ID, ExecutionStatus.Running, tx)
	if err != nil {
		return err
	}

	for _, dbstep := range waitingtates {

		// 重试
		// 使用子流程修改状态，更加完善
		err = svc.ChangeStepStatus(dbstep.ID, StepStatus.Unblocking, tx)
		if err != nil {
			return err
		}
		err = svc.CleanStepToken(dbstep.ID, tx)
		if err != nil {
			return err
		}
		// 	递归修改组状态
		err = svc.ChangeMasterGroupStatus(dbstep.ExecutionID, dbstep.GroupID, ExecutionStatus.Running, tx)
		if err != nil {
			return err
		}
		finishtime := time.Now()
		event := vo.ExecutionEvent{
			ExecutionID: dbexecution.ID,
			StepID:      dbstep.ID,
			StepName:    dbstep.Name,
			StartTime:   starttime,
			FinishTime:  finishtime,
			Data: EventContent_UnblockTask{
				RequestInfo: requestinfo,
			},
		}
		events = append(events, event)

		msg := StepExecuteMessage{
			Block: true,
		}

		message := NewStepMessage(dbstep.ExecutionID, MessageType.StateExecute, dbstep.ID, msg)
		messages = append(messages, message)
	}
	tx.Commit()

	svc.SendExecutionEvents(events...)

	// 写入message queue ,发送手动处理消息
	for _, message := range messages {
		err = svc.SendInnerMessage(message, time.Now())
		if err != nil {
			return err
		}
	}
	return nil
}
