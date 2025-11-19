// package executor provides the main executor for workflow event processing.
package executor

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/mmtbak/microlibrary/rdb"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/config"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/parser"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/parser/states"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/po"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/repository/queue"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/vo"
)

// Executor is the main executor for workflow event processing.
// It is responsible for processing events and managing the execution of workflows.
// It uses the rdb.DBClient for database operations and queue.InnerMessageQueue for message transport.
// The Mode field indicates the execution mode, which can be set to "express" for faster processing.
// It is designed to be extensible, allowing for different execution strategies and configurations.
// The Executor struct encapsulates the necessary components for executing workflow events,
type Executor struct {
	ExecutionService *executionService
	Config           *config.Config
	Parser           *parser.Parser
}

// NewExecutor creates a new Executor instance
func NewExecutor(config *config.Config) *Executor {

	executor := &Executor{
		Config: config,
		Parser: parser.NewParser(config),
	}
	return executor
}

func (executor *Executor) ProcessEvent(state *po.Step, event queue.InnerMessageBody) error {
	return nil
}

// NewTaskFromToken NewTaskFromToken
func (executor *Executor) NewTaskFromToken(token string, fields []string, session rdb.Tx) (*Task, error) {

	var err error
	var tasktoken = po.TaskToken{
		Token: token,
	}
	var dbStep *po.Step

	// 增加控制session
	tx, maker := executor.ExecutionService.MetaDB.NewTxMaker(session)
	defer maker.Close(&err)

	err = tx.Where(tasktoken).Take(&tasktoken).Error
	if rdb.IsErrRecordNotFound(err) {
		return nil, fmt.Errorf("%w: %s", vo.ErrorTaskTokenNotFound, token)
	}
	if err != nil {
		return nil, err
	}

	dbStep, err = executor.ExecutionService.QueryStepByID(tasktoken.StepID, fields, tx)
	if err != nil {
		return nil, err
	}

	state, err := NewTaskFromData(dbStep, executor)
	if err != nil {
		return nil, err
	}
	return state, err
}

func (executor *Executor) SendInnerMessage(message queue.InnerMessageBody, sendtime *time.Time) error {

	return executor.ExecutionService.SendInnerMessage(message, sendtime)
}

// SendExecutionEvents 发送event
func (executor *Executor) SendExecutionEvents(events ...vo.ExecutionEvent) {
	executor.ExecutionService.Exporter.SendExecutionEvents(events)
}

// ProcessEventStepInit process step event 'Init'
// nolint: funlen
func (executor *Executor) ProcessEventStepInit(msg queue.InnerMessageBody) error {

	var err error
	starttime := time.Now()

	var dbStep *po.Step
	var dbExecution *po.Execution

	// 提前取出 Input字段， 避免在事务中查询
	dbStep, err = executor.ExecutionService.QueryStepByID(msg.StepID, append(StepFields.L2, StepFieldNames.Definition), nil)

	if err != nil {
		slog.Error(err.Error())
		return err
	}

	parserState, err := executor.Parser.ParseState(dbStep.Definition)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	// 判断 是否超过 MaxExecuteTimes, 如果超过， 则拒绝执行
	if dbStep.ExecuteCount >= int(parserState.GetBaseState().MaxExecuteTimes) {
		err = fmt.Errorf("step execute count reach 'MaxExecuteTimes' argument")
		return err
	}

	var EnableExecuteIndex = executor.Config.Option.EnableExecuteIndex

	tx, maker := executor.ExecutionService.MetaDB.NewTxMaker(nil)
	defer maker.Close(&err)

	// 如果打开了 ExecuteIndex 开关，需要计算每个步骤的ExecuteIndex
	if EnableExecuteIndex {
		// 要计算执行的Index , 需要加全局锁
		txf := rdb.ForUpdate(tx)
		dbExecution, err = executor.ExecutionService.QueryExecutionByID(msg.ExecutionID,
			[]string{ExecutionFieldNames.ID, ExecutionFieldNames.MaxExecuteIndex}, txf)
		if err != nil {
			return err
		}
	} else {
		// 减少一次额外的查询
		dbExecution, err = executor.ExecutionService.QueryExecutionByID(msg.ExecutionID,
			[]string{ExecutionFieldNames.ID, ExecutionFieldNames.MaxExecuteIndex}, tx)
		if err != nil {
			return err
		}
	}

	// 开始准备初始化
	// 先计算 ExecuteIndex 和 ExecuteCount
	var maxindex, executionindex, executecount int

	maxindex = dbExecution.MaxExecuteIndex
	// maxindex 代表当前最大的执行Index。
	// 如果当前是未执行节点，  则dbStep.ExecuteIndex = maxindex +1; maxindex = dbStep.ExecuteIndex
	// 如果当前是已经执行的节点，则 都保持不变 。
	if EnableExecuteIndex {
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

	if err != nil {
		slog.Error(err.Error())
		return err
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
		updateExe := po.Execution{
			MaxExecuteIndex: executionindex,
		}
		err = tx.Where(po.Execution{ID: dbStep.ExecutionID}).Updates(&updateExe).Error
		if err != nil {
			return err
		}
	}
	tx.Commit()

	event1 := vo.ExecutionEvent{
		ExecutionID: dbStep.ExecutionID,
		StepID:      dbStep.ID,
		StepName:    dbStep.Name,
		StartTime:   starttime,
		FinishTime:  now,
		Data: EventContent_StateEntered{
			Input: dbStep.Input,
		},
	}
	event2 := vo.ExecutionEvent{
		ExecutionID: dbStep.ExecutionID,
		StepID:      dbStep.ID,
		StepName:    dbStep.Name,
		StartTime:   starttime,
		FinishTime:  now,
		Data: EventContent_StateInit{
			Action:         "StartState",
			ExecutionIndex: executionindex,
			ExecutionCount: updatestate.ExecuteCount,
		},
	}
	executor.SendExecutionEvents(event1, event2)

	var stateExeMsg = StepExecuteMessage{
		Block: false,
	}
	// message queue send create message
	message := NewStepMessage(dbStep.ExecutionID, MessageType.StateExecute, dbStep.ID, stateExeMsg)
	err = executor.SendInnerMessage(message, nil)
	if err != nil {
		return err
	}
	return nil
}
