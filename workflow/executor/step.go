package executor

import (
	"fmt"

	"github.com/skyflow-workflow/skyflow_backend/workflow/parser/states"
	"github.com/skyflow-workflow/skyflow_backend/workflow/po"
	"github.com/skyflow-workflow/skyflow_backend/workflow/repository/queue"
)

// Step 状态节点的抽象类
type Step interface {
	// 执行初始化
	Init(queue.InnerMessageBody) error
	// 执行
	Run(queue.InnerMessageBody) error
	// 获得bone信息
	GetBone() StepBone
	GetNextStep(output any) (NextStep, error)
	// 处理消息
	ProcessEvent(queue.InnerMessageBody) error
}

// NewStepFromData ...
func NewStepFromData(dbStep *po.Step, executor *Executor) (Step, error) {
	var step Step
	var err error
	switch states.StateType(dbStep.Type) {
	case states.StateTypes.Task:
		step, err = NewTaskFromData(dbStep, executor)
	case states.StateTypes.Choice:
		step, err = NewChoiceFromData(dbStep, executor)
	// case grammar.StateType.Map:
	// 	step, err = NewMapFromState(dbStep, node, svc)
	case states.StateTypes.Pass:
		step, err = NewPassFromData(dbStep, executor)
	// case grammar.StateType.Parallel:
	// 	step, err = NewParallelFromState(dbStep, node.(*grammar.ParallelState), svc)
	// case grammar.StateType.StateGroup:
	// 	step, err = NewStepGroupFromState(dbStep, node.(*grammar.StateGroup), svc)
	// case grammar.StateType.Suspend:
	// 	step, err = NewSuspendFromState(dbStep, node.(*grammar.SuspendState), svc)
	case states.StateTypes.Wait:
		step, err = NewWaitFromData(dbStep, executor)
	// case grammer.StateType.Fail:
	// 	state, err = NewFailFromData(dbstep)
	// case grammer.StateType.Succeed:
	// 	state, err = NewSucceedFromData(dbstep)
	default:
		err = fmt.Errorf("step id: %d  name: %s type: %s  not supported",
			dbStep.ID, dbStep.Name, dbStep.Type)

	}
	return step, err
}
