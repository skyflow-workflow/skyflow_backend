package executor

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/dolthub/vitess/go/vt/log"
	"github.com/skyflow-workflow/skyflow_backbend/pkg/toolkit"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/parser/states"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/po"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/repository/queue"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/vo"
)

// Wait ExecutionWaitState
type Wait struct {
	*ExecutionStep
	State *states.WaitState
}

// WaitData  State Data for Wait State
type WaitData struct {
	ExecuteCounter int `json:"execute_counter"`
}

// NewWaitFromID NewWaitFromID
func NewWaitFromID(id int, executor *Executor) (*Wait, error) {

	var err error

	var dbstep = &po.Step{}

	dbstep, err = executor.ExecutionService.QueryStepByID(id, []string{}, nil)
	if err != nil {
		return nil, err
	}

	state, err := NewWaitFromData(dbstep, executor)

	return state, err
}

// NewWaitFromData NewWaitFromData
func NewWaitFromData(dbStep *po.Step, executor *Executor) (*Wait, error) {

	baseStep, err := NewExecutionStep(dbStep, executor)
	if err != nil {
		return nil, err
	}
	state, ok := baseStep.State.(*states.WaitState)
	if !ok {
		return nil, fmt.Errorf("step state is not 'Wat' type")
	}
	waitState := &Wait{
		ExecutionStep: baseStep,
		State:         state,
	}
	return waitState, err
}

// GetWakeupTime GetWakeupTime
func (w *Wait) GetWakeupTime() (time.Time, error) {

	var err error
	var dest time.Time
	var input interface{}
	input, err = w.DecodeInput()
	if err != nil {
		return dest, err
	}
	dest, err = w.State.GetWakeupTime(input)
	return dest, err
}

// IsEnd IsEnd
func (w *Wait) IsEnd() bool {
	return w.State.IsEnd()
}

// ProcessEvent ProcessEvent
func (w *Wait) ProcessEvent(message queue.InnerMessageBody) error {

	var err error
	switch message.Type {
	case MessageType.WaitStateWakeup:
		err = w.WaitStateWakeup(message)
	default:
		log.Errorf("igore message: ", message)
	}
	return err
}

// GetNextState  return (NextState string , output interface{}, err error )
func (w *Wait) GetNextState() (NextStep, error) {

	var err error
	var sns NextStep
	var input interface{}
	input, err = w.DecodeInput()
	if err != nil {
		return sns, err
	}
	ns, err := w.State.GetNextState(input)
	if err != nil {
		return sns, err
	}
	sns = NextStep{
		NextState: ns,
		GroupID:   w.Data.GroupID,
	}
	return sns, nil
}

// Run Run
func (w *Wait) Run(_ queue.InnerMessageBody) error {
	var err error

	dbstate := w.Data

	starttime := time.Now()

	wakeupTime, err := w.GetWakeupTime()

	if err != nil {
		log.Error(err.Error())
		return err
	}
	now := time.Now()
	//
	// 如果醒来时间 < 当前时间， 说明当前时间已过， 直接成功
	if wakeupTime.Before(now) {
		err = w.Finish()
		if err != nil {
			return err
		}
		return nil
	}

	tx, maker := w.ExecutionService.MetaDB.NewTxMaker(nil)
	defer maker.Close(&err)
	var updatestep po.Step

	updatestep = po.Step{
		Status: string(StepStatus.Running),
	}
	err = tx.Where(po.Step{ID: dbstate.ID}).Updates(&updatestep).Error
	if err != nil {
		log.Error(err.Error())
		return err
	}
	tx.Commit()

	finishtime := time.Now()
	event := vo.ExecutionEvent{
		ExecutionID: dbstate.ExecutionID,
		StartTime:   starttime,
		FinishTime:  finishtime,
		StepID:      dbstate.ID,
		StepName:    dbstate.Name,
		Data: EventContent_WaitStateExecuted{
			ExecuteCount: dbstate.ExecuteCount,
			WakeupTime:   wakeupTime,
		},
	}
	w.ExecutionService.SendExecutionEvents(event)

	swm := StepWakeupMessage{
		ExecuteCount: dbstate.ExecuteCount,
	}
	// message queue send create message

	message := NewStepMessage(dbstate.ExecutionID, MessageType.WaitStateWakeup, dbstate.ID, swm)
	err = w.ExecutionService.SendInnerMessage(message, &wakeupTime)
	if err != nil {
		return err
	}
	return nil
}

// WaitStateWakeup   wait state wakeup
func (w *Wait) WaitStateWakeup(message queue.InnerMessageBody) error {

	var err error

	var dbstate = w.Data
	// var dbexecution *db.Execution

	/*
		当前wait节点当前不是Running ,说明状态被强行终止了。
		忽略这次计时器事件
	*/
	if dbstate.Status != string(StepStatus.Running) {
		return nil
	}

	mdec := StepWakeupMessage{}
	err = json.Unmarshal([]byte(message.Data), &mdec)
	if err != nil {
		return err
	}
	/*
		如果计数器与内部计数器不一致，说明此次计时器消息并不是当前正在执行的。
		说明当前当前状态已经被强制重试过。
		直接忽略历史计时器
	*/
	if mdec.ExecuteCount != dbstate.ExecuteCount {
		return nil
	}

	err = w.Finish()
	if err != nil {
		return err
	}
	return nil

}

func (w *Wait) Finish() error {
	var err error

	var dbstate = w.Data

	var starttime = time.Now()
	sns, err := w.GetNextState()
	if err != nil {
		return err
	}
	outputstr, err := toolkit.ToString(sns.Output)
	if err != nil {
		return err
	}

	now := time.Now()
	tx, maker := w.ExecutionService.MetaDB.NewTxMaker(nil)
	defer maker.Close(&err)

	updatestate := po.Step{
		Status:     string(StepStatus.Success),
		FinishTime: &now,
		Output:     outputstr,
	}

	err = tx.Where(po.Step{ID: dbstate.ID}).Updates(&updatestate).Error
	if err != nil {
		return err
	}
	tx.Commit()

	finishtime := time.Now()
	event1 := vo.ExecutionEvent{
		ExecutionID: dbstate.ExecutionID,
		StartTime:   starttime,
		FinishTime:  finishtime,
		StepID:      dbstate.ID,
		StepName:    dbstate.Name,
		Data:        EventContent_WaitStateWakeup{},
	}
	event2 := vo.ExecutionEvent{
		ExecutionID: dbstate.ExecutionID,
		StartTime:   starttime,
		FinishTime:  finishtime,
		StepID:      dbstate.ID,
		StepName:    dbstate.Name,
		Data: EventContent_StateExited{
			Output: outputstr,
		},
	}
	w.ExecutionService.SendExecutionEvents(event1, event2)

	fns := FindNextStep{
		Name:    sns.Name,
		GroupID: sns.GroupID,
	}

	// message queue send create message
	newmessage := NewStepMessage(dbstate.ExecutionID, MessageType.FindNextStep, dbstate.ID, fns)
	err = w.ExecutionService.SendInnerMessage(newmessage, nil)
	if err != nil {
		return err
	}
	return nil

}
