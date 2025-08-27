package executor

type ExecutorConfig struct {
	Option *Option
	Quota  *Quota
}

// ParserConfig ...
type Option struct {
	// AllowActivity specifies whether to allow activity Task.
	AllowActivity bool
	// AllowWait specifies whether to allow Wait State.
	AllowWait bool
	// AllowSuspend specifies whether to allow Suspend State.
	AllowSuspend bool
	// AllowParallel specifies whether to allow Parallel State.
	AllowParallel bool
	// AllowMap specifies whether to allow Map State.
	AllowMap bool
	// AllowChoice specifies whether to allow Choice State.
	AllowChoice bool
	// AllowFail specifies whether to allow Fail State.
	AllowFail bool
	// AllowSucceed specifies whether to allow Succeed State.
	AllowSucceed bool
	// AllowPass specifies whether to allow Pass State.
	AllowPass bool
	// AbortOnFail specifies whether to stop execution on state failure.
	// in operation mode, it will try to execute all steps, user can retry the failed steps,
	// but in express mode, it will stop on the first failure. use fail fast mode.
	AbortOnFail bool
	// PersistenceStep specifies whether to persist step data.
	PersistenceStep bool
}

// Quota defines the quota for the workflow parser
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

// DefaultQuota default quota
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

// StandardExecutorConfig standard model workflow executor config
var StandardExecutorConfig = ExecutorConfig{
	Option: &Option{
		AllowActivity:   true,
		AllowWait:       true,
		AllowSuspend:    true,
		AllowParallel:   true,
		AllowMap:        true,
		AllowChoice:     true,
		AllowFail:       true,
		AllowSucceed:    true,
		AllowPass:       true,
		AbortOnFail:     false,
		PersistenceStep: true,
	},
	Quota: &DefaultQuota,
}

// ExpressParserConfig express model workflow executor config
var ExpressExecutorConfig = ExecutorConfig{
	Option: &Option{
		AllowActivity:   false,
		AllowWait:       false,
		AllowSuspend:    false,
		AllowParallel:   true,
		AllowMap:        true,
		AllowChoice:     true,
		AllowFail:       true,
		AllowSucceed:    true,
		AllowPass:       true,
		AbortOnFail:     true,
		PersistenceStep: false,
	},
	Quota: &DefaultQuota,
}
