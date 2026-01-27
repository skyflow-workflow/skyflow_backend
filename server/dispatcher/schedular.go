package dispatcher

import (
	"fmt"
	"slices"
	"time"

	"log/slog"

	"github.com/skyflow-workflow/skyflow_backend/workflow/executor"
	"github.com/skyflow-workflow/skyflow_backend/workflow/po"
	"github.com/skyflow-workflow/skyflow_backend/workflow/repository/queue"
)

// DispatcherService start schedular worker
func (svc *DispatcherService) StartSchedularWorkerManager() {
	// 同步组管理
	svc.receiverWg.Add(1)

	go func() {
		var err error
		// wg -1
		defer svc.receiverWg.Done()
		var message queue.InnerMessage
		for {
			// 如果内存限制器存在，检查内存限制
			// 如果内存检查不通过，等待1秒
			if svc.memlimiter != nil && !svc.memlimiter.CheckAvailable() {
				slog.Error("check  memory size:",
					"current_mb", int(svc.memlimiter.GetCurrentUsedMemoryByte()/1024/1024),
					"limit_mb", int(svc.memlimiter.LimitSizeByte/1024/1024))
				time.Sleep(1 * time.Second)
				continue
			}
			select {
			case <-svc.ctx.Done():
				slog.Error("Schedular Worker Manager Stopped")
				return
			case message = <-svc.msgChan:
				// wg +1 ,记录当前有运行的消息
				err = svc.workerPool.Invoke(message)
				if err != nil {
					slog.Error(err.Error())
				}
				svc.eventWg.Add(1)
				// 计数器 +1
				svc.IncreaseEventCounter()
			}
		}
	}()

}

func LogMessageEvent(msgBody queue.InnerMessageBody, err error) []any {

	args := []any{
		"ExecutionID", msgBody.ExecutionID,
		"StepID", msgBody.StepID,
		"Class", msgBody.Class,
		"Type", msgBody.Type,
		"Data", msgBody.Data,
	}
	if err != nil {
		args = append(args, "Error", err.Error())
	}
	return args
}

// ProcessMessage 处理消息
func (svc *DispatcherService) ProcessMessage(i interface{}) {
	var err error
	// 同步组管理
	defer svc.eventWg.Done()
	starttime := time.Now()
	// 如果不是innermessage ，忽略消息
	message, ok := i.(queue.InnerMessage)
	if !ok {
		return
	}
	msgbody := message.Body()
	// 如果execution_id  =0 , 非法消息，忽略
	if msgbody.ExecutionID == 0 {
		return
	}

	logEventFields := LogMessageEvent(msgbody, nil)
	logEventFields = append(logEventFields,
		"StartTime", starttime.Format(time.RFC3339Nano),
		"EventID", message.ID(),
	)

	// var has bool
	defer func() {
		slog.Info("ProcessMessage Finish", logEventFields...)
		// Ack Message
		err = message.Ack()
		if err != nil {
			slog.Error("Ack Message Failed", append(logEventFields, "Error", err)...)
			// panic(fmt.Errorf(errStr))
		}
		if svc.config.Debug {
			finishtime := time.Now()
			duration := finishtime.Sub(starttime)
			slog.Debug("ProcessMessage Speed", append(logEventFields, "Duration", duration.String())...)
		}
	}()
	defer func() {
		// 如果panic 了。 捕获panic
		if r := recover(); r != nil {
			slog.Error(fmt.Sprintf("%v", r))
			err = svc.workflowService.ExecutionService.ExecutionErrorProcess(err, msgbody)
			slog.Error("Panic Process Error", "error", err.Error())
		}
	}()
	// debug info
	slog.Debug("Start ProcessMessage", logEventFields...)
	// 交给具体执行函数
	err = svc.processInnerMessage(msgbody)
	slog.Info("Finish ProcessMessage", append(logEventFields, "Error", err)...)

	if err != nil {
		slog.Error("Process Message Failed Failed: ", append(logEventFields, "Error", err)...)
		// 处理不了就失败整个任务
		err2 := svc.ExecutionService.ExecutionErrorProcess(err, msgbody)
		if err2 != nil {
			slog.Error(err2.Error())
			return
		}

	}

}

// 调用step接口处理消息
// NOCC:golint/fnsize("设计如此")
func (svc *DispatcherService) processInnerMessage(msgBody queue.InnerMessageBody) (err error) {

	var dbStep *po.Step
	var dbExecution *po.Execution

	dbExecution, err = svc.workflowService.ExecutionService.QueryExecutionByID(msgBody.ExecutionID, executor.ExecutionFields.L1, nil)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	logEventFields := LogMessageEvent(msgBody, nil)

	// 如果 消息类型是正常消息， 而 Execution状态不在 [ created  running ] , 忽略消息。 不能接收Failed/Abort/Success 等其他状态的消息
	if !slices.Contains(
		[]string{string(executor.ExecutionStatus.Running), string(executor.ExecutionStatus.Created)},
		dbExecution.Status) {
		msg := fmt.Sprintf("Execution Status Is [ %s ] , Ignore This Message: ", dbExecution.Status)
		slog.Info(msg, logEventFields...)
		return
	}
	// 根据消息类型做出响应的动作
	switch msgBody.Class {
	case queue.MessageClass.Execution:
		exe, err := svc.ExecutionService.NewExecutionFromID(msgBody.ExecutionID, executor.ExecutionFields.L1, nil)
		if err != nil {
			return err
		}
		// 如果是状态异常，忽略消息
		if err == executor.ErrorExecutionStatus {

			slog.Error("Execution Status is unexpect , Ignore This Message.", append(logEventFields, "Error", err)...)
			return nil
		}
		err = exe.ProcessEvent(msgBody)
		return err
	case queue.MessageClass.Step:

		// 如果StateID == 0 , 忽略消息
		if msgBody.StepID == 0 {
			slog.Error("StepID is 0 , Ignore This Message.", append(logEventFields, "Error", err)...)
			return
		}
		// 先查询最小数据，判断 step 的状态。是否符合预期。
		dbStep, err = svc.ExecutionService.QueryStepByID(msgBody.StepID, executor.StepFields.L1, nil)
		if err != nil {
			slog.Error(err.Error())
			return
		}

		//判断State 状态 是否符合预期， 增加对重复消息的可用性判断
		candState, ok := executor.StateEventCheckStatus[msgBody.Type]
		if !ok {
			// 如果当前状态不在预期的候选列表中，说明消息类型异常，直接报错
			fields := append(logEventFields,
				"StepID", dbStep.ID,
				"StepName", dbStep.Name,
				"StepStatus", dbStep.Status,
				"StepType", dbStep.Type,
				"EventType", msgBody.Type,
			)

			err = fmt.Errorf("step event status is unexpect, event type is '%s', step status is '%s'",
				msgBody.Type, dbStep.Status)

			slog.Error(err.Error(), fields...)
			return err
		}

		// 如果State状态 不在预期的 event 状态， 有可能是重复消息， 忽略消息
		if !slices.Contains(candState, dbStep.Status) {
			err = fmt.Errorf(
				"step event status is not match, ingore this message, event type is '%s', step status is '%s'",
				msgBody.Type, dbStep.Status,
			)
			logEventFields = append(logEventFields,
				"Error", err,
				"StepID", dbStep.ID,
				"StepStatus", dbStep.Status,
				"StepType", dbStep.Type,
				"CandState", candState,
				"EventType", msgBody.Type,
			)
			slog.Error(err.Error(), logEventFields...)
			return
		}
		var step executor.Step
		// 查询全量数据， 初始化 step
		dbStep, err = svc.ExecutionService.QueryStepByID(msgBody.StepID, executor.StepFields.L5, nil)
		if err != nil {
			slog.Error(err.Error())
			return
		}
		step, err = executor.NewStepFromData(dbStep, svc.ExecutionService.StandardExecutor)
		if err != nil {
			slog.Error(err.Error())
			return err
		}

		switch msgBody.Type {
		case executor.MessageType.StateNewTurn:
			err = step.Init(msgBody)
		case executor.MessageType.StateExecute:
			err = step.Run(msgBody)
		case executor.MessageType.FindNextStep:
			err = svc.workflowService.ExecutionService.ProcessFindNextStep(msgBody)
		case executor.MessageType.ReportStepSuspend:
			err = svc.workflowService.ExecutionService.ProcessReportStepSuspend(msgBody)
		case executor.MessageType.ReportStepBlocked:
			err = svc.workflowService.ExecutionService.ProcessReportStepBlocked(msgBody)
		default:
			err = step.ProcessEvent(msgBody)
		}
		if err != nil {
			err = svc.workflowService.ExecutionService.StepErrorProcess(err, dbStep, msgBody)
		}
		return err
	}
	return nil

}
