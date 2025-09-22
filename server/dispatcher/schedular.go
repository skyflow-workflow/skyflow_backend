package dispatcher

import (
	"fmt"
	"time"

	"log/slog"

	"github.com/skyflow-workflow/skyflow_backbend/workflow/po"
)

// DispatcherService start schedular worker
func (svc *DispatcherService) StartSchedularWorkerManager() {
	// 同步组管理
	svc.receiverwg.Add(1)

	go func() {
		var err error
		// wg -1
		defer svc.receiverwg.Done()
		//从 innerqueue 接收消息,发送消息 workerpool，并发处理
		slog.Info("Dispatcher Message Receiver Start Running.")
		var message queue.InnerMessageInterface
		messagechan, err := svc.workflowService.Innerqueue.ReceiveInnerMessage()
		if err != nil {
			slog.Error("Acquire Inner Message Channel Failed:", err.Error())
			return
		}
		for {
			// 如果内存限制器存在，检查内存限制
			// 如果内存检查不通过，等待1秒
			if svc.memlimiter != nil && !svc.memlimiter.CheckAvailable() {
				slog.Error("current memory size: %d mb, limit size: %d mb",
					int(svc.memlimiter.GetCurrentMemoryByte()/1024/1024), int(svc.memlimiter.LimitSizeByte/1024/1024))
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
				svc.eventwg.Add(1)
				// 计数器 +1
				svc.IncreaseEventCounter()
			}
		}
	}()

}

func innerMsgLog(prefix string, msg queue.InnerMessage, err error) string {
	if err != nil {
		return fmt.Sprintf("%s, msg_id %s, ExecutionID %d, StepID %d, Priority %s, Type %s, ErrorDetail: %s",
			prefix, msg.ID, msg.ExecutionID, msg.StepID, msg.Class, msg.Type, err.Error())
	}
	return fmt.Sprintf("%s,msg_id %s, ExecutionID %d, StepID %d, Priority %s, Type %s",
		prefix, msg.ID, msg.ExecutionID, msg.StepID, msg.Class, msg.Type)
}

func loggermsg(prefix string, msg queue.InnerMessage) {
	msgstr := innerMsgLog(prefix, msg, nil)
	slog.Debug(msgstr)
}

// ProcessMessage 处理消息
func (svc *DispatcherService) ProcessMessage(i interface{}) {
	var err error
	// 同步组管理
	defer svc.eventwg.Done()
	starttime := time.Now()
	// 如果不是innermessage ，忽略消息
	message, ok := i.(queue.InnerMessageInterface)
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
		// 计算处理时间到监控系统
		processEventDuration.Set(time.Since(starttime).Seconds())
		// Ack Message
		err = message.Ack()
		if err != nil {
			errStr := innerMsgLog("Ack Message Failed: ", msgbody, err)
			slog.Error(errStr)
			// panic(fmt.Errorf(errStr))
		}
		if svc.option.Debug {
			finishtime := time.Now()
			duration := finishtime.Sub(starttime)
			jstr, _ := toolkit.ToString(msgbody)
			msg := fmt.Sprintf("ProcessMessage Speed  %s , Message Content %s", duration.String(), jstr)
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
func (svc *DispatcherService) processInnerMessage(msgbody queue.InnerMessage) (err error) {

	var dbstep po.Step
	var dbexecution po.Execution

	dbexecution, err = svc.workflowService.ExecutionService.QueryExecutionByID(msgbody.ExecutionID, execution.ExecutionFields.L1, nil)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	// 如果 消息类型是正常消息， 而 Execution状态不在 [ created  running ] , 忽略消息。 不能接收Failed/Abort/Success 等其他状态的消息
	if !toolkit.StringInSlice(
		[]string{string(execution.ExecutionStatus.Running), string(execution.ExecutionStatus.Created)},
		dbexecution.Status) {
		msg := fmt.Sprintf("Execution Status Is [ %s ] , Ignore This Message: ", dbexecution.Status)
		logStr := innerMsgLog(msg, msgbody, nil)
		slog.Info(logStr)
		return
	}
	// 根据消息类型做出响应的动作
	if msgbody.Class == queue.MessageClass.Execution {
		exe, err := svc.ExecutionService.NewExecutionFromID(msgbody.ExecutionID, execution.ExecutionFields.L1, nil)
		if err != nil {
			return err
		}
		// 如果是状态异常，忽略消息
		if err == execution.ErrorExecutionStatus {
			slog.Error(fmt.Sprintf("%w : %s ", err, msgbody))
			return nil
		}
		err = exe.ProcessEvent(msgbody)
		return err
	} else if msgbody.Class == queue.MessageClass.State {

		// 如果StateID == 0 , 忽略消息
		if msgbody.StepID == 0 {
			return
		}
		// 先查询最小数据，判断 step 的状态。是否符合预期。
		dbstep, err = svc.ExecutionService.QueryStepByID(msgbody.StepID, execution.StepFields.L1, nil)
		if err != nil {
			slog.Error(err.Error())
			return
		}

		//判断State 状态 是否符合预期， 增加对重复消息的可用性判断
		candState, ok := execution.StateEventEventCheckStatus[msgbody.Type]
		if ok {
			// 如果State状态 不在预期的event 状态， 忽略该消息
			if !toolkit.StringInSlice(candState, dbstep.Status) {
				jstr, _ := toolkit.ToString(msgbody)
				msg := fmt.Sprintf("Ignore This Message, StateID : %d , Current State: %s ,   Message Content: %s ",
					dbstep.ID, dbstep.Status, jstr)
				slog.Error(msg)
				return
			}
		}
		var step execution.Step
		// 查询全量数据， 初始化 step
		dbstep, err = svc.ExecutionService.QueryStepByID(msgbody.StepID, execution.StepFields.L5, nil)
		if err != nil {
			slog.Error(err.Error())
			return
		}
		step, err = execution.NewStepFromData(&dbstep, svc.workflowService.ExecutionService, dbexecution.FlowType)
		if err != nil {
			return err
		}

		if msgbody.Type == execution.MessageType.StateNewTurn {
			err = step.Init(msgbody)
		} else if msgbody.Type == execution.MessageType.StateExecute {
			err = step.Run(msgbody)
		} else if msgbody.Type == execution.MessageType.FindNextState {
			err = svc.workflowService.ExecutionService.ProcessFindNextState(msgbody)
		} else if msgbody.Type == execution.MessageType.ReportStepSuspend {
			err = svc.workflowService.ExecutionService.ProcessReportStepSuspend(msgbody)
		} else if msgbody.Type == execution.MessageType.ReportStepBlocked {
			err = svc.workflowService.ExecutionService.ProcessReportStepBlocked(msgbody)
		} else {
			err = step.ProcessEvent(msgbody)
		}
		if err != nil {

			err = svc.workflowService.ExecutionService.StepErrorProcess(err, &dbstep, msgbody)
		}
		return err
	}

	return nil

}
