package config

type Config struct {
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
	// EnableExecuteIndex specifies whether to enable execute index for step.
	// in standard mode, it will use execute index to track the execution order of steps.
	// in express mode, it will always not use execute index.
	EnableExecuteIndex bool
}
