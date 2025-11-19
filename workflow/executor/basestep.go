package executor

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/mmtbak/microlibrary/rdb"
	"github.com/skyflow-workflow/skyflow_backend/pkg/toolkit"
	"github.com/skyflow-workflow/skyflow_backend/workflow/parser/states"
	"github.com/skyflow-workflow/skyflow_backend/workflow/repository/queue"
	"github.com/skyflow-workflow/skyflow_backend/workflow/vo"
	"trpc.group/trpc-go/tnet/log"

	"github.com/skyflow-workflow/skyflow_backend/workflow/po"
)

// ExecutionStep execution step
type ExecutionStep struct {
	Data             *po.Step
	ExecutionService ExecutionService
	Executor         *Executor
	State            states.State
	// BaseState        *states.BaseState
}

// NewBaseExecutionStep New一个可执行状态实例
func NewExecutionStep(dbStep *po.Step, executor *Executor) (*ExecutionStep, error) {
	var err error

	state, err := executor.Parser.ParseState(dbStep.Definition)
	if err != nil {
		return nil, err
	}
	exeStep := &ExecutionStep{
		State:            state,
		Executor:         executor,
		Data:             dbStep,
		ExecutionService: executor.ExecutionService,
	}
	return exeStep, nil

}

// GetBone  Get Execution Step Status and return StepBone  struct
func (step *ExecutionStep) GetBone() StepBone {

	bone := StepBone{
		BaseBone: BaseBone{
			BaseBone: step.State.GetBone().BaseBone,
			Status:   step.Data.Status,
			StepID:   step.Data.ID,
		},
	}
	if step.Data.ExecuteIndex > 0 {
		bone.Index = step.Data.ExecuteIndex
	}
	return bone
}

// ProcessEvent process step event
// ProcessEvent is a placeholder method for processing events.
// It currently returns an error indicating that the method is not implemented.
func (step *ExecutionStep) ProcessEvent(message queue.InnerMessageBody) error {
	return fmt.Errorf("method not implement")
}

// DecodeInput decode intput to json
func (step *ExecutionStep) DecodeInput() (any, error) {
	i, err := toolkit.DecodeStringToMap(step.Data.Input)
	return i, err
}

// DecodeOutput decode output to json
func (step *ExecutionStep) DecodeOutput() (any, error) {
	i, err := toolkit.DecodeStringToMap(step.Data.Output)
	return i, err
}

// GetInput 获得经过计算的input
func (step *ExecutionStep) GetInput() (any, error) {
	i, err := step.DecodeInput()
	if err != nil {
		return nil, err
	}
	ii, err := step.State.GetBaseState().GetParametersInput(i)
	return ii, err
}

// GetNextStep  return (NextState string , output any, err error )
func (step *ExecutionStep) GetNextStep(output any) (NextStep, error) {

	var err error
	var sns NextStep
	var input any
	input, err = step.DecodeInput()
	if err != nil {
		return sns, nil
	}
	ns, err := step.State.GetBaseState().GetNextState(input, output)
	if err != nil {
		return sns, err
	}
	nextstep := NextStep{
		NextState: ns,
		GroupID:   step.Data.GroupID,
	}
	return nextstep, nil

}

// Init Init初始化state
// NOCC:golint/fnsize("设计如此")
func (step *ExecutionStep) Init(msg queue.InnerMessageBody) error {

	var err error
	starttime := time.Now()

	var dbStep *po.Step
	var baseDBStep *po.Step
	var dbExecution *po.Execution

	var stateExeMsg = StepExecuteMessage{
		Block: false,
	}
	// 兼容历史消息
	if msg.Data != "" {
		err = json.Unmarshal([]byte(msg.Data), &stateExeMsg)
		if err != nil {
			return err
		}
	}
	// 提前取出 Input字段， 避免在事务中查询
	baseDBStep, err = step.ExecutionService.QueryStepByID(step.Data.ID, append(StepFields.L1, "input"), nil)

	if err != nil {
		log.Error(err.Error())
		return err
	}

	tx, maker := step.ExecutionService.MetaDB.NewTxMaker(nil)
	defer maker.Close(&err)

	// 如果打开了 ExecuteIndex 开关，需要计算每个步骤的ExecuteIndex
	if step.Executor.Config.Option.EnableExecuteIndex {
		// 要计算执行的Index , 需要加全局锁
		txf := rdb.ForUpdate(tx)
		dbExecution, err = step.ExecutionService.QueryExecutionByID(step.Data.ExecutionID, []string{"id", "max_execute_index"}, txf)
		if err != nil {
			return err
		}
	} else {
		// 减少一次额外的查询
		dbExecution, err = step.ExecutionService.QueryExecutionByID(dbStep.ExecutionID, []string{"id", "max_execute_index"}, tx)
		if err != nil {
			return err
		}
	}

	// 从数据库重新读取数据， 避免并发问题， 减少读取的字段数
	fields := []string{"id", "type", "execute_count", "execute_index", "execution_id"}
	dbStep, err = step.ExecutionService.QueryStepByID(step.Data.ID, fields, tx)

	if err != nil {
		log.Error(err.Error())
		return err
	}

	// 判断 是否超过 MaxExecuteTimes, 如果超过， 则拒绝执行
	if dbStep.ExecuteCount >= int(step.State.GetBaseState().MaxExecuteTimes) {
		err = fmt.Errorf("step execute count reach 'MaxExecuteTimes' argument")
		return err
	}
	// 开始准备初始化
	// 当前V1阶段
	// 不论何种情况，都先查一下
	var maxindex, executionindex, executecount int

	maxindex = dbExecution.MaxExecuteIndex
	// maxindex 代表当前最大的执行Index。
	// 如果当前是未执行节点，  则dbstep.ExecuteIndex = maxindex +1; maxindex = dbstep.ExecuteIndex
	// 如果当前是已经执行的节点，则 都保持不变 。
	if step.Executor.Config.Option.EnableExecuteIndex {

		//初始化 基础结构数据, 计算当前的ExecuteIndex
		// maxindex 是
		if dbStep.ExecuteIndex < 0 {
			if maxindex <= 0 {
				executionindex = 1
			} else {
				executionindex = maxindex + 1
			}
		} else {
			// 如果当前的已经 > 0
			// 则保持不变
			executionindex = dbStep.ExecuteIndex
		}
	} else {
		if dbStep.ExecuteIndex < 0 {
			executionindex = 0 - dbStep.ExecuteIndex
		} else {
			executionindex = dbStep.ExecuteIndex
		}
	}
	// 非Task 状态需要在处理流程中计数器 +1 ，
	// 因为task节点可以自动重试，可以在重试逻辑中自动累加计数器，所以需要在执行逻辑中 +1
	// 但是其他类型节点不需要，可以在全局初始化的时候 +1， 这些节点重试的时候，进入初始化流程中， 在这里再次 +1
	if dbStep.Type == string(states.StateTypes.Task) {
		executecount = dbStep.ExecuteCount
	} else {
		executecount = dbStep.ExecuteCount + 1
	}

	// update state information
	now := time.Now()
	updatestate := po.Step{
		Status:       string(StepStatus.Initialize),
		ExecuteIndex: executionindex,
		ExecuteCount: executecount,
		StartTime:    &now,
		Data:         "",
	}

	err = tx.Where(po.Step{ID: dbStep.ID}).Updates(&updatestate).Error
	if err != nil {
		return err
	}

	//如果打开了 ExecuteIndex 开关
	// 那么更新Execution的max_execute_index 值 ,不然就放弃，保持 0
	// 如果当前的index > maxindex , 更新maxindex
	if maxindex < executionindex {
		updateexe := po.Execution{
			MaxExecuteIndex: executionindex,
		}
		err = tx.Where(po.Execution{ID: dbStep.ExecutionID}).Updates(&updateexe).Error
		if err != nil {
			return err
		}
	}
	tx.Commit()

	event0 := vo.ExecutionEvent{
		ExecutionID: dbStep.ExecutionID,
		StepID:      dbStep.ID,
		StepName:    baseDBStep.Name,
		StartTime:   starttime,
		FinishTime:  now,
		Data: EventContent_StateEntered{
			Input: baseDBStep.Input,
		},
	}
	event1 := vo.ExecutionEvent{
		ExecutionID: dbStep.ExecutionID,
		StepID:      dbStep.ID,
		StepName:    baseDBStep.Name,
		StartTime:   starttime,
		FinishTime:  now,
		Data: EventContent_StateInit{
			Action:         "StartStep",
			ExecutionIndex: executionindex,
			ExecutionCount: updatestate.ExecuteCount,
		},
	}
	step.ExecutionService.SendExecutionEvents(event0, event1)

	// message queue send create message
	message := NewStepMessage(dbStep.ExecutionID, MessageType.StateExecute, dbStep.ID, stateExeMsg)
	err = step.ExecutionService.SendInnerMessage(message, nil)
	if err != nil {
		return err
	}
	return nil
}
