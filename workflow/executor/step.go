package executor

import (
	"github.com/skyflow-workflow/skyflow_backbend/workflow/repository/queue"
)

// Step 状态节点的抽象类
type Step interface {
	// 执行初始化
	Init(queue.InnerMessage) error
	// 执行
	Run(queue.InnerMessage) error
	// 获得bone信息
	GetBone() StepBone
	GetNextStep(output any) (NextStep, error)
	// 处理消息
	ProcessEvent(queue.InnerMessage) error
}

// // NewStepFromData ...
// func NewStepFromData(dbStep *po.Step, executor *Executor) (Step, error) {
// 	var step Step
// 	var err error

// 	wfstep, err := parser.ParserStep(dbStep.Definition, workflowtype)
// 	if err != nil {
// 		return nil, err
// 	}
// 	node := wfstep.GetNode()

// 	switch dbStep.Type {
// 	case grammar.StateType.Task:
// 		step, err = NewTaskFromState(dbStep, node.(*grammar.TaskState), svc)
// 	case grammar.StateType.Choice:
// 		step, err = NewChoiceFromState(dbStep, node.(*grammar.ChoiceState), svc)
// 	case grammar.StateType.Map:
// 		step, err = NewMapFromState(dbStep, node, svc)
// 	case grammar.StateType.Pass:
// 		step, err = NewPassFromState(dbStep, node.(*grammar.PassState), svc)
// 	case grammar.StateType.Parallel:
// 		step, err = NewParallelFromState(dbStep, node.(*grammar.ParallelState), svc)
// 	case grammar.StateType.StateGroup:
// 		step, err = NewStepGroupFromState(dbStep, node.(*grammar.StateGroup), svc)
// 	case grammar.StateType.Suspend:
// 		step, err = NewSuspendFromState(dbStep, node.(*grammar.SuspendState), svc)
// 	case grammar.StateType.Wait:
// 		step, err = NewWaitFromState(dbStep, node.(*grammar.WaitState), svc)
// 	// case grammer.StateType.Fail:
// 	// 	state, err = NewFailFromData(dbstep)
// 	// case grammer.StateType.Succeed:
// 	// 	state, err = NewSucceedFromData(dbstep)
// 	default:
// 		err = fmt.Errorf("state [ %s ] type : [ %s ] unrecognize ", dbStep.Name, dbStep.Type)

// 	}
// 	return step, err
// }
