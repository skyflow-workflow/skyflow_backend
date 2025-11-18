package dispatcher

import (
	"fmt"
	"slices"
	"time"

	"log/slog"

	"github.com/skyflow-workflow/skyflow_backbend/pkg/toolkit"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/executor"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/po"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/repository/queue"
)

// DispatcherService start schedular worker
func (svc *DispatcherService) StartSchedularWorkerManager() {
	// 同步组管理
	svc.receiverWg.Add(1)

	go func() {
		var err error
		// wg -1
		defer svc.receiverWg.Done()
		//从 innerqueue 接收消息,发送消息 workerpool，并发处理
		slog.Info("Dispatcher Message Receiver Start Running.")
		var message queue.InnerMessage
		messagechan, err := svc.workflowService.InnerQueue.ReceiveInnerMessage()
		if err != nil {
			slog.Error("Acquire Inner Message Channel Failed:",
				"error", err.Error())
			return
		}
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
			case message = <-messagechan:
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

func innerMsgLog(prefix string, msg queue.InnerMessageBody, err error) string {
	if err != nil {
		return fmt.Sprintf("%s, ExecutionID %d, StepID %d, Priority %s, Type %s, ErrorDetail: %s",
			prefix, msg.ExecutionID, msg.StepID, msg.Class, msg.Type, err.Error())
	}
	return fmt.Sprintf("%s, ExecutionID %d, StepID %d, Priority %s, Type %s",
		prefix, msg.ExecutionID, msg.StepID, msg.Class, msg.Type)
}

func loggermsg(prefix string, msg queue.InnerMessageBody) {
	msgstr := innerMsgLog(prefix, msg, nil)
	slog.Debug(msgstr)
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

	// var has bool
	defer func() {
		loggermsg("ProcessMessage Finish: ", msgbody)
		// Ack Message
		err = message.Ack()
		if err != nil {
			errStr := innerMsgLog("Ack Message Failed: ", msgbody, err)
			slog.Error(errStr)
			// panic(fmt.Errorf(errStr))
		}
		if svc.config.Debug {
			finishtime := time.Now()
			duration := finishtime.Sub(starttime)
			msgBodyStr, err := toolkit.ToString(msgbody)
			if err != nil {
				slog.Error("ProcessMessage Failed: ", "error", err.Error())
				return
			}
			msg := fmt.Sprintf("ProcessMessage Speed  %s , Message Content %s", duration.String(), msgBodyStr)
			slog.Debug(msg)
		}
	}()
	defer func() {
		// 如果panic 了。 捕获panic
		if r := recover(); r != nil {
			slog.Error(fmt.Sprintf("%v", r))
			err = svc.workflowService.ExecutionService.ExecutionErrorProcess(err, msgbody.ExecutionID)
			slog.Error(fmt.Sprintf("Panic Process Error %v", err))
		}
	}()
	// debug info
	slog.Info(fmt.Sprintf("Start ProcessMessage :  %v", msgbody))
	// 交给具体执行函数
	err = svc.processInnerMessage(msgbody)
	slog.Info(fmt.Sprintf("Finish ProcessMessage :  %v, err: %v", msgbody, err))

	if err != nil {
		errStr := innerMsgLog("Process Event ProcessFailed Failed: ", msgbody, err)
		slog.Error(errStr)
		// 处理不了就失败整个任务
		err2 := svc.ExecutionService.ExecutionErrorProcess(err, msgbody.ExecutionID)
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

	// 如果 消息类型是正常消息， 而 Execution状态不在 [ created  running ] , 忽略消息。 不能接收Failed/Abort/Success 等其他状态的消息
	if !slices.Contains(
		[]string{string(executor.ExecutionStatus.Running), string(executor.ExecutionStatus.Created)},
		dbExecution.Status) {
		msg := fmt.Sprintf("Execution Status Is [ %s ] , Ignore This Message: ", dbExecution.Status)
		logStr := innerMsgLog(msg, msgBody, nil)
		slog.Info(logStr)
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

			slog.Error("Execution Status is unexpect , Ignore This Message.",
				"ExecutionID", msgBody.ExecutionID,
				"ExecutionStatus", dbExecution.Status,
				"MessageType", msgBody.Type,
				"error", err,
			)
			return nil
		}
		err = exe.ProcessEvent(msgBody)
		return err
	case queue.MessageClass.Step:

		// 如果StateID == 0 , 忽略消息
		if msgBody.StepID == 0 {
			return
		}
		// 先查询最小数据，判断 step 的状态。是否符合预期。
		dbStep, err = svc.ExecutionService.QueryStepByID(msgBody.StepID, executor.StepFields.L1, nil)
		if err != nil {
			slog.Error(err.Error())
			return
		}

		//判断State 状态 是否符合预期， 增加对重复消息的可用性判断
		candState, ok := executor.StateEventEventCheckStatus[msgBody.Type]
		if ok {
			// 如果State状态 不在预期的event 状态， 忽略该消息
			if !slices.Contains(candState, dbStep.Status) {
				msgBodyStr, err2 := toolkit.ToString(msgBody)
				if err2 != nil {
					slog.Error("msgBody to json string failed",
						"error", err2)
					return err
				}
				msg := fmt.Sprintf("Ignore This Message, StateID : %d , Current State: %s ,   Message Content: %s ",
					dbStep.ID, dbStep.Status, msgBodyStr)
				slog.Error(msg)
				return
			}
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
