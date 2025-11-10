package executor

import (
	"fmt"
	"slices"
	"time"

	"github.com/mmtbak/microlibrary/rdb"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/cache"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/domain"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/exporter"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/po"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/repository/lock"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/repository/queue"
	"github.com/skyflow-workflow/skyflow_backbend/workflow/vo"
)

type ExecutionService = *executionService

type executionService struct {
	MetaDB           *rdb.DBClient
	InnerQueue       queue.InnerMessageQueue
	Exporter         exporter.ExporterService
	LockService      lock.LockService
	DomainService    domain.DomainService
	StandardExecutor *Executor
	ExpressExecutor  *Executor
	CacheService     cache.TaskCacheService
}

func NewExecutionService(
	MetaDB *rdb.DBClient,
	InnerQueue queue.InnerMessageQueue,
	Exporter exporter.ExporterService,
) ExecutionService {
	var svc = &executionService{
		MetaDB:        MetaDB,
		InnerQueue:    InnerQueue,
		Exporter:      Exporter,
		LockService:   lock.NewDBLockService(MetaDB),
		DomainService: domain.DefaultDomainService,
	}
	//binding executor
	StandardExecutor.ExecutionService = svc
	ExpressExecutor.ExecutionService = svc
	svc.StandardExecutor = StandardExecutor
	svc.ExpressExecutor = ExpressExecutor

	return svc
}

// SendExecutionEvents 发送event
func (svc *executionService) SendExecutionEvents(events ...vo.ExecutionEvent) {
	if svc.Exporter == nil {
		return
	}
	svc.Exporter.SendExecutionEvents(events)
}

// SendExecutionEvents 发送event
func (svc *executionService) SendInnerMessage(message queue.InnerMessageBody, sendtime *time.Time) error {

	if svc.InnerQueue == nil {
		return fmt.Errorf("inner queue is not initialized")
	}
	return svc.InnerQueue.SendInnerMessage(message, sendtime)

}

// NewExecutionFromID Create Execution by execution id
func (svc *executionService) NewExecutionFromID(id int, fields []string, tx rdb.Tx) (*Execution, error) {

	var dbExecution = &po.Execution{}
	dbExecution, err := svc.QueryExecutionByID(id, fields, tx)
	if err != nil {
		return nil, err
	}
	return NewExecutionFromData(dbExecution, svc)
}

// NewExecutionFromUUID Create Execution by execution uuid
func (svc *executionService) NewExecutionFromUUID(uuid string, fields []string, session rdb.Tx) (*Execution, error) {

	var dbExecution = &po.Execution{}

	dbExecution, err := svc.QueryExecutionByUUID(uuid, fields, session)

	if err != nil {
		return nil, err
	}
	exe, err := NewExecutionFromData(dbExecution, svc)
	return exe, err
}

// NewStepFromID Create ExecutionState by execution id
func (svc *executionService) NewStepFromID(step_id int, session rdb.Tx) (Step, error) {

	var dbStep = &po.Step{}
	var err error

	tx, maker := svc.MetaDB.NewTxMaker(session)
	defer maker.Close(&err)
	dbStep, err = svc.QueryStepByID(step_id, []string{}, tx)
	if err != nil {
		return nil, err
	}

	return NewStepFromData(dbStep, svc.StandardExecutor)
}

// NewTaskFromToken NewTaskFromToken
func (svc *executionService) NewTaskFromToken(token string, session rdb.Tx) (*Task, error) {

	var err error
	var tasktoken = po.TaskToken{
		Token: token,
	}
	var dbstep = &po.Step{}

	// 增加控制session
	tx, maker := svc.MetaDB.NewTxMaker(session)
	defer maker.Close(&err)

	err = tx.Where(tasktoken).Take(&tasktoken).Error
	if rdb.IsErrRecordNotFound(err) {
		return nil, fmt.Errorf("%w: %s", vo.ErrorTaskTokenNotFound, token)
	}
	if err != nil {
		return nil, err
	}

	dbstep, err = svc.QueryStepByID(tasktoken.StepID, []string{}, tx)
	if err != nil {
		return nil, err
	}

	state, err := NewTaskFromData(dbstep, svc.StandardExecutor)
	if err != nil {
		return nil, err
	}
	return state, err
}

func (svc *executionService) SendEventsMessages(events []vo.ExecutionEvent, msgs []queue.InnerMessageBody) error {

	svc.SendExecutionEvents(events...)

	for _, msg := range msgs {
		err := svc.InnerQueue.SendInnerMessage(msg, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

// func (svc *executionService) NewParallelFromID(step_id int, session rdb.Tx) (*Parallel, error) {

// 	step, err := NewParallelFromID(step_id, svc, session)
// 	return step, err

// }

// JudgeExecutionRunningStatus 判断Execution是否是可执行的状态
func (svc *executionService) JudgeExecutionRunningStatus(execution_id int) error {

	var err error
	dbexecution, err := svc.QueryExecutionByID(execution_id, ExecutionFields.L1, nil)
	if err != nil {
		return err
	}

	// 如果 消息类型是正常消息， 而 Execution状态不在 [ created  running ] , 忽略消息。 不能接收Failed/Abort/Success 等其他状态的消息
	if !slices.Contains([]string{string(ExecutionStatus.Running)},
		dbexecution.Status) {
		return fmt.Errorf("%w: current status '%s'", vo.ErrorExecutionStatus, dbexecution.Status)
	}
	return nil
}
