// package quota  defined cloudflow quota here
package quota

// Quota ...
type Quota struct {
	// MaxStateNameSize specifies the maximum size of the state name.
	MaxStateNameSize int
	// MaxStateLimit specifies the maximum number of states allowed in a workflow.
	// If the number of states exceeds this limit, the parser will return an error.
	// If MaxStateLimit is 0, the parser will not check the number of states.
	MaxStateNumber int
	// MaxParallelBranchNumber specifies the maximum number of branches allowed in a parallel state.
	MaxParallelBranchNumber int
	// MaxMapBranchNumber specifies the maximum number of branches allowed in a map state.
	MaxMapBranchNumber int
	// MaxMapConncurrency specifies the maximum number of branches allowed in a map state.
	MaxMapConncurrency int
	// MaxStateMachineDepth specifies the maximum depth of the state machine.
	MaxStateMachineDepth int
	// MaxActivityURISize specifies the maximum size of the activity URI.
	MaxActivityURISize int
	// MaxStartExecutionInputSize specifies the maximum size of the start execution input.
	MaxStartExecutionInputSize int
	// MaxExecuteTimes specifies the maximum number of times a workflow can be executed.
	MaxExecuteTimes int
	// MaxExecutionUUIDSize specifies the maximum size of the execution UUID.
	MaxExecutionUUIDSize int
	//  MaxWorkflowSize specifies the maximum size of the workflow.
	MaxWorkflowSize int
	// MaxWorkflowURISize specifies the maximum size of the workflow URI.
	MaxWorkflowURISize int
	// MaxStepNameSize specifies the maximum size of the step name.
	MaxStepNameSize int
	// MaxRunningExecution specifies the maximum number of running
	// executions allowed in the system.
	MaxRunningExecution int
	// MaxTitleSize specifies the maximum size of the title.
	MaxTitleSize int
	// MaxInputSize specifies the maximum size of the input.
	MaxInputSize int
	// MaxOutputSize specifies the maximum size of the output.
	MaxOutputSize int
	// MaxStepDataSize specifies the maximum size of the step stored data.
	MaxStepDataSize int
	// MaxTaskInputSize specifies the maximum size of the task input.
	MaxTaskInputSize int
	// MaxTaskOutputSize specifies the maximum size of the task output.
	MaxTaskOutputSize int
}

// hard limit

// DefaultQuota ...
var DefaultQuota = Quota{
	MaxStateNameSize:           200,
	MaxStateNumber:             1000,
	MaxStateMachineDepth:       100,
	MaxParallelBranchNumber:    500,
	MaxMapBranchNumber:         500,
	MaxMapConncurrency:         100,
	MaxActivityURISize:         200,
	MaxStartExecutionInputSize: 128 * 1024,
	MaxExecuteTimes:            10000,
	MaxExecutionUUIDSize:       200,
	// max workflow size 512KB
	MaxWorkflowSize:     512 * 1024,
	MaxWorkflowURISize:  200,
	MaxStepNameSize:     200,
	MaxRunningExecution: 1000,
	MaxTitleSize:        200,
	// max input size 128KB
	MaxInputSize: 128 * 1024,
	// max output size 128KB
	MaxOutputSize: 128 * 1024,
	// max step data size 32KB
	MaxStepDataSize: 32 * 1024,
	// max task input size 32KB
	MaxTaskInputSize: 32 * 1024,
	// max task output size 32KB
	MaxTaskOutputSize: 32 * 1024,
}
